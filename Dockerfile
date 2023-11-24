FROM fraunhoferiosb/frost-server-http:latest

# Env for database
ENV POSTGRES_USER=postgres
ENV POSTGRES_NAME=frost-db
ENV POSTGRES_PASSWORD=frost-password
ENV POSTGRES_DB=frost-db
ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432

# Env for the FROST server
# During the preheating, the FROST server is only running on localhost
ENV FROST_SERVICE_ROOT_URL=http://localhost:8080/FROST-Server
ENV FROST_HTTP_CORS_ENABLE=true
ENV FROST_HTTP_CORS_ALLOWED_ORIGINS=*
ENV FROST_PERSISTENCE_DB_DRIVER=org.postgresql.Driver
ENV FROST_PERSISTENCE_DB_URL=jdbc:postgresql://localhost:5432/frost-db
ENV FROST_PERSISTENCE_DB_USERNAME=postgres
ENV FROST_PERSISTENCE_DB_PASSWORD=frost-password
ENV FROST_PERSISTENCE_AUTO_UPDATE_DATABASE=true
ENV FROST_PERSISTENCE_IDGENERATIONMODE=ServerAndClientGenerated

# Env for the sync script
ENV FROST_HAMBURG_URL=https://tld.iot.hamburg.de/v1.1/
# During the preheating, the FROST server is only running on localhost
ENV FROST_PROXY_URL=http://localhost:8080/FROST-Server/v1.1/
ENV EXCLUDE_LIST_FILE=/root/exclude_list.xlsx
ENV VR_LIST_FILE=/root/vr_list.xlsx

# Change user to root
USER root

# Must happen before installing PostGIS
RUN useradd -ms /bin/bash postgres

# Install PostGIS
# See: https://trac.osgeo.org/postgis/wiki/UsersWikiPostGIS3UbuntuPGSQLApt
RUN apt update -y
RUN apt install -y postgresql-14-postgis-3 postgresql-client

COPY --chown=postgres:postgres init-db.sh /postgres/init-db.sh

COPY preheat.sh /root/preheat.sh
COPY run.sh /root/run.sh

# Install python3 and pip and the required packages
RUN apt-get install -y python3 python3-pip
COPY requirements.txt /root/requirements.txt
RUN python3 -m pip install -r /root/requirements.txt
COPY sync.py /root/sync.py
COPY exclude_list.xlsx /root/exclude_list.xlsx
COPY vr_list.xlsx /root/vr_list.xlsx

# Run preheating
RUN /root/preheat.sh

CMD ["/root/run.sh"]