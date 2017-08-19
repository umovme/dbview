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
		sampleConnString = "user=sebastian dbname=sebastian sslmode=verify-full"
	)

	Context("When I connect to the database", func() {

		It("Should convert a the connection info to the lib/pq connection", func() {

			Expect(dbConnectionInfo.toString()).To(Equal(sampleConnString))
		})

		It("Should not set a value in the connection string when it is empty", func() {

			dbConnectionInfo := ConnectionDetails{userName: "sebastian"}

			Expect(dbConnectionInfo.toString()).To(Equal("user=sebastian sslmode=verify-full"))
		})

		It("Should connect to a database", func() {
			err := connect(dbConnectionInfo)
			Expect(err).To(BeNil())
		})

	})

	Context("When I install the dbview at the database", func() {

		It("Should create a user in the database", func() {

			err := CreateUser(dbConnectionInfo, "dbview", []string{"SUPERUSER"})
			Expect(err).To(BeNil())
		})
	})
})
