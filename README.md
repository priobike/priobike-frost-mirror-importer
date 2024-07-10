# priobike-frost-mirror-importer

Service that imports all things currently marked as "good" for the Frost Server. As "good" marked things are maintained [here](https://daten-hamburg.de/tlf_public/).

PrioBike-Services can then use this information to only request those "good" things.

[Learn more about PrioBike](https://github.com/priobike)

## Quickstart

The easiest way to run the frost mirror importer is to use the contained `docker-compose`:
```
docker-compose up
```

These are the build arguments we use to configure the frost server build:
- `serviceRootUrl`
- `http_cors_enable`
- `http_cors_allowed.origins`
- `persistence_db_driver`
- `persistence_db_url`
- `persistence_db_username`
- `persistence_db_password`
- `persistence_autoUpdateDatabase`
- `persistence_idGenerationMode`

For more information see: [Frost Server Settings](https://fraunhoferiosb.github.io/FROST-Server/settings/settings.html)

## API and CLI

Frost server running under:
`localhost:8080/`

Can be tested with Frost API requests https://fraunhoferiosb.github.io/FROST-Server/sensorthingsapi/1_Home.html.

For example:
`http://localhost:8080/FROST-Server/v1.1/Things`

## Contributing

We highly encourage you to open an issue or a pull request. You can also use our repository freely with the `MIT` license.

Every service runs through testing before it is deployed in our release setup. Read more in our [PrioBike deployment readme](https://github.com/priobike/.github/blob/main/wiki/deployment.md) to understand how specific branches/tags are deployed.

## Anything unclear?

Help us improve this documentation. If you have any problems or unclarities, feel free to open an issue.
