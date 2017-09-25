# DBView

This project contains a CLI to to help our customers replicate their data on it own servers.

## Build Details

### Depencies

This projects uses the packages bellow to build:

```bash
go get github.com/spf13/cobra
go get github.com/goreleaser/goreleaser
go get github.com/apex/log
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
$ dbview install -D ~/tmp/dbview_dump_customer_1329_20170131.pgbkp
  INFO[0000] Using config file: /Users/sebastian/.dbview.toml
  INFO[0000] INSTALLING DBVIEW AND DEPENDENCIES
  INFO[0000] Validating parameters...
  INFO[0000] STARTING UP
  INFO[0000] Creating the 'dbview' user
  INFO[0000] Creating the 'u1329' user
  INFO[0000] Fixing permissions
  INFO[0000] Updating the 'search_path'
  INFO[0000] Creating the 'umovme_dbview_db' database
  INFO[0000] Creating the necessary extensions
  INFO[0002] Restoring the dump file
  INFO[0004] Done.
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
      --dump-file string   Database dump file
      --force-cleanup      Remove the database and user before starts (DANGER)

Global Flags:
      --config string                           config file (default is $HOME/.dbview.yaml)
      --customer int                            Your customer ID
      --help                                    Show this help message
  -d, --local-database.database string          Local maintenance database. Used for administrative tasks. (default "postgres")
  -h, --local-database.host string              Local Database host (default "127.0.0.1")
  -P, --local-database.password string          Local Database password
  -p, --local-database.port string              Local Database password
      --local-database.ssl string               Local SSL connection: 'require', 'verify-full', 'verify-ca', and 'disable' supported (default "disable")
      --local-database.target_database string   Local target database. (default "umovme_dbview_db")
      --local-database.target_username string   Local target username. (default "dbview")
  -U, --local-database.username string          Local Database user (default "postgres")
```

## Getting support

If you need any help, contact us at [http://www.umov.me/suporte/](http://www.umov.me/suporte/).