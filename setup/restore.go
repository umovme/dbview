package setup

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

/*
RestoreOptions : Define the options for restore a dump file into a database
*/
type RestoreOptions struct {
	exePath    string
	customArgs []string
}

/*
RestoreDumpFile : Calls the 'pg_restore' to restore a dump file gererated by pg_dump
*/
func RestoreDumpFile(connDetail ConnectionDetails, dumpFile string, options RestoreOptions) error {

	pgRestoreBin := "pg_restore"

	if options.exePath != "" {
		pgRestoreBin = options.exePath
	}

	args := fmt.Sprintf(
		"-U %s -d %s %s %s",
		connDetail.userName,
		connDetail.database,
		strings.Join(options.customArgs, " "),
		dumpFile)

	cmd := exec.Command(pgRestoreBin, strings.Split(args, " ")...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	return nil
}
