package setup

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
CreateUser : Creates a new user in the database
*/
func CreateUser(connDetail ConnectionDetails, userName string, options []string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if exists {
		// returns if the user already exists
		return nil
	}
	_, err = db.Exec(fmt.Sprintf("CREATE USER %s %s;", userName, strings.Join(options, " ")))
	return err
}

/*
DropUser : Drops a user in the database
*/
func DropUser(connDetail ConnectionDetails, userName string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if exists {
		// returns if the user already exists
		return nil
	}
	_, err = db.Exec(
		fmt.Sprintf("DROP USER IF EXISTS %s;", userName))
	return err
}

func checkIfUserExists(connDetail ConnectionDetails, userName string) (bool, error) {

	var db *sql.DB
	var err error

	found := false

	if db, err = connect(connDetail); err != nil {
		return found, err
	}
	defer db.Close()

	totalRows := 0
	err = db.QueryRow("SELECT count(1) FROM pg_roles WHERE rolname = $1", userName).Scan(&totalRows)

	return (totalRows > 0), err
}

/*
GrantRolesToUser : Grant some roles privilieges for a user
*/
func GrantRolesToUser(connDetail ConnectionDetails, userName string, roles []string) error {
	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if !exists {
		// returns if the user not exists
		return nil
	}

	for _, role := range roles {
		_, err = db.Exec(fmt.Sprintf("GRANT %s to %s;", role, userName))
		if err != nil {
			return err
		}
	}

	return nil
}

/*
SetSearchPathForUser : Set the 'search_path' variable for an user
*/
func SetSearchPathForUser(connDetail ConnectionDetails, userName string, schemas []string) error {

	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	var exists bool

	if exists, err = checkIfUserExists(connDetail, userName); err != nil {
		return err
	} else if !exists {
		// returns if the user not exists
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("ALTER ROLE %s SET search_path TO %s;", userName, strings.Join(schemas, ",")))
	return err
}
