package setup

import (
	"database/sql"
	"fmt"

	"github.com/apex/log"
)

// CreateExtensionsInDatabase : Creates extensions in the target database
func CreateExtensionsInDatabase(connDetail ConnectionDetails, extensions []string) error {

	var db *sql.DB
	var err error

	if db, err = connect(connDetail); err != nil {
		return err
	}
	defer db.Close()

	for _, extension := range extensions {

		available, err := checkExtensionInDatabase(connDetail, extension)
		if err != nil {
			return err
		}

		if available {

			query := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s;", extension)

			log.Debugf("Running %s", query)
			_, err = db.Exec(query)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("The extension '%s' are not available on the server.", extension)
		}

	}
	return nil
}

// checkExtensionInDatabase : Checks if a extension is available on the database server
// This is necessary to avoid errors when installing a extension which is not available at the server.
func checkExtensionInDatabase(connDetail ConnectionDetails, extension string) (bool, error) {

	var db *sql.DB
	var err error

	found := false

	if db, err = connect(connDetail); err != nil {
		return found, err
	}
	defer db.Close()

	totalRows := 0
	err = db.QueryRow("SELECT count(1) FROM pg_available_extensions WHERE name = $1", extension).Scan(&totalRows)

	return (totalRows > 0), err
}
