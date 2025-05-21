#!/bin/sh

# Load environment variables from app.env
set -a
. ./app.env
set +a

# Extract host and port from DB_SOURCE
DB_HOST=$(echo $DB_SOURCE | sed -E 's/.*@([^:]+):([0-9]+).*/\1/')
DB_PORT=$(echo $DB_SOURCE | sed -E 's/.*@([^:]+):([0-9]+).*/\2/')

echo "Waiting for database at $DB_HOST:$DB_PORT..."

# Wait for DB to be ready
until nc -z "$DB_HOST" "$DB_PORT"; do
  echo "Waiting for database at $DB_HOST:$DB_PORT..."
  sleep 2
done

echo "Database is up! Running migrations and starting server..."

make migrateup || echo "Migration step skipped or failed"

./server
