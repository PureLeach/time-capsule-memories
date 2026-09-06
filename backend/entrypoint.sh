#!/bin/sh
set -eu

echo "Waiting for the database to become available..."
until pg_isready -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT}" -d "${POSTGRES_DB_NAME}" >/dev/null 2>&1; do
  sleep 2
done

echo "Applying migrations..."
goose -dir ./migrations postgres "${DATABASE_URL}" up

echo "Starting application..."
exec ./time-capsule-memories
