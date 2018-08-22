// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/go-playground/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/umovme/dbview/setup"
)

func updateReplicationFunc() {
	localConn := setup.ConnectionDetails{
		Username: viper.GetString("local-database.username"),
		Host:     viper.GetString("local-database.host"),
		Port:     viper.GetInt("local-database.port"),
		Database: viper.GetString("local-database.target_database"),
		SslMode:  viper.GetString("local-database.ssl"),
		Password: viper.GetString("local-database.password"),
	}

	log.Info("Updating the database functions...")
	log.Debugf("Running %s", setup.ReplicationLogFunction)
	if err := setup.ExecuteQuery(localConn, setup.ReplicationLogFunction); err != nil {
		log.WithError(err).Error("error updating replication function on database")
	}
	log.Info("Done.")
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Updates replication functions on the target database",
	Run: func(cmd *cobra.Command, args []string) {
		updateReplicationFunc()
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
