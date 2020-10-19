# DBView

This project contains a CLI to help our customers replicate their data on it own servers.

## Build Details

### Depencies

This projects uses the packages bellow to build:

```bash
make tools
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
./publish TAG
```
> Remember to use a tag in the [semantic versioning](http://semver.org/) pattern.

## How to use it

For now, this CLI only implements a possibility to install a new server and update the replication. Further versions are planned.

### Installation

Database extensions required:
```
dblink
hstore
postgis
tablefunc
unaccent
```

For the `install` process you need to input the dump file (sended by our support team), your customer id and the database credentials. For example:

```bash 
$ dbview install --config /tmp/dbview.toml /tmp/dbview_dump_customer_1329_20170131.pgbkp
  INFO[0000] Using config file: /tmp/dbview.toml
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

A full detail of options are available with the `--help` option. For example:

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

## Updating the replication


For the `replicate` process you need configure the `config.toml` adjuting the `[remote-database]` section then run:

```bash 
$ dbview replicate --config /tmp/dbview.toml
  INFO[0000] Using config file: /tmp/dbview.toml
  INFO[0000] Updating Replication Data...
  INFO[0001] Done.
```
> **IMPORTANT:** the user and password credentials are provided by our support team.

A full detail of options are available with the `--help` option. For example:

```bash
$ dbview replicate --help
Runs the replication functions and updates the target database at the latest version

Usage:
  dbview replicate [flags]

Flags:
      --daemon                            Run as daemon
  -l, --options.row_limit int32           row limit of each replication action (default 100)
      --refresh-interval duration         Refresh interval for daemon mode (default 30s)
      --remote-database.database string   Remote Database name (default "prod_umov_dbview")
      --remote-database.host string       Remote Database Host (default "dbview.umov.me")
      --remote-database.password string   Remote Database password
      --remote-database.port string       Remote Database Port (default "9999")
      --remote-database.ssl string        Remote SSL connection: 'require', 'verify-full', 'verify-ca', and 'disable' supported (default "disable")
      --remote-database.username string   Remote Database User (default "postgres")

Global Flags:
      --config string                    config file (default is $HOME/.dbview.yaml)
      --customer int                     Your customer ID
      --debug                            Show debug messages
      --help                             Show this help message
  -d, --local-database.database string   Local maintenance database. Used for administrative tasks. (default "postgres")
  -h, --local-database.host string       Local Database Host (default "127.0.0.1")
  -P, --local-database.password string   Local Database password
  -p, --local-database.port string       Local Database Port
      --local-database.ssl string        Local SSL connection: 'require', 'verify-full', 'verify-ca', and 'disable' supported (default "disable")
  -U, --local-database.username string   Local Database User (default "postgres")
      --pgsql-bin string                 PostgreSQL binaries PATH
```

### SystemD Service

With the `--daemon` option its possible start it as a systemD service. Create a file `/usr/lib/systemd/system/dbview.service` with the content bellow:

```
[Unit]
Description=DBView Replication
After=network.target

[Service]
Environment=CMD_OPTIONS=--daemon --duration 35s
Environment=CONFIG_FILE=/opt/dbview/config.toml
Type=simple
User=root
Group=root
ExecStart="/opt/dbview/dbview" replicate --config "${CONFIG_FILE}" ${CMD_OPTIONS}

[Install]
WantedBy=multi-user.target
```
> Update `CMD_DIR`, `CMD_OPTIONS`, and `CONFIG_FILE` for different PATH and configurations.

Then reload systemd configurations to use your service:

```
systemctl daemon-reload
```


## Getting support

If you need any help, contact us at [http://www.umov.me/suporte/](http://www.umov.me/suporte/).
