package cmd

import (
	"runtime"
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
