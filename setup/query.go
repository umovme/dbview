package setup

import "database/sql"

/*
ExecuteQuery : Runs a query at the database
*/
func ExecuteQuery(connDetail ConnectionDetails, query string) error {

	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query)
	return err
}

/*
ReplicationLogFunction : Contains a SQL Command to create the replication function
*/
const ReplicationLogFunction string = `DROP FUNCTION IF EXISTS do_replication_log(TEXT, TEXT, TEXT, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS do_replication_log(TEXT, TEXT, TEXT);
CREATE OR REPLACE FUNCTION do_replication_log(
	remote_connection_info TEXT,
	schema_name TEXT
)
RETURNS TEXT AS
$$
DECLARE
	rDeltas 				RECORD;
	remote_connection_id	TEXT;
	rows_limit 				INTEGER DEFAULT 100;
	applied_deltas 			INTEGER DEFAULT 0;
	query					TEXT;
	last_transactionlog		BIGINT;
BEGIN
	-- LOCK to prevent concurrent running in the same environment
	IF pg_try_advisory_lock(substr(schema_name,2)::bigint) IS FALSE THEN
		RAISE EXCEPTION '(%) Replication already running for this customer', schema_name;
	END IF;

	remote_connection_id := 'do_remote_replication_log';

	RAISE LOG '(%) Stablishing REMOTE connection to uMov.me', schema_name;
	-- Connect to the remote host (uMov.me)
	PERFORM public.dblink_connect_u(remote_connection_id, remote_connection_info);

	-- Adjust local search_path
	PERFORM set_config('search_path', schema_name || ', dbview, public', true);

	RAISE LOG '(%) Getting last applied transaction to check your DBView replica consistency', schema_name;
	-- Get Last Applied TransactionLog
	SELECT	COALESCE(max(trl_id), 0)
	INTO	last_transactionlog
	FROM	transactionlog;

	-- Query to get deltas to be applied in local copy
	query := 'SELECT trl_id, trl_datehour, ';
	query := query || E' CASE WHEN trl_statements ~ \'^BEGIN;\' THEN substr(trl_statements, 8, length(trl_statements)-15) ELSE trl_statements END, ';
	query := query || ' trl_txid FROM transactionlog ';
	query := query || ' WHERE trl_id > '|| last_transactionlog || ' ORDER BY trl_id LIMIT ' || rows_limit;

	RAISE LOG '(%) Getting last % deltas do be applied in your local copy of DBView', schema_name, rows_limit;
	FOR rDeltas IN
		SELECT	*
		FROM	public.dblink(remote_connection_id, query)
				AS transaction(
					trl_id			BIGINT,
					trl_datehour	TIMESTAMPTZ,
					trl_statements	TEXT,
					trl_txid		BIGINT
				)
	LOOP
		RAISE DEBUG '(%) %', schema_name, rDeltas;

		-- Check the order of the remote and local transactionlog do be applied
		IF applied_deltas = 0 AND rDeltas.trl_id <> (last_transactionlog + 1) AND last_transactionlog != 0 THEN
			PERFORM public.dblink_disconnect(remote_connection_id);
			RAISE EXCEPTION
				'(%) Expected transaction % does not exist in remote host. Please contact the uMov.me Support Team to get a new dump!',
				schema_name, (last_transactionlog + 1);
		END IF;

		RAISE LOG '(%) . Applying delta % from dbview remote transactionlog table', schema_name, rDeltas.trl_id;

		EXECUTE rDeltas.trl_statements;

		INSERT INTO transactionlog(trl_id, trl_datehour, trl_statements, trl_txid)
		VALUES (rDeltas.trl_id, rDeltas.trl_datehour, rDeltas.trl_statements, rDeltas.trl_txid);

		applied_deltas := applied_deltas + 1;
	END LOOP;

	PERFORM public.dblink_disconnect(remote_connection_id);
	PERFORM pg_advisory_unlock(substr(schema_name,2)::bigint);

	RAISE LOG '(%) Applied % deltas from dbview remote transactionlog table', schema_name, applied_deltas;

	RETURN format('(%s) Applied %s deltas from dbview remote transactionlog table', schema_name, applied_deltas::text);
END;
$$
LANGUAGE plpgsql;
`
