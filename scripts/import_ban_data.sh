#!/bin/bash

export PGPASSWORD="adr"

psql -U adr -c "DROP SCHEMA IF EXISTS ban_ign CASCADE;"

psql -U adr < /tmp/scripts/create_schemas_tables.sql

for f in /tmp/ban-ign-id/*.csv.gz; do
    echo "[ban] integrating $f"
    psql -U adr -c "COPY ban_ign.ban FROM PROGRAM 'gzip -dc $f' DELIMITER ';' CSV HEADER"
done

psql -U adr -c "CREATE EXTENSION IF NOT EXISTS postgis"
echo "[ban] Add geometry column"
psql -U adr -c "ALTER TABLE ban_ign.ban ADD COLUMN geom geometry(Point,4326)"
echo "[ban] Fill geometry column"
psql -U adr -c "UPDATE ban_ign.ban SET geom=ST_SetSRID(ST_MakePoint(lon,lat),4326)"
echo "[ban] Remove unnecessary columns"
psql -U adr -c "ALTER TABLE ban_ign.ban DROP COLUMN x, DROP COLUMN y, DROP COLUMN lon, DROP COLUMN lat"
echo "[ban] Create lowercase index on cle_interop"
psql -U adr -c "CREATE INDEX cle_interop_lower_idx ON ban_ign.ban (lower(cle_interop))"
echo "[ban] Create index on id_ban_adresse"
echo "VACUUM ANALYZE"
psql -U adr -c "VACUUM FULL ANALYZE"

# Optional for now
# echo "[ban] Create geographic index"
# psql -U adr -c "CREATE INDEX ban_idx ON ban_ign.ban USING GIST(geom)"
# echo "[ban] Add primary key on ban"
# psql -U adr -c "ALTER TABLE ban_ign.ban ADD PRIMARY KEY(id_ban_position)"