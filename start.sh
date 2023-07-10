#!/bin/sh

set -e

echo "Run database migration"
echo "$DB_SOURCE"
migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"