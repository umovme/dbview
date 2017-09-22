# DBView

This project contains a CLI to to help our customers replicate their data on it own servers.

## Build Details

### Depencies

This projects uses the packages bellow to build:

```bash
go get -v github.com/spf13/cobra
go get github.com/goreleaser/goreleaser
```

And (for now) have a ruby gem dependency:

```bash
gem install fpm
```

#### osx related dependecies:

```bash
brew install rpm
brew install dpkg
```
> More details here: http://timperrett.com/2014/03/23/enabling-rpmbuild-on-mac-osx/

### Deploying new versions

```bash
goreleaser --rm-dist
```

## How to use it

For now, this CLI only implements a possibility to install a new server. Further versions are planned.

### Installation

For the `install` process you need to input the dump file (sended by our support team), your customer id and the database credentials. For example:

```bash 
$ dbview install -U sebastian -P "seba" -c 1329 -D ~/tmp/dbview_dump_customer_1329_20170131.pgbkp
Validating parameters...
Creating the 'dbview' user
Creating the 'u1329' user
Fixing permissions
Updating the 'search_path'
Creating the 'umovme_dbview_db' database
Creating the necessary extensions
Restoring the dump file
Done.
```

A full detail of options are avaliable with the `--help` option. For example:

```bash
$ dbview install --help
Install all dependencies of the dbview environment like,
users, permissions, database and restores the database dump.

The database dump are provided by the uMov.me support team.

Please contact us with you have any trouble.

Usage:
  dbview install [flags]

Flags:
  -c, --customer int             Your customer ID
  -d, --database string          Database name (default "postgres")
  -D, --dump-file string         Database dump file
      --force-cleanup            Remove the database and user before starts (DANGER)
      --host string              Database host (default "127.0.0.1")
  -P, --password string          Username password
  -p, --port int                 Database port (default 5432)
  -S, --ssl-mode string          SSL connection: 'require', 'verify-full', 'verify-ca', and 'disable' supported (default "disable")
      --target-database string   The target database (default "umovme_dbview_db")
      --target-username string   The target username (default "dbview")
  -U, --username string          Database user (default "postgres")

Global Flags:
      --config string   config file (default is $HOME/.dbview.yaml)
```

## Getting support

If you need any help, contact us at [http://www.umov.me/suporte/](http://www.umov.me/suporte/).