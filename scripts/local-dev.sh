#!/bin/bash
set -e

# This script helps with local development setup
echo "Setting up local development environment..."

# Start dependencies with docker-compose
docker-compose up -d postgres redis etcd

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 3

# Run migrations if needed
echo "Running database migrations..."
# TODO: Add migration command here

# Run the gateway in dev mode
echo "Starting gateway in development mode..."
CONFIG_PATH=configs ENVIRONMENT=dev go run cmd/gateway/main.go
