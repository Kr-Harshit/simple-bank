#!/bin/sh

set -e

echo "load env variables"
source /app/app.env

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"