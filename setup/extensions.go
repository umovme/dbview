package setup

import (
	"database/sql"
	"fmt"
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
		_, err = db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s;", extension))
		if err != nil {
			return err
		}
	}
	return nil
}
