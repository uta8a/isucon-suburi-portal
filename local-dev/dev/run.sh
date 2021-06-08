#!/bin/sh
set -ue

export $(cat ./local-dev/dev/.env)
timeout 30 sh -c "until nc -vz $DB_HOST $DB_PORT; do sleep 1; done" && sql-migrate up

go run main.go
