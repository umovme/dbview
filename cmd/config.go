package cmd

import (
	"github.com/umovme/dbview/setup"
)

// Config contains all input from the config file
type Config struct {
	PGBinPath string                             `yaml:"pgsql-bin"`
	LocalDB   setup.ConnectionDetails            `yaml:"local-database"`
	Customers map[string]setup.ConnectionDetails `yaml:"customers"`
	RemoteDB  setup.ConnectionDetails            `yaml:"remote-database"`
	Options   map[string]string                  `yaml:"options"`
}
