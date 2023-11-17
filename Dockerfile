FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o main .

ENV POSTGRES_DB=sensorthings
ENV POSTGRES_USER=sensorthings
ENV POSTGRES_PASSWORD=ChangeMe

FROM postgis/postgis:14-3.2 AS runner

COPY --from=builder /app/main .

# Install OpenJDK-8
RUN apt-get update && \
    apt-get install -y default-jre && \
    apt-get install -y wget

# Setup JAVA_HOME -- useful for docker commandline
# ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
# RUN export JAVA_HOME

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

RUN wget https://repo1.maven.org/maven2/de/fraunhofer/iosb/ilt/FROST-Server/FROST-Server.HTTP/2.0.10/FROST-Server.HTTP-2.0.10-classes.jar
RUN wget https://repo1.maven.org/maven2/de/fraunhofer/iosb/ilt/FROST-Server/FROST-Server.HTTP/2.0.10/FROST-Server.HTTP-2.0.10-javadoc.jar
RUN wget https://repo1.maven.org/maven2/de/fraunhofer/iosb/ilt/FROST-Server/FROST-Server.HTTP/2.0.10/FROST-Server.HTTP-2.0.10-sources.jar

RUN java -cp FROST-Server.HTTP-2.0.10-sources.jar de.fraunhofer.iosb.ilt.FROST-Server

 # TODO start frost server and import data