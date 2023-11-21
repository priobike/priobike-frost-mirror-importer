FROM fraunhoferiosb/frost-server-http:latest

ENV POSTGRES_USER=postgres
ENV POSTGRES_NAME=frost-db
ENV POSTGRES_PASSWORD=frost-password
ENV POSTGRES_DB=frost-db
ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432

ENV FROST_SERVICE_ROOT_URL=http://localhost:8080/FROST-Server
ENV FROST_HTTP_CORS_ENABLE=true
ENV FROST_HTTP_CORS_ALLOWED_ORIGINS=*
ENV FROST_PERSISTENCE_DB_DRIVER=org.postgresql.Driver
ENV FROST_PERSISTENCE_DB_URL=jdbc:postgresql://localhost:5432/frost-db
ENV FROST_PERSISTENCE_DB_USERNAME=postgres
ENV FROST_PERSISTENCE_DB_PASSWORD=frost-password
ENV FROST_PERSISTENCE_AUTO_UPDATE_DATABASE=true

# Change user to root
USER root

# Must happen before installing PostGIS
RUN useradd -ms /bin/bash postgres

# Install PostGIS
# See: https://trac.osgeo.org/postgis/wiki/UsersWikiPostGIS3UbuntuPGSQLApt
RUN apt update -y
RUN apt install -y postgresql-14-postgis-3 postgresql-client

COPY --chown=postgres:postgres init-db.sh /postgres/init-db.sh

# TODO: Remove this mock
COPY demoEntities.json /root/demoEntities.json

COPY preheat.sh /root/preheat.sh
COPY run.sh /root/run.sh

# Run preheating
RUN /root/preheat.sh

CMD ["/root/run.sh"]