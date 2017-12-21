package setup

import (
	"testing"
)

func TestCheckIfDatabaseExists(t *testing.T) {

	var testCases = []struct {
		description string
		dbName      string
		expected    bool
	}{
		{"Test on a existing database", "template1", true},
		{"Test on a existing empty database", "", false},
		{"Test on a non existing  database", "missing_db_from_server_125126616_1", false},
	}

	dbConnectionInfo := ConnectionDetails{
		Username: "seba",
		Database: "seba",
		SslMode:  "disable",
		Port:     5432}

	for _, tc := range testCases {
		t.Name()
		t.Log(tc.description)
		found, _ := checkIfDatabaseExists(dbConnectionInfo, tc.dbName)

		if found != tc.expected {
			t.Errorf("Expected %v got %v", found, tc.expected)
		}
	}
}
