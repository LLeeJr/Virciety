# Message Service
This service is currently only an example. However, it introduces multiple core concepts of go-kit, also insuring a well-built structure.
In order to run the [main.go](main.go) file sync the service's dependencies and run ``run go main.go`` from the root.
Additionally, there is a [docker-compose.yml](docker-compose.yml) which launches on execution a postgres database and an adminer on port 8080.
The service connects to the database all by itself, however, data-storage happens only locally at the moment.