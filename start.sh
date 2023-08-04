#!/bin/sh

set -e

# echo "starting database migration..."
# /app/migrate -path /app/migration -database "$ENV_DATABASE_SOURCE" -verbose up
# echo "database migration completed"

echo "starting the app..."
exec "$@"