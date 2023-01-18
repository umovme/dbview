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
	"time"

	"github.com/go-playground/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/umovme/dbview/setup"
)

const (
	daemonInterval = 3 * time.Second
	maxRowsAllowed = 10000
)

// replicateCmd represents the replicate command
var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "Runs the replication functions",
	Long:  `Runs the replication functions and updates the target database at the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		logInfoBold("Starting dbview replication")

		if viper.GetBool("daemon") {

			dur := viper.GetDuration("refresh-interval")
			if dur < daemonInterval {
				log.Fatalf("Refresh interval (%s) must greater or equals than %s.\n", dur, daemonInterval)
			}

			logInfoBold("Starting in daemon mode")

			ticker := time.NewTicker(dur)
			for ; true; <-ticker.C {
				runReplicate()
			}
		}
		updateReplicationFunc()
		runReplicate()

		log.Info("Done.")
	},
}

func runReplicate() {

	rowLimit := viper.GetInt32("options.row_limit")

	if rowLimit > maxRowsAllowed {
		log.Warnf("The maximum value for row limit is 10.000 rows. Fixing it in 10k rows.")
		rowLimit = maxRowsAllowed
	}

	localConn := setup.ConnectionDetails{
		Username: viper.GetString("local-database.username"),
		Host:     viper.GetString("local-database.host"),
		Port:     viper.GetInt("local-database.port"),
		Database: viper.GetString("local-database.target_database"),
		SslMode:  viper.GetString("local-database.ssl"),
		Password: viper.GetString("local-database.password"),
		AppName: "local",
	}

	log.Debugf("Using local connection with '%s'", localConn.ToString())

	remoteConn := setup.ConnectionDetails{
		Username: viper.GetString("remote-database.username"),
		Host:     viper.GetString("remote-database.host"),
		Port:     viper.GetInt("remote-database.port"),
		Database: viper.GetString("remote-database.database"),
		SslMode:  viper.GetString("remote-database.ssl"),
		Password: viper.GetString("remote-database.password"),
		AppName:  "dbview_client_" + viper.GetString("customer"),
	}

	log.Debugf("Using remote connection with '%s'", remoteConn.ToString())
	log.Debugf("Remember to use a remote user with '%s' in their search_path variable!", customerUser)

	newQuery := fmt.Sprintf(
		"SELECT do_replication_log('%s', 'u%s', %d);",
		remoteConn.ToString(),
		viper.GetString("customer"),
		rowLimit,
	)

	log.Debugf("QUERY: %s", newQuery)

	log.Info("Updating Replication Data...")
	if err := setup.ExecuteQuery(localConn, newQuery); err != nil {
		log.WithError(err).Error("fail to replicate the data")
	}
}

func init() {
	RootCmd.AddCommand(replicateCmd)

	// daemon mode related
	replicateCmd.PersistentFlags().Bool("daemon", false, "Run as daemon ")
	viper.BindPFlag("daemon", replicateCmd.PersistentFlags().Lookup("daemon"))

	replicateCmd.PersistentFlags().Duration("refresh-interval", daemonInterval, "Refresh interval for daemon mode")
	viper.BindPFlag("refresh-interval", replicateCmd.PersistentFlags().Lookup("refresh-interval"))

	replicateCmd.PersistentFlags().String("remote-database.ssl", "disable", fmt.Sprintf("Remote %s", sslConnectionLabel))
	viper.BindPFlag("remote-database.ssl", replicateCmd.PersistentFlags().Lookup("remote-database.ssl"))

	replicateCmd.PersistentFlags().StringP("remote-database.username", "", "postgres", fmt.Sprintf("Remote %s", dbUserLabel))
	viper.BindPFlag("remote-database.username", replicateCmd.PersistentFlags().Lookup("remote-database.username"))

	replicateCmd.PersistentFlags().StringP("remote-database.port", "", "9999", fmt.Sprintf("Remote %s", dbPortLabel))
	viper.BindPFlag("remote-database.port", replicateCmd.PersistentFlags().Lookup("remote-database.port"))

	replicateCmd.PersistentFlags().StringP("remote-database.password", "", "", fmt.Sprintf("Remote %s", dbUserPasswordLabel))
	viper.BindPFlag("remote-database.password", replicateCmd.PersistentFlags().Lookup("remote-database.password"))

	replicateCmd.PersistentFlags().StringP("remote-database.host", "", "dbview.umov.me", fmt.Sprintf("Remote %s", dbHostLabel))
	viper.BindPFlag("remote-database.host", replicateCmd.PersistentFlags().Lookup("remote-database.host"))

	replicateCmd.PersistentFlags().StringP("remote-database.database", "", "prod_umov_dbview", "Remote Database name")
	viper.BindPFlag("local-database.database", replicateCmd.PersistentFlags().Lookup("local-database.database"))

	replicateCmd.PersistentFlags().Int32P("options.row_limit", "l", 100, "row limit of each replication action")
	viper.BindPFlag("options.row_limit", replicateCmd.PersistentFlags().Lookup("options.row_limit"))

}
