// Copyright © 2017 Sebastian Webber <sebastian@swebber.me>
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

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/umovme/dbview/setup"
)

// replicateCmd represents the replicate command
var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "Runs the replication functions",
	Long:  `Runs the replication functions and updates the target database at the latest version`,
	Run: func(cmd *cobra.Command, args []string) {

		localConn := setup.ConnectionDetails{
			Username: viper.GetString("local-database.username"),
			Host:     viper.GetString("local-database.host"),
			Port:     viper.GetInt("local-database.port"),
			Database: viper.GetString("local-database.target_database"),
			SslMode:  viper.GetString("local-database.ssl"),
			Password: viper.GetString("local-database.password")}

		remoteConn := setup.ConnectionDetails{
			Username: viper.GetString("remote-database.username"),
			Host:     viper.GetString("remote-database.host"),
			Port:     viper.GetInt("remote-database.port"),
			Database: viper.GetString("remote-database.database"),
			SslMode:  viper.GetString("remote-database.ssl"),
			Password: viper.GetString("remote-database.password")}

		newQuery := fmt.Sprintf(
			"SELECT do_replication_log('%s', '%s', '%s');",
			remoteConn.ToString(),
			localConn.ToString(),
			fmt.Sprintf("u%s", viper.GetString("customer")))

		log.Info("Updating Replication Data...")
		abort(
			setup.ExecuteQuery(localConn, newQuery))

		log.Info("Done.")
	},
}

func init() {
	RootCmd.AddCommand(replicateCmd)

	replicateCmd.PersistentFlags().String("remote-database.ssl", "disable", f("Remote %s", sslConnectionLabel))
	viper.BindPFlag("remote-database.ssl", replicateCmd.PersistentFlags().Lookup("remote-database.ssl"))

	replicateCmd.PersistentFlags().StringP("remote-database.username", "", "postgres", f("Remote %s", dbUserLabel))
	viper.BindPFlag("remote-database.username", replicateCmd.PersistentFlags().Lookup("remote-database.username"))

	replicateCmd.PersistentFlags().StringP("remote-database.port", "", "9999", f("Remote %s", dbPortLabel))
	viper.BindPFlag("remote-database.port", replicateCmd.PersistentFlags().Lookup("remote-database.port"))

	replicateCmd.PersistentFlags().StringP("remote-database.password", "", "", f("Remote %s", dbUserPasswordLabel))
	viper.BindPFlag("remote-database.password", replicateCmd.PersistentFlags().Lookup("remote-database.password"))

	replicateCmd.PersistentFlags().StringP("remote-database.host", "", "dbview.umov.me", f("Remote %s", dbHostLabel))
	viper.BindPFlag("remote-database.host", replicateCmd.PersistentFlags().Lookup("remote-database.host"))

	replicateCmd.PersistentFlags().StringP("remote-database.database", "", "prod_umov_dbview", "Remote Database name")
	viper.BindPFlag("local-database.database", replicateCmd.PersistentFlags().Lookup("local-database.database"))

}
