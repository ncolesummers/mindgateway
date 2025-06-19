# MindGateway

MindGateway is an enterprise LLM inference gateway that provides a unified interface for interacting with various LLM backends through a hybrid microservices architecture.

## Architecture Overview

MindGateway uses a hybrid architecture combining a modular monolith with extracted microservices:

- **Gateway (Modular Monolith)**: Handles incoming requests, authentication, routing, and communication with workers
- **Auth Service (Microservice)**: Manages authentication and authorization
- **Worker Registry (Microservice)**: Handles worker registration, discovery, and health monitoring

Key features:

- **OpenAI-compatible API**: Drop-in replacement for OpenAI services
- **Outbound-only worker connections**: Workers connect to the gateway, eliminating the need for inbound ports
- **Flexible routing**: Intelligent request routing based on model availability, load, and priorities
- **Enterprise security**: Role-based access control, token validation, and comprehensive audit logging

![Architecture Diagram](./docs/architecture/architecture.png)

## Getting Started

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- Make
- protoc (for gRPC code generation)
- golangci-lint

### Quick Start

1. Clone the repository:

```bash
git clone https://github.com/ncolesummers/mindgateway.git
cd mindgateway
```

2. Set up the local development environment:

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Generate protobuf code
./scripts/generate-proto.sh

# Start local development environment
./scripts/local-dev.sh
```

3. Access the API at http://localhost:8080

## Development

### Project Structure

```
mindgateway/
├── cmd/                    # Application entry points
├── internal/               # Internal packages
│   ├── gateway/            # Gateway implementation
│   ├── auth/               # Auth service implementation
│   ├── registry/           # Worker registry implementation
│   └── shared/             # Shared utilities
├── pkg/                    # Public packages
│   ├── api/                # API types and clients
│   └── proto/              # Protocol buffer definitions
├── deployments/            # Deployment configurations
├── scripts/                # Utility scripts
├── test/                   # Test files
└── configs/                # Configuration files
```

### Building

Build all services:

```bash
make build
```

### Testing

Run unit tests:

```bash
make test
```

Run integration tests:

```bash
make test-integration
```

## Deployment

### Docker

Build Docker images:

```bash
docker build -t mindgateway/gateway -f deployments/docker/gateway.Dockerfile .
docker build -t mindgateway/auth -f deployments/docker/auth.Dockerfile .
docker build -t mindgateway/registry -f deployments/docker/registry.Dockerfile .
```

### Kubernetes

Deploy to Kubernetes:

```bash
# Create namespace
kubectl apply -f deployments/kubernetes/base/namespace.yaml

# Deploy secrets and configmaps
kubectl apply -f deployments/kubernetes/base/

# Deploy services
kubectl apply -f deployments/kubernetes/gateway/
kubectl apply -f deployments/kubernetes/auth/
kubectl apply -f deployments/kubernetes/registry/
```

## API Documentation

MindGateway provides an OpenAI-compatible API. The API documentation is available at:

- `/docs/api` - OpenAPI specification
- `/docs/swagger` - Swagger UI (when running in development mode)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
