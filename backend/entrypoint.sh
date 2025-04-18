#!/bin/sh
set -e

echo "💤 Waiting for the database to become available..."
until pg_isready -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -d ${POSTGRES_DB_NAME}; do
  echo "The database is unavailable, we're waiting..."
  sleep 2
done


echo "📦 Performing migrations..."
goose -dir ./migrations postgres ${DATABASE_URL} up

echo "🚀 Launch the application..."
./time-capsule-memories

echo "🛑 The application has terminated"
