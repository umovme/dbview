DBVIEW in docker Steps:

1. Adjust config files in config folder (config.toml and .env) with your infos.

2. Request to uMov.Me a dump file and save it into config folder with "dbview.pgbkp" name

3. To start database container run:
    docker-compose up -d db

4. To restore dump and initialize dbview run:
    docker-compose up dbview-install

5. To run replicate run:
    docker-compose up -d dbview-replicate
    