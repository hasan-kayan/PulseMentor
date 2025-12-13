#!/bin/bash

# Migration script for PulseMentor Backend
# Usage: ./scripts/migrate.sh

set -e

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    echo "Example: export DATABASE_URL='postgres://user:password@localhost:5432/pulsementor?sslmode=disable'"
    exit 1
fi

MIGRATIONS_DIR="migrations"
MIGRATION_FILES=$(ls -1 $MIGRATIONS_DIR/*.sql 2>/dev/null | sort)

if [ -z "$MIGRATION_FILES" ]; then
    echo "No migration files found in $MIGRATIONS_DIR"
    exit 1
fi

echo "Running migrations..."

for migration in $MIGRATION_FILES; do
    echo "Executing: $migration"
    psql "$DATABASE_URL" -f "$migration" || {
        echo "Error executing $migration"
        exit 1
    }
done

echo "Migrations completed successfully!"

