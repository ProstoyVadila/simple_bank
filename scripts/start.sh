#!/bin/sh

set -e

echo "run db migrations"
source /app/prod.env
echo "$(cat /app/prod.env)"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the web server"
exec "$@"
