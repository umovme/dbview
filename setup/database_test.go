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

		testUserName = "dbview"
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
			db, err := connect(dbConnectionInfo)
			db.Close()
			Expect(err).To(BeNil())
		})

	})

	Context("When I Create the dbview user at the database", func() {

		It("It check if user exists before try to create a new one", func() {
			_, err := checkIfUserExists(dbConnectionInfo, testUserName)
			Expect(err).To(BeNil())
		})

		It("Should create a user in the database if not exists", func() {
			err := CreateUser(dbConnectionInfo, testUserName, []string{"PASSWORD 'super_senha'", "SUPERUSER"})
			Expect(err).To(BeNil())
		})
	})
})
