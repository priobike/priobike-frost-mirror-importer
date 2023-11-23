#!/bin/bash

# MUST BE EXECUTED AS root USER

# Init database in background as postgres user
su - postgres -c "POSTGRES_USER=${POSTGRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} POSTGRES_DB=${POSTGRES_DB} POSTGRES_NAME=${POSTGRES_NAME} POSTGRES_HOST=${POSTGRES_HOST} POSTGRES_PORT=${POSTGRES_PORT} /postgres/init-db.sh"

# Run FROST-Server in background as root user
serviceRootUrl=${FROST_SERVICE_ROOT_URL} \
http_cors_enable=${FROST_HTTP_CORS_ENABLE} \
http_cors_allowed_origins=${FROST_HTTP_CORS_ALLOWED_ORIGINS} \
persistence_db_driver=${FROST_PERSISTENCE_DB_DRIVER} \
persistence_db_url=${FROST_PERSISTENCE_DB_URL} \
persistence_db_username=${FROST_PERSISTENCE_DB_USERNAME} \
persistence_db_password=${FROST_PERSISTENCE_DB_PASSWORD} \
persistence_autoUpdateDatabase=${FROST_PERSISTENCE_AUTO_UPDATE_DATABASE} \
persistence_idGenerationMode=${FROST_PERSISTENCE_IDGENERATIONMODE} \
catalina.sh run &

# TODO: Properly wait for the FROST server to start, but 10 seconds is enough for now
sleep 10

# Download, filter and push things from remote FROST server to local FROST server
python3 /root/sync.py
