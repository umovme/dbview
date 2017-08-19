package setup

import (
	"database/sql"
	"fmt"
	"strings"

	// Required package for the PostgreSQL database
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
