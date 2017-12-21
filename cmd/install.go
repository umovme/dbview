// Copyright © 2017 uMov.me Team <devteam-umovme@googlegroups.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/umovme/dbview/setup"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the dbview in the database",
	Long: `

	Install all dependencies of the dbview environment like, 
users, permissions, database and restores the database dump.
	
The database dump are provided by the uMov.me support team. 
	
Please contact us with you have any trouble.`,
	Run: func(cmd *cobra.Command, args []string) {

		logInfoBold("Installing dbview and dependencies")

		log.Info("Validating parameters...")
		if !checkInputParameters() {
			return
		}

		// fmt.Println(viper.GetString("local-database.ssl"), viper.GetString("author"), viper.GetString("local-database.ssl"))
		// return

		conn := setup.ConnectionDetails{
			Username: viper.GetString("local-database.username"),
			Host:     viper.GetString("local-database.host"),
			Port:     viper.GetInt("local-database.port"),
			Database: viper.GetString("local-database.database"),
			SslMode:  viper.GetString("local-database.ssl"),
			Password: viper.GetString("local-database.password")}

		cleanup(conn, customerUser)

		logInfoBold("Starting up")
		for _, user := range []string{viper.GetString("local-database.target_username"), customerUser} {
			log.Infof("Creating the '%s' user", user)
			abort(
				setup.CreateUser(conn, user, nil))
		}

		log.Info("Fixing permissions")
		abort(
			setup.GrantRolesToUser(conn, customerUser, []string{viper.GetString("local-database.target_username")}))

		log.Info("Updating the 'search_path'")
		abort(
			setup.SetSearchPathForUser(conn, customerUser, []string{customerUser, "public"}))

		log.Infof("Creating the '%s' database", viper.GetString("local-database.target_database"))
		abort(
			setup.CreateNewDatabase(conn, viper.GetString("local-database.target_database"), []string{"OWNER " + viper.GetString("local-database.target_username"), "TEMPLATE template0"}))

		log.Info("Creating the necessary extensions")
		conn.Database = viper.GetString("local-database.target_database")

		abort(
			setup.CreateExtensionsInDatabase(conn, []string{"hstore", "dblink", "pg_freespacemap", "postgis", "tablefunc", "unaccent"}))

		exists, err := setup.CheckIfSchemaExists(conn, "dbview")
		abort(err)

		restoreArgs := []string{"-Fc"}

		if exists {
			// if exists the dbview schema, this is not a first user schema on this database
			// then just create a new schema and restore only it
			abort(
				setup.CreateSchema(conn, customerUser))

			restoreArgs = append(restoreArgs, fmt.Sprintf("--schema=%s", customerUser))
		}

		pgPath := viper.GetString("pgsql-bin")

		if pgPath != "" {
			setup.SetPgsqlBinPath(pgPath)
		}

		log.Info("Restoring the dump file")
		abort(
			setup.RestoreDumpFile(conn, pDumpFile, setup.RestoreOptions{CustomArgs: restoreArgs}))

		log.Info("Installing the database functions")

		abort(
			setup.ExecuteQuery(conn, setup.ReplicationLogFunction))

		log.Info("Done.")
	},
}

func checkInputParameters() bool {

	if viper.GetInt("customer") == 0 {
		fmt.Println("Missing the customer id!")
		return false
	}

	if pDumpFile == "" {
		fmt.Println("Missing the dump file!")
		return false

	}

	return true
}

func cleanup(conn setup.ConnectionDetails, customerUser string) {
	if pCleanInstall {

		logWarnBold("Cleanup old stuff")

		log.Warnf("Dropping the '%s' database", viper.GetString("local-database.target_database"))
		abort(
			setup.DropDatabase(conn, viper.GetString("local-database.target_database")))
		for _, user := range []string{viper.GetString("local-database.target_username"), customerUser} {
			log.Warnf("Dropping the '%s' user", user)
			abort(
				setup.DropUser(conn, user))
		}
	}
}

var (
	pCleanInstall bool
	pDumpFile     string
)

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&pCleanInstall, "force-cleanup", "", false, "Remove the database and user before starts (DANGER)")
	installCmd.Flags().StringVar(&pDumpFile, "dump-file", "", "Database dump file")

	installCmd.PersistentFlags().String("local-database.target_database", "umovme_dbview_db", "Local target database.")
	viper.BindPFlag("local-database.target_database", installCmd.PersistentFlags().Lookup("local-database.target_database"))

	installCmd.PersistentFlags().String("local-database.target_username", "dbview", "Local target username.")
	viper.BindPFlag("local-database.target_username", installCmd.PersistentFlags().Lookup("local-database.target_username"))

}
