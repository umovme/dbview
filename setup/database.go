package setup

import (
	"database/sql"
	"fmt"
	"strings"
	// needed for a PostgreSQL connection
	_ "github.com/lib/pq"
)

/*
ConnectionDetails : defines details about a new connection
*/
type ConnectionDetails struct {
	userName, host, database, sslmode string
	port                              int
}

func (c ConnectionDetails) toString() string {

	returnData := ""

	if c.userName != "" {
		returnData += fmt.Sprintf("user=%s ", c.userName)
	}

	if c.database != "" {
		returnData += fmt.Sprintf("dbname=%s ", c.database)
	}

	if c.host != "" {
		returnData += fmt.Sprintf("host=%s ", c.host)
	}

	if c.sslmode != "" {
		returnData += fmt.Sprintf("sslmode=%s", c.sslmode)
	}

	return strings.Trim(returnData, " ")
}

func connect(connDetail ConnectionDetails) (*sql.DB, error) {
	return sql.Open("postgres", connDetail.toString())
}

/*
CreateUser : Creates a new user in the database
*/
func CreateUser(connDetail ConnectionDetails, userName string, options []string) error {

	db, err := connect(connDetail)

	if err != nil {
		return err
	}

	exists, err := checkIfUserExists(connDetail, userName)

	if err != nil {
		return err
	} else if exists {
		return nil
	}

	_, err = db.Exec("CREATE USER dbview " + strings.Join(options, " ") + ";")

	if err != nil {
		return err
	}

	return nil
}

func checkIfUserExists(connDetail ConnectionDetails, userName string) (bool, error) {

	found := false
	db, err := connect(connDetail)

	if err != nil {
		return found, err
	}

	totalRows := 0
	err = db.QueryRow("SELECT count(1) FROM pg_roles WHERE rolname = $1", userName).Scan(&totalRows)

	if totalRows > 0 {
		found = true
	}

	return found, err
}
