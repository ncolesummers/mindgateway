#!/bin/bash
set -e

# This script generates Go code from proto files
echo "Generating protobuf code..."

# Create output directories if they don't exist
mkdir -p pkg/proto/auth
mkdir -p pkg/proto/registry

# Generate auth service proto
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  pkg/proto/auth/auth.proto

# Generate registry service proto
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  pkg/proto/registry/registry.proto

echo "Protobuf code generation complete!"
