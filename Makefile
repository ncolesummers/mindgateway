.PHONY: build test lint proto run-local

# Build all services
build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/auth-service ./cmd/auth-service
	go build -o bin/worker-registry ./cmd/worker-registry

# Run tests
test:
	go test -v -race ./...

# Run linter
lint:
	golangci-lint run

# Generate protobuf
proto:
	./scripts/generate-proto.sh

# Run local development
run-local:
	docker-compose up -d
	go run ./cmd/gateway

# Run integration tests
test-integration:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit