# SimpleBank

This repository contains the codes of the Backend Master Class course by TECH SCHOOL.

## Setup local environment

### Install Tools:

- Golang
- Postgres
- Docker
- [Golang Migrate](https://github.com/golang-migrate/migrate)
- [SQLC](https://github.com/kyleconroy/sqlc)

### Setup Services

- Start all basic services: `$ make services `
- Run db migration up all versions: `$ make migrateup`

## Developer commands

| Command            | Description                                            |
|--------------------|--------------------------------------------------------|
| `make services`    | Run the required services using Docker-compose         |
| `make test`        | Run all the defined tests                              |
| `make createdb`    | Create the simple_bank db in the docker postgres       |
| `make dropdb`      | Drop the development db                                |
| `make migrateup`   | Run db migration to the latest migration               |
| `make migratedown` | Run db migration one version down                      |
| `make sqlc`        | Run SQLC generate to create crud after add a new query |
