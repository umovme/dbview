package setup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
RestoreOptions : Define the options for restore a dump file into a database
*/
type RestoreOptions struct {
	ExePath    string
	CustomArgs []string
}

/*
RestoreDumpFile : Calls the 'pg_restore' to restore a dump file gererated by pg_dump
*/
func RestoreDumpFile(connDetail ConnectionDetails, dumpFile string, options RestoreOptions) error {

	pgRestoreBin := "pg_restore"

	if options.ExePath != "" {
		pgRestoreBin = options.ExePath
	}

	args := fmt.Sprintf(
		"-U %s -d %s %s %s",
		connDetail.Username,
		connDetail.Database,
		strings.Join(options.CustomArgs, " "),
		dumpFile)

	if connDetail.Password != "" {
		err := os.Setenv("PGPASSWORD", connDetail.Password)

		if err != nil {
			return err
		}
	}
	/// ... at the right means turn the slide in a variadic variable
	cmd := exec.Command(pgRestoreBin, strings.Split(args, " ")...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ". " + stderr.String() + "(" + args + ")")
	}

	return nil
}
