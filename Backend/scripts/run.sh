#!/bin/bash

# Run script for PulseMentor Backend
# Usage: ./scripts/run.sh

set -e

cd "$(dirname "$0")/.."

# Check if required environment variables are set
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    echo "Example: export DATABASE_URL='postgres://user:password@localhost:5432/pulsementor?sslmode=disable'"
    exit 1
fi

if [ -z "$JWT_SECRET" ]; then
    echo "Error: JWT_SECRET environment variable is not set"
    echo "Example: export JWT_SECRET='your-super-secret-jwt-key'"
    exit 1
fi

echo "Starting PulseMentor Backend..."
echo "Server will be available at: http://${SERVER_HOST:-0.0.0.0}:${SERVER_PORT:-8080}"

go run cmd/server/main.go

