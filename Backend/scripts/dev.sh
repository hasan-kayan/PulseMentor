#!/bin/bash

# Development run script for PulseMentor Backend
# Usage: ./scripts/dev.sh

set -e

cd "$(dirname "$0")/.."

# Set default environment variables if not set
export APP_ENV="${APP_ENV:-dev}"
export SERVER_HOST="${SERVER_HOST:-0.0.0.0}"
export SERVER_PORT="${SERVER_PORT:-8080}"
export DATABASE_URL="${DATABASE_URL:-postgres://localhost/pulsementor?sslmode=disable}"
export JWT_SECRET="${JWT_SECRET:-dev-secret-key-change-in-production-min-32-chars-long}"
export JWT_ISSUER="${JWT_ISSUER:-pulsementor}"
export ACCESS_TOKEN_TTL="${ACCESS_TOKEN_TTL:-24h}"
export REFRESH_TOKEN_TTL="${REFRESH_TOKEN_TTL:-168h}"
export BCRYPT_COST="${BCRYPT_COST:-12}"

echo "=========================================="
echo "PulseMentor Backend - Development Server"
echo "=========================================="
echo ""
echo "Environment: $APP_ENV"
echo "Server: http://$SERVER_HOST:$SERVER_PORT"
echo "Database: $DATABASE_URL"
echo ""
echo "Starting server..."
echo ""

go run cmd/server/main.go

