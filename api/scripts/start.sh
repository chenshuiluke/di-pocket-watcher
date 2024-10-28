#!/bin/bash
# scripts/start.sh

echo "Waiting for database..."
until pg_isready -h db -U postgres -d di_pocket_watcher > /dev/null 2>&1; do
    sleep 1
done
echo "Database is ready!"

echo "Running database migrations..."
goose -dir ./db/migrations up

echo "Starting application..."
air -c .air.toml