#!/bin/bash
set -e

# This script runs integration tests
echo "Running integration tests..."

# Start test dependencies
docker-compose -f docker-compose.test.yml up -d postgres-test redis-test etcd-test

# Wait for services to be ready
echo "Waiting for test services to be ready..."
sleep 3

# Run integration tests
TEST_POSTGRES_URI="postgres://mindgateway_test:test@localhost:5433/mindgateway_test?sslmode=disable" \
TEST_REDIS_URI="redis://localhost:6380/0" \
TEST_ETCD_ENDPOINTS="localhost:2380" \
go test -v ./test/integration/...

# Clean up
echo "Cleaning up test environment..."
docker-compose -f docker-compose.test.yml down

echo "Integration tests completed!"
