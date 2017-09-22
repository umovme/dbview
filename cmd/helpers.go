package cmd

import (
	"strings"

	"github.com/apex/log"
)

// abort: aborts this program on any error
func abort(err error) {
	if err != nil {
		log.WithError(err).Fatal("Failed")
	}
}

func logInfoBold(message string) {
	log.Infof("\033[1m%s\033[0m", strings.ToUpper(message))
}

func logWarnBold(message string) {
	log.Warnf("\033[1m%s\033[0m", strings.ToUpper(message))
}
