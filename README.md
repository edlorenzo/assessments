# assessments

> A service for Philadelphia's bike-sharing program.

## Prerequisites

Download and install the following:

- Go 1.21 - https://golang.org/
- GORM - https://gorm.io/
- go-migrate - [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Setup

Go recommends having a single workspace directory for all your Go projects. Before cloning, [set up your Go workspace first](https://golang.org/doc/code.html).

Create the project directory inside the workspace:

```zsh
$ git clone git@github.com:edlorenzo/assessments.git
```

Declare environment variables

- Refer to the [.env.sample](https://github.com/edlorenzo/assessments/blob/main/.env.example).
- Change the values accordingly
- Rename the .env.sample to .env

```
APP_HOST=127.0.0.1
APP_PORT=3000
APP_NAME=assessments-service
APP_ENVIRONMENT=dev
APP_READ_TIMEOUT=60
APP_WRITE_TIMEOUT=60

PROCESS=WORKER

DB_HOST=localhost
DB_DATABASE=assessment
DB_USERNAME=assessment
DB_PASSWORD=P@ssw0rd
DB_PORT=5433

MIGRATION_PATH=D:/_0.0.7-Golang/2024/technical-exam/assessments/cmd/migrations

OPEN_WEATHER_MAP_API_KEY=008d34dc370d3d63b2bd19e1dd13620e
OPEN_WEATHER_MAP_API_URL=http://api.openweathermap.org/data/2.5/weather
OPEN_WEATHER_MAP_API_CITY_PARAM=Philadelphia

INDEGO_API_URL=https://bts-status.bicycletransit.workers.dev/
INDEGO_API_ABBREVIATION_PARAM=phl

JOB_DELAY_MIN=60
```

- Note: make sure the `MIGRATION_PATH, and OPEN_WEATHER_MAP_API_KEY` is updated.

To start the server locally:

```zsh
$ cd assessments
$ go mod tidy
$ Update the .env file and set the value for PROCESS to SERVER. Ex: PROCESS=SERVER
$ go run cmd/main.go or make run 
```

To start Job locally:

```zsh
$ cd assessments
$ go mod tidy
$ Update the .env file and set the value for PROCESS to WORKER. Ex: PROCESS=WORKER
$ As well the JOB_DELAY_MIN=60 is configurable per minutes. The default value is 1hr.
$ go run cmd/main.go or make run 
```

Start the server via docker-compose:

```
docker-compose up -d --force-recreate
```

## Healthcheck

- Access the healthcheck endpoint (http://localhost:8003/health/liveness | http://localhost:8003/health/readiness) to see if the service is up.

## Postman Collection

- Import the [collection](https://www.postman.com/dark-desert-384866/workspace/new-team-workspace/collection/2409862-4c9f8d54-396c-462f-8b0e-0bf22522c580) to Postman to see the sample payload and response.
