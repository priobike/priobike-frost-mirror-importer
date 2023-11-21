FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main .

ARG POSTGRES_DB=sensorthings
ARG POSTGRES_USER=sensorthings
ARG POSTGRES_PASSWORD=ChangeMe

ENV POSTGRES_DB=${POSTGRES_DB}
ENV POSTGRES_USER={POSTGRES_USER}
ENV POSTGRES_PASSWORD={POSTGRES_PASSWORD}

ENV serviceRootUrl=http://localhost:8080/FROST-Server
ENV plugins_multiDatastream.enable=false
ENV http_cors_enable=true
ENV http_cors_allowed_origins=*
ENV persistence_db_driver=org.postgresql.Driver
ENV persistence_db_url=jdbc:postgresql://localhost:5432/sensorthings
ENV persistence_db_username=sensorthings
ENV persistence_db_password=ChangeMe
ENV persistence_autoUpdateDatabase=true
ENV persistence_idGenerationMode=ServerAndClientGenerated
 
ENV defaultCount=false
ENV defaultTop=100
 
ENV useAbsoluteNavigationLinks=true
ENV countMode=FULL
ENV extension_customLinks=true
ENV extension_customLinks_recurseDepth=0
ENV extension_filterDelete_enable=false
ENV plugins_odata_enable=false
ENV plugins_openApi_enable=true
ENV resources_requests_cpu=300m
ENV resources_requests_memory=900Mi
ENV resources_limits_cpu=300m
ENV resources_limits_memory=900Mi
 
ENV db_autoUpdate=true
ENV alwaysOrderbyId=false
ENV db_maximumConnection=10
ENV db_maximumIdleConnection=10
ENV db_minimumIdleConnection=10

FROM tomcat:9-jdk17 AS runner

RUN apt-get update && apt-get install -y \
  binutils \
  libproj-dev \
  gdal-bin

# Install Postgres client to check liveness of the database
RUN apt-get install -y postgresql-client

# Install postgres dev utils
RUN apt-get install -y libpq-dev

# Install postgis
RUN apt-get install -y postgresql-14-postgis-3

RUN apt-get -y install systemctl

RUN apt-get install -y sudo

RUN service postgresql start

RUN echo "listen_addresses = '*'" >> /etc/postgresql/14/main/postgresql.conf

# Connect to the default PostgreSQL template database
RUN sudo -u postgres psql template1

RUN ALTER USER {POSTGRES_USER} with password {POSTGRES_PASSWORD};

RUN exit

RUN systemctl restart postgresql.service

# Install CURL for healthcheck
RUN apt-get install -y curl

# Copy go script
COPY --from=builder /app/main .

# Install wget
RUN apt-get install -y wget

# Setup JAVA_HOME -- useful for docker commandline
# ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
# RUN export JAVA_HOME

RUN wget https://repo1.maven.org/maven2/de/fraunhofer/iosb/ilt/FROST-Server/FROST-Server.HTTP/2.0.10/FROST-Server.HTTP-2.0.10.war

RUN apt-get update && apt-get install unzip && apt-get clean

# Copy to images tomcat path
RUN unzip -d ${CATALINA_HOME}/webapps/FROST-Server FROST-Server.HTTP-2.0.10.war \
    && addgroup --system --gid 1000 tomcat \
    && adduser --system --uid 1000 --gid 1000 tomcat \
    && chgrp -R 0 $CATALINA_HOME \
    && chmod -R g=u $CATALINA_HOME

USER tomcat

# Expose the default Tomcat port
EXPOSE 8080

# Start Tomcat when the container starts
CMD ["catalina.sh", "run"]

# TODO import data