package setup

import (
	"fmt"
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

	if c.sslmode == "" {
		c.sslmode = "verify-full"
	}

	if c.sslmode != "" {
		returnData += fmt.Sprintf("sslmode=%s", c.sslmode)
	}

	return returnData
}

func connect(connDetail ConnectionDetails) error {
	// db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	return nil
}

/*
CreateUser : Creates a new user in the database
*/
func CreateUser(connDetail ConnectionDetails, userName string, options []string) error {
	return nil
}
