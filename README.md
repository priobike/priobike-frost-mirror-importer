# PrioBike-Frost-Mirror-Importer
Service that imports all things currently marked as "good" for the Frost Server. PrioBike-Services can then use this information to only request those "good" things.

# Test locally
`docker-compose up`

Frost server running under:
`localhost:8080/`

Can be tested with Frost API requests https://fraunhoferiosb.github.io/FROST-Server/sensorthingsapi/1_Home.html.

For example:
`http://localhost:8080/FROST-Server/v1.1/Things`
