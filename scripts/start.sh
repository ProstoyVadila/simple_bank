#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the web server"
exec "$@"
