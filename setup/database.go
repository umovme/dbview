package setup

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
CreateNewDatabase : Creates a new database
*/
func CreateNewDatabase(connDetail ConnectionDetails, dbName string, options []string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}

	var exists bool

	if exists, err = checkIfDatabaseExists(connDetail, dbName); err != nil {
		return err
	} else if exists {
		// returns if the database already exists
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s %s;", dbName, strings.Join(options, " ")))
	return err
}

func checkIfDatabaseExists(connDetail ConnectionDetails, dbName string) (bool, error) {

	var db *sql.DB
	var err error

	found := false

	if db, err = connect(connDetail); err != nil {
		return found, err
	}

	totalRows := 0
	err = db.QueryRow("SELECT count(1) FROM pg_database WHERE datname = $1", dbName).Scan(&totalRows)

	return (totalRows > 0), err
}

/*
CreateExtensionsInDatabase : Creates extensions in the target database
*/
func CreateExtensionsInDatabase(connDetail ConnectionDetails, extensions []string) error {

	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}

	for _, extension := range extensions {
		_, err = db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s;", extension))
		if err != nil {
			return err
		}
	}
	return nil
}
