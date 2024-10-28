#!/bin/bash
# scripts/start.sh

# Function to test database connection
wait_for_db() {
    echo "Waiting for database to be ready..."
    while ! goose -dir ./db/migrations postgres "host=db user=postgres password=di_pocket_watcher dbname=di_pocket_watcher port=5432 sslmode=disable" status >/dev/null 2>&1; do
        echo "Database not ready... waiting"
        sleep 2
    done
    echo "Database is ready!"
}

# Wait for database
wait_for_db

# Run migrations
echo "Running database migrations..."
goose -dir ./db/migrations postgres "host=db user=postgres password=di_pocket_watcher dbname=di_pocket_watcher port=5432 sslmode=disable" up

# Start the application with Air
echo "Starting application..."
air -c .air.toml