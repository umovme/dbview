package setup

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setup database user and groups", func() {

	var (
		dbConnectionInfo = ConnectionDetails{
			Username: "sebastian",
			Database: "sebastian"}
		sampleConnString = "user=sebastian dbname=sebastian"

		testUserName  = "dbview"
		wrongUserName = "missing_user_for_this_database"

		createTempDBName = func() string {
			hasher := md5.New()
			hasher.Write([]byte(time.Now().Local().Format(time.UnixDate)))
			return "temp_" + hex.EncodeToString(hasher.Sum(nil))
		}
	)

	Context("When I connect to the database", func() {

		It("Should convert a the connection info to the lib/pq connection", func() {
			Expect(dbConnectionInfo.toString()).To(Equal(sampleConnString))
		})

		It("Should not set a value in the connection string when it is empty", func() {
			Expect(ConnectionDetails{Username: "sebastian", Host: ""}.toString()).To(Equal("user=sebastian"))
		})

		It("Should connect to a database", func() {
			dbConnectionInfo.SslMode = "disable"
			_, err := connect(dbConnectionInfo)

			Expect(err).To(BeNil())
		})

	})

	Context("When I Create a user at the database", func() {

		It("It check if user exists before try to create a new one", func() {
			_, err := checkIfUserExists(dbConnectionInfo, testUserName)
			Expect(err).To(BeNil())
		})

		It("Should create a user in the database", func() {
			err := CreateUser(dbConnectionInfo, testUserName, []string{"PASSWORD 'super_senha'", "SUPERUSER"})
			Expect(err).To(BeNil())
		})

		It("Should grant some roles to a new user", func() {
			err := GrantRolesToUser(dbConnectionInfo, "sebastian", []string{"dbview"})
			Expect(err).To(BeNil())
		})

		It("Should check if user exists before grant a role", func() {
			err := GrantRolesToUser(dbConnectionInfo, wrongUserName, []string{"dbview"})
			Expect(err).To(BeNil())
		})

		It("Should be possible to update the 'search_path' configuration for the user", func() {
			err := SetSearchPathForUser(dbConnectionInfo, "sebastian", []string{"dbview", "public"})
			Expect(err).To(BeNil())
		})

		It("Should check if user exists before update the 'search_path'", func() {
			err := SetSearchPathForUser(dbConnectionInfo, wrongUserName, []string{"dbview", "public"})
			Expect(err).To(BeNil())
		})
	})

	Context("When I create a database", func() {
		It("Should create a new database", func() {
			err := CreateNewDatabase(dbConnectionInfo, "dbview", nil)
			Expect(err).To(BeNil())
		})

		It("Should drop a database", func() {
			err := DropDatabase(dbConnectionInfo, "dbview")
			Expect(err).To(BeNil())
		})

		It("Should check if the database exists before create a new one", func() {
			_, err := checkIfDatabaseExists(dbConnectionInfo, "template1")
			Expect(err).To(BeNil())
		})

		It("Should create a some extensions", func() {
			err := CreateExtensionsInDatabase(dbConnectionInfo, []string{"plpgsql"})
			Expect(err).To(BeNil())
		})

		It("Should check if a schema exists", func() {
			_, err := CheckIfSchemaExists(dbConnectionInfo, "public")
			Expect(err).To(BeNil())
		})

		It("Should remove a schema and all of this contents", func() {
			err := RemoveSchema(dbConnectionInfo, "u1325")
			Expect(err).To(BeNil())
		})

		It("Should create a new schema", func() {
			err := CreateSchema(dbConnectionInfo, "public")
			Expect(err).To(BeNil())
		})
	})

	Context("When I restore a database", func() {

		It("Should restore a dump file", func() {

			// dbConnectionInfo.Database = "sebastian"

			tempDBName := createTempDBName()
			var err error
			err = CreateNewDatabase(dbConnectionInfo, tempDBName, nil)

			dbConnectionInfo.Database = tempDBName

			options := RestoreOptions{exePath: "/usr/local/bin/pg_restore", customArgs: []string{"-Fc"}}
			err = RestoreDumpFile(dbConnectionInfo, "/Users/sebastian/tmp/file.dump", options)
			Expect(err).To(BeNil())

			dbConnectionInfo.Database = "postgres"
			err = DropDatabase(dbConnectionInfo, tempDBName)
			Expect(err).To(BeNil())
		})
	})

	Context("When I connect to the database", func() {

		It("Should run a query", func() {
			err := ExecuteQuery(dbConnectionInfo, "select 1;")
			Expect(err).To(BeNil())
		})

		It("Should create the replication function", func() {

			tempDBName := createTempDBName()
			var err error

			// step 1: create a new database
			err = CreateNewDatabase(dbConnectionInfo, tempDBName, nil)
			Expect(err).To(BeNil())

			// step 2: import the functions
			err = ExecuteQuery(dbConnectionInfo, ReplicationLogFunction)
			Expect(err).To(BeNil())

			// step 3: drop the new database
			dbConnectionInfo.Database = "postgres"
			err = DropDatabase(dbConnectionInfo, tempDBName)
			Expect(err).To(BeNil())
		})
	})
})
