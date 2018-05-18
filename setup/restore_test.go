package setup

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	testDB   = "super_test_db"
	testPort = 9876
	testHost = "127.0.0.1"
	testUser = "test_user"
)

func makeConnection(user, host, db, port bool) (out ConnectionDetails) {

	if user {
		out.Username = testUser
	}

	if host {
		out.Host = testHost
	}

	if db {
		out.Database = testDB
	}

	if port {
		out.Port = testPort
	}

	return
}

func Test_formatConnectionOptions(t *testing.T) {

	tests := []struct {
		name    string
		args    ConnectionDetails
		wantOut []string
	}{
		{name: "empty connection", args: makeConnection(false, false, false, false), wantOut: []string{}},
		{name: "custom connection with db and user", args: makeConnection(true, false, true, false), wantOut: []string{"--user=" + testUser, "--dbname=" + testDB}},
		{name: "custom connection with all options", args: makeConnection(true, true, true, true), wantOut: []string{"--user=" + testUser, "--host=" + testHost, "--dbname=" + testDB, fmt.Sprintf("--port=%d", testPort)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := formatConnectionOptions(tt.args); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("formatConnectionOptions() = %#v, want %#v", gotOut, tt.wantOut)
			}
		})
	}
}
