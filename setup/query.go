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

	_, err = db.Exec(query)
	return err
}
