ARG postgres_version=13
FROM postgres:${postgres_version}

ARG postgis_version=3

RUN apt-get update \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${postgis_version} \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${postgis_version}-scripts \
 && apt-get install -y postgresql-${PG_MAJOR}-pgrouting \
 && apt-get install -y --no-install-recommends postgis \
 && rm -rf /var/lib/apt/lists/*

ADD conf.sql /docker-entrypoint-initdb.d/conf.sql
