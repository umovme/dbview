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
	Username string
	Password string
	Host     string
	Database string
	SslMode  string
	port     int
}

func (c ConnectionDetails) toString() string {

	returnData := ""

	if c.Username != "" {
		returnData += fmt.Sprintf("user=%s ", c.Username)
	}

	if c.Password != "" {
		returnData += fmt.Sprintf("password=%s ", c.Password)
	}

	if c.Database != "" {
		returnData += fmt.Sprintf("dbname=%s ", c.Database)
	}

	if c.Host != "" {
		returnData += fmt.Sprintf("host=%s ", c.Host)
	}

	if c.SslMode != "" {
		returnData += fmt.Sprintf("sslmode=%s", c.SslMode)
	}

	return strings.Trim(returnData, " ")
}

func connect(connDetail ConnectionDetails) (*sql.DB, error) {
	return sql.Open("postgres", connDetail.toString())
}
