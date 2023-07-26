#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$ENV_DATABASE_SOURCE" -verbose up

echo "start the app"
exec "$@"