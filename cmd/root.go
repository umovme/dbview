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
	"os"
	"runtime"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dbview",
	Short: "Manages the local dbview database",
	Long: `Manages the local dbview database providing tools to install,
configure and update the replication system.
	
Please contact us with you have any trouble.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// ugly messages on windows forces me to disable this pretty messages
	if runtime.GOOS != "windows" {
		log.SetHandler(text.Default)
	}
	log.SetLevel(log.InfoLevel)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbview.yaml)")

	RootCmd.PersistentFlags().String("local-database.ssl", "disable", sslConnectionLabel)
	viper.BindPFlag("local-database.ssl", RootCmd.PersistentFlags().Lookup("local-database.ssl"))

	RootCmd.PersistentFlags().StringP("local-database.username", "U", "postgres", f("Local %s", dbUserLabel))
	viper.BindPFlag("local-database.username", RootCmd.PersistentFlags().Lookup("local-database.username"))

	RootCmd.PersistentFlags().StringP("local-database.port", "p", "", f("Local %s", dbPortLabel))
	viper.BindPFlag("local-database.port", RootCmd.PersistentFlags().Lookup("local-database.port"))

	RootCmd.PersistentFlags().StringP("local-database.password", "P", "", f("Local %s", dbUserPasswordLabel))
	viper.BindPFlag("local-database.password", RootCmd.PersistentFlags().Lookup("local-database.password"))

	RootCmd.PersistentFlags().StringP("local-database.host", "h", "127.0.0.1", f("Local %s", dbHostLabel))
	viper.BindPFlag("local-database.host", RootCmd.PersistentFlags().Lookup("local-database.host"))

	RootCmd.PersistentFlags().StringP("local-database.database", "d", "postgres", "Local maintenance database. Used for administrative tasks.")
	viper.BindPFlag("local-database.database", RootCmd.PersistentFlags().Lookup("local-database.database"))

	RootCmd.PersistentFlags().String("local-database.target_database", "umovme_dbview_db", "Local target database.")
	viper.BindPFlag("local-database.target_database", RootCmd.PersistentFlags().Lookup("local-database.target_database"))

	RootCmd.PersistentFlags().String("local-database.target_username", "dbview", "Local target username.")
	viper.BindPFlag("local-database.target_username", RootCmd.PersistentFlags().Lookup("local-database.target_username"))

	RootCmd.PersistentFlags().Int("customer", 0, "Your customer ID")
	viper.BindPFlag("customer", RootCmd.PersistentFlags().Lookup("customer"))

	RootCmd.PersistentFlags().String("pgsql-bin", "", "PostgreSQL binaries PATH")
	viper.BindPFlag("pgsql-bin", RootCmd.PersistentFlags().Lookup("pgsql-bin"))

	RootCmd.PersistentFlags().Bool("help", false, "Show this help message")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".dbview") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file: " + viper.ConfigFileUsed())
	} else {
		log.WithError(err).Fatal("fail to load the config file")
	}
}
