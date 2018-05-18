package setup

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/log"
)

// RestoreOptions : Define the options for restore a dump file into a database
type RestoreOptions struct {
	CustomArgs []string
}

/*
RestoreDumpFile : Calls the 'pg_restore' to restore a dump file gererated by pg_dump
*/
func RestoreDumpFile(connDetail ConnectionDetails, dumpFile string, options RestoreOptions) error {

	pgRestoreBin := "pg_restore"

	if pgsqlBinPATH != "" {
		pgRestoreBin = fmt.Sprintf("%s/pg_restore", pgsqlBinPATH)
	}

	// conn := formatConnectionOptions(connDetail)
	args := formatConnectionOptions(connDetail)

	args = append(args, options.CustomArgs...)
	args = append(args, dumpFile)

	if connDetail.Password != "" {
		err := os.Setenv("PGPASSWORD", connDetail.Password)

		if err != nil {
			return err
		}
	}

	log.Debugf("%s %#v\n", pgRestoreBin, args)

	cmd := exec.Command(pgRestoreBin, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"%s. %s\nCMD: %s %s",
				fmt.Sprint(err),
				stderr.String(),
				pgRestoreBin,
				args))
	}

	return nil
}

func formatConnectionOptions(connDetail ConnectionDetails) []string {

	out := []string{}

	if connDetail.Username != "" {
		out = append(out, "--user="+connDetail.Username)
	}
	if connDetail.Host != "" {
		out = append(out, "--host="+connDetail.Host)
	}
	if connDetail.Database != "" {
		out = append(out, "--dbname="+connDetail.Database)
	}
	if connDetail.Port > 0 {
		out = append(out, fmt.Sprintf("--port=%d", connDetail.Port))
	}

	return out
}
