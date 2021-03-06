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
DROP FUNCTION IF EXISTS do_replication_log(TEXT, TEXT);
CREATE OR REPLACE FUNCTION do_replication_log(
	remote_connection_info TEXT,
	schema_name TEXT,
	rows_limit INTEGER
)
RETURNS TEXT AS
$$
-- Version 0.7
DECLARE
	rDeltas 					RECORD;
	remote_connection_id		TEXT;
	applied_deltas 				INTEGER DEFAULT 0;
	query						TEXT;
	last_transactionlog			BIGINT;
	remote_transaction_count 	BIGINT;
	start_time 					TIMESTAMPTZ;
	elapsed_time 				INTERVAL;
BEGIN
	-- LOCK to prevent concurrent running in the same environment
	IF pg_try_advisory_lock(substr(schema_name,2)::bigint) IS FALSE THEN
		RAISE EXCEPTION '(%) Replication already running for this customer', schema_name;
	END IF;

	remote_connection_id 	:= 'do_remote_replication_log';
	start_time 				:= clock_timestamp();

	RAISE LOG '(%) Stablishing REMOTE connection to uMov.me', schema_name;
	-- Connect to the remote host (uMov.me)
	PERFORM public.dblink_connect(remote_connection_id, remote_connection_info);

	-- Adjust local search_path
	PERFORM set_config('search_path', schema_name || ', dbview, public', true);

	RAISE LOG '(%) Getting last applied transaction to check your DBView replica consistency', schema_name;
	-- Get Last Applied TransactionLog
	EXECUTE FORMAT('SELECT  COALESCE(max(trl_id), 0) FROM %I.transactionlog', schema_name)
	INTO  last_transactionlog;
  
	RAISE LOG '(%) Cleaning up "transactionlog" older than %', schema_name, last_transactionlog;

	EXECUTE FORMAT('DELETE FROM  %I.transactionlog  WHERE trl_id < %L', schema_name, last_transactionlog);

	-- Query to get deltas to be applied in local copy
	SELECT INTO QUERY
		FORMAT($QUERY$
SELECT
  trl_id,
  trl_datehour,
  CASE WHEN trl_statements ~ '^BEGIN;' THEN substr(trl_statements, 8, length(trl_statements)-15) ELSE trl_statements END,
  trl_txid
FROM %I.transactionlog
WHERE trl_id > %s
ORDER BY trl_id
LIMIT %s;
$QUERY$, schema_name, last_transactionlog, rows_limit);

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

		RAISE LOG '(%) . Applying delta % from dbview remote transactionlog table', schema_name, rDeltas.trl_id;

		EXECUTE rDeltas.trl_statements;

		EXECUTE FORMAT('INSERT INTO %I.transactionlog(trl_id, trl_datehour, trl_statements, trl_txid) VALUES (%s,%L,%L,%s)',
		schema_name, rDeltas.trl_id, rDeltas.trl_datehour, rDeltas.trl_statements, rDeltas.trl_txid);
  
		applied_deltas := applied_deltas + 1;
	END LOOP;

	PERFORM public.dblink_disconnect(remote_connection_id);
	PERFORM pg_advisory_unlock(substr(schema_name,2)::bigint);

	elapsed_time := clock_timestamp() - start_time;

	RAISE LOG '(%) Applied % deltas from dbview remote transactionlog table (elapsed %)', schema_name, applied_deltas, elapsed_time;

	RETURN format('(%s) Applied %s deltas from dbview remote transactionlog table (elapsed %s)', schema_name, applied_deltas::text, elapsed_time::text);
END;
$$
LANGUAGE plpgsql;
`
