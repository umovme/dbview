package cmd

import (
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/apex/log"
)

const (
	sslConnectionLabel  string = "SSL connection: 'require', 'verify-full', 'verify-ca', and 'disable' supported"
	dbUserLabel         string = "Database User"
	dbUserPasswordLabel string = "Database password"
	dbPortLabel         string = "Database Port"
	dbHostLabel         string = "Database Host"
)

func preventAbort() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Fatalf("Captured %v event, aborting...", sig)
		}
	}()
}

var customerUser string

// abort: aborts this program on any error
func abort(err error) {
	if err != nil {
		log.WithError(err).Fatal("Failed")
	}
}

func logInfoBold(message string) {
	if runtime.GOOS != "windows" {
		log.Infof("\033[1m%s\033[0m", strings.ToUpper(message))
	} else {
		log.Infof("%s", strings.ToUpper(message))
	}
}

func logWarnBold(message string) {
	if runtime.GOOS != "windows" {
		log.Warnf("\033[1m%s\033[0m", strings.ToUpper(message))
	} else {
		log.Warnf("%s", strings.ToUpper(message))
	}
}
