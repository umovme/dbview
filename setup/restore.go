package setup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/go-playground/log"
)

// RestoreOptions : Define the options for restore a dump file into a database
type RestoreOptions struct {
	CustomArgs []string
}
/*
RestoreSQLFile : Calls the 'psql' to restore a gzipped file gererated by pg_dump + gzip
*/
func RestoreSQLFile(connDetail ConnectionDetails, dumpFile string, exists bool) error {
	//
	// new dumfile format:
	//
	// ---------------------------------------------------------------------
	// |                            tar file                               |
	// ---------------------------------------------------------------------
	// |   internal_dump_schema_dbview.gz | internal_dump_schema_user.gz   |
	// ---------------------------------------------------------------------

	psqlBin := "psql"
	psqlArgs := []string{"-v", "ON_ERROR_STOP=1", "-1", "-X"}
	psqlConn := formatConnectionOptions(connDetail)

	if pgsqlBinPATH != "" {
		psqlBin = fmt.Sprintf("%s/psql", pgsqlBinPATH)
	}

	psqlArgs = append(psqlArgs, psqlConn...)

	if connDetail.Password != "" {
		err := os.Setenv("PGPASSWORD", connDetail.Password)

		if err != nil {
			return err
		}
	}

	log.Debugf("%s %#v\n", psqlBin, psqlArgs)

	f, err := os.Open(dumpFile)
	defer f.Close()

	if err != nil {
		return err
	}

	tarFile := tar.NewReader(f)

	for {
		var out bytes.Buffer
		var stderr bytes.Buffer
		psqlCmd := exec.Command(psqlBin, psqlArgs...)
		psqlCmd.Stdout = &out
		psqlCmd.Stderr = &stderr

		hdr, err := tarFile.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Debugf("restoring: %s", hdr.Name)
		if exists && strings.HasPrefix(hdr.Name, "internal_dump_schema_dbview") {
			log.Warn("dbview schema already exists, ignoring to restore it")
			continue
		} else {
			reader, err := gzip.NewReader(tarFile)

			if err != nil {
				return err
			}

			psqlCmd.Stdin = reader

			if err := psqlCmd.Run(); err != nil {
				return fmt.Errorf(
					fmt.Sprintf(
						"%s. %s\nCMD: %s %s",
						fmt.Sprint(err),
						stderr.String(),
						psqlBin,
						psqlArgs))
			}
		}
	}

	return nil
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
		out = append(out, "--username="+connDetail.Username)
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
