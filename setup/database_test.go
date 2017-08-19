package setup

import (
	// . "github.com/sebastianwebber/dbview/setup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Setup database user and groups", func() {

	var (
		dbConnectionInfo = ConnectionDetails{
			userName: "sebastian",
			database: "sebastian"}
		sampleConnString = "user=sebastian dbname=sebastian"

		testUserName  = "dbview"
		wrongUserName = "missing_user_for_this_database"
	)

	Context("When I connect to the database", func() {

		It("Should convert a the connection info to the lib/pq connection", func() {
			Expect(dbConnectionInfo.toString()).To(Equal(sampleConnString))
		})

		It("Should not set a value in the connection string when it is empty", func() {
			Expect(ConnectionDetails{userName: "sebastian", host: ""}.toString()).To(Equal("user=sebastian"))
		})

		It("Should connect to a database", func() {
			dbConnectionInfo.sslmode = "disable"
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
})
