package setup

import (
	"database/sql"
	"fmt"
	"strings"

	// Required package for the PostgreSQL database
	_ "github.com/lib/pq"
)

var (
	pgsqlBinPATH string
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
	Port     int
}

var c *ConnectionDetails

// ToString : Creates a connection string in the lib/pq format
func ToString() string { return c.ToString() }
func (c ConnectionDetails) ToString() string {

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
		returnData += fmt.Sprintf("sslmode=%s ", c.SslMode)
	}

	if c.Port > 0 {
		returnData += fmt.Sprintf("port=%d ", c.Port)
	}

	return strings.Trim(returnData, " ")
}

func connect(connDetail ConnectionDetails) (*sql.DB, error) {
	return sql.Open("postgres", connDetail.ToString())
}

// SetPgsqlBinPath : Sets the PostgreSQL binary path
// This option exists to force the full binaries path when the
// binaries are not present in the OS PATH environment variable.
func SetPgsqlBinPath(path string) {
	pgsqlBinPATH = path
}
