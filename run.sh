#!/bin/bash

# MUST BE EXECUTED AS root USER

# Init database in background as postgres user
su - postgres -c "POSTGRES_USER=${POSTGRES_USER} POSTGRES_PASSWORD=${POSTGRES_PASSWORD} POSTGRES_DB=${POSTGRES_DB} POSTGRES_NAME=${POSTGRES_NAME} POSTGRES_HOST=${POSTGRES_HOST} POSTGRES_PORT=${POSTGRES_PORT} /postgres/init-db.sh"

# Run the server
catalina.sh run
