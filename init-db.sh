#!/bin/bash

# MUST BE EXECUTED AS postgres USER

# See: https://trac.osgeo.org/postgis/wiki/UsersWikiPostGIS3UbuntuPGSQLApt
pg_lsclusters 14 main start
service postgresql restart

# Wait for the postgres database to startup
echo "Waiting for postgres server..."
# Await PostGreSQL server to become available
RETRIES=20
while [ "$RETRIES" -gt 0 ]
do
  echo "Waiting for postgres: pg_isready -d ${POSTGRES_NAME} -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER}"
  PG_STATUS="$(pg_isready -d ${POSTGRES_NAME} -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U ${POSTGRES_USER})"
  PG_EXIT=$(echo $?)
  if [ "$PG_EXIT" = "0" ];
    then
      RETRIES=0
  fi
  sleep 0.5
done
echo "Postgres server is up!"

# Create the database
echo "Creating database ${POSTGRES_DB}..."
createdb ${POSTGRES_DB} -O ${POSTGRES_USER}

# Load PostGIS into $POSTGRES_DB
for DB in "$POSTGRES_DB" "${@}"; do
    psql --dbname="$DB" -c "
        CREATE EXTENSION IF NOT EXISTS postgis;
    "
done

# Set the password for the postgres user
echo "Setting password for postgres user..."
psql -c "ALTER USER ${POSTGRES_USER} WITH PASSWORD '${POSTGRES_PASSWORD}';"
