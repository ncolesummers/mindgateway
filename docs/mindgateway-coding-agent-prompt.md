# Coding Agent Prompt: Initialize MindGateway Project Structure

## Project Overview
Create the initial directory structure and boilerplate code for MindGateway, an enterprise LLM inference gateway using a hybrid microservices architecture with Go and gRPC.

## Architecture Context
- **Hybrid Architecture**: Two extracted microservices (Auth Service, Worker Registry) + Modular Monolith (Gateway)
- **Technology Stack**: Go 1.22+, gRPC, Gin web framework, PostgreSQL, Redis, etcd
- **Deployment**: Kubernetes-native with Docker support
- **Key Feature**: Workers connect outbound-only (no inbound ports required)

## Required Directory Structure

```
mindgateway/
├── cmd/
│   ├── gateway/
│   │   └── main.go                 # Main gateway application
│   ├── auth-service/
│   │   └── main.go                 # Auth microservice
│   └── worker-registry/
│       └── main.go                 # Worker registry microservice
├── internal/
│   ├── gateway/
│   │   ├── server/
│   │   │   ├── server.go          # HTTP server setup
│   │   │   └── middleware.go      # Common middleware
│   │   ├── handlers/
│   │   │   ├── chat.go           # OpenAI-compatible chat endpoint
│   │   │   ├── health.go         # Health check endpoints
│   │   │   └── metrics.go        # Prometheus metrics
│   │   ├── routing/
│   │   │   ├── router.go         # Request routing logic
│   │   │   ├── scorer.go         # Worker scoring algorithm
│   │   │   └── loadbalancer.go  # Load balancing strategies
│   │   └── queue/
│   │       ├── queue.go          # Priority queue implementation
│   │       └── scheduler.go      # Request scheduling
│   ├── auth/
│   │   ├── service.go            # Auth service implementation
│   │   ├── jwt.go                # JWT token handling
│   │   ├── oidc.go               # OIDC/Entra ID integration
│   │   └── rbac.go               # Role-based access control
│   ├── registry/
│   │   ├── service.go            # Worker registry service
│   │   ├── worker.go             # Worker state management
│   │   └── health.go             # Worker health checks
│   └── shared/
│       ├── config/
│       │   └── config.go         # Configuration management
│       ├── logging/
│       │   └── logger.go         # Structured logging
│       └── errors/
│           └── errors.go         # Custom error types
├── pkg/
│   ├── api/
│   │   ├── openai/
│   │   │   └── types.go         # OpenAI API types
│   │   └── ollama/
│   │       └── client.go        # Ollama client
│   └── proto/
│       ├── auth/
│       │   ├── auth.proto       # Auth service gRPC definitions
│       │   └── auth.pb.go       # Generated code
│       └── registry/
│           ├── registry.proto   # Registry service gRPC definitions
│           └── registry.pb.go   # Generated code
├── deployments/
│   ├── docker/
│   │   ├── gateway.Dockerfile
│   │   ├── auth.Dockerfile
│   │   └── registry.Dockerfile
│   ├── kubernetes/
│   │   ├── base/
│   │   │   ├── namespace.yaml
│   │   │   ├── configmap.yaml
│   │   │   └── secrets.yaml
│   │   ├── gateway/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   └── ingress.yaml
│   │   ├── auth/
│   │   │   ├── deployment.yaml
│   │   │   └── service.yaml
│   │   └── registry/
│   │       ├── deployment.yaml
│   │       └── service.yaml
│   └── helm/
│       └── mindgateway/
│           ├── Chart.yaml
│           ├── values.yaml
│           └── templates/
├── scripts/
│   ├── generate-proto.sh         # Generate gRPC code
│   ├── local-dev.sh             # Local development helper
│   └── test-integration.sh      # Run integration tests
├── test/
│   ├── integration/
│   │   ├── auth_test.go
│   │   ├── routing_test.go
│   │   └── e2e_test.go
│   └── load/
│       └── stress_test.go
├── configs/
│   ├── dev.yaml                 # Development config
│   ├── staging.yaml             # Staging config
│   └── prod.yaml                # Production config
├── docs/
│   ├── architecture/
│   │   ├── README.md
│   │   └── decisions/           # ADRs
│   ├── api/
│   │   └── openapi.yaml        # OpenAPI specification
│   └── deployment/
│       └── README.md
├── .github/
│   ├── workflows/
│   │   ├── ci.yml              # CI pipeline
│   │   ├── release.yml         # Release pipeline
│   │   └── security.yml        # Security scanning
│   └── PULL_REQUEST_TEMPLATE.md
├── docker-compose.yml           # Local development
├── docker-compose.test.yml      # Integration testing
├── Makefile                     # Build commands
├── go.mod
├── go.sum
├── .gitignore
├── .golangci.yml               # Linter configuration
└── README.md
```

## Initial File Contents

### 1. `go.mod`
```go
module github.com/company/mindgateway

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/spf13/viper v1.18.2
    github.com/prometheus/client_golang v1.18.0
    github.com/go-redis/redis/v9 v9.4.0
    google.golang.org/grpc v1.60.1
    google.golang.org/protobuf v1.32.0
    github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1
    github.com/sirupsen/logrus v1.9.3
    github.com/stretchr/testify v1.8.4
    go.etcd.io/etcd/client/v3 v3.5.11
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
)
```

### 2. `Makefile`
```makefile
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
```

### 3. `cmd/gateway/main.go`
```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/company/mindgateway/internal/gateway/server"
    "github.com/company/mindgateway/internal/shared/config"
    "github.com/company/mindgateway/internal/shared/logging"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize logger
    logger := logging.NewLogger(cfg.LogLevel)

    // Create server with modular components
    srv, err := server.New(
        server.WithConfig(cfg),
        server.WithLogger(logger),
    )
    if err != nil {
        logger.Fatalf("Failed to create server: %v", err)
    }

    // Start server
    go func() {
        if err := srv.Start(); err != nil {
            logger.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        logger.Errorf("Server forced to shutdown: %v", err)
    }
    
    logger.Info("Server exited")
}
```

### 4. `docker-compose.yml`
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: mindgateway
      POSTGRES_PASSWORD: localdev
      POSTGRES_DB: mindgateway
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  etcd:
    image: quay.io/coreos/etcd:v3.5.11
    command: >
      etcd
      --advertise-client-urls http://0.0.0.0:2379
      --listen-client-urls http://0.0.0.0:2379
    ports:
      - "2379:2379"

  # Mock Ollama for local development
  ollama-mock:
    image: nginx:alpine
    volumes:
      - ./test/mocks/ollama:/usr/share/nginx/html
    ports:
      - "11434:80"

volumes:
  postgres_data:
```

### 5. `internal/gateway/server/server.go`
```go
package server

import (
    "context"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/company/mindgateway/internal/gateway/handlers"
    "github.com/company/mindgateway/internal/shared/config"
    "github.com/sirupsen/logrus"
)

type Server struct {
    router *gin.Engine
    config *config.Config
    logger *logrus.Logger
    
    // Modular components
    authClient     AuthClient
    registryClient RegistryClient
    routingEngine  RoutingEngine
    queueManager   QueueManager
}

type Option func(*Server)

func New(opts ...Option) (*Server, error) {
    s := &Server{
        router: gin.New(),
    }
    
    for _, opt := range opts {
        opt(s)
    }
    
    s.setupRoutes()
    s.setupMiddleware()
    
    return s, nil
}

func (s *Server) setupRoutes() {
    // Health checks
    s.router.GET("/health", handlers.Health())
    s.router.GET("/ready", handlers.Ready())
    
    // Metrics
    s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    // API routes
    v1 := s.router.Group("/v1")
    {
        // OpenAI compatible endpoints
        v1.POST("/chat/completions", s.handleChatCompletion)
        v1.POST("/completions", s.handleCompletion)
        v1.POST("/embeddings", s.handleEmbeddings)
    }
    
    // Admin routes
    admin := s.router.Group("/admin")
    admin.Use(s.requireAdmin())
    {
        admin.GET("/workers", s.listWorkers)
        admin.GET("/queue", s.queueStatus)
    }
}

func (s *Server) Start() error {
    return s.router.Run(s.config.Server.Address)
}

func (s *Server) Shutdown(ctx context.Context) error {
    // Graceful shutdown logic
    return nil
}
```

### 6. `.gitignore`
```
# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of go coverage tool
*.out

# Go vendor directory
vendor/

# IDE files
.idea/
.vscode/
*.swp
*.swo

# Environment files
.env
.env.local

# Generated files
*.pb.go
*.pb.gw.go

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Local data
data/
```

## Additional Requirements

1. **Create Initial Proto Files**: 
   - Define gRPC services for Auth and Registry
   - Include streaming for worker connections

2. **Setup GitHub Actions CI**:
   - Test, lint, and build on every push
   - Security scanning with gosec
   - Docker image building

3. **Include README.md** with:
   - Architecture overview
   - Quick start guide
   - Development setup instructions
   - API documentation links

4. **Create docker-compose.test.yml** for integration testing environment

5. **Add example configuration files** in `configs/` directory

## Development Environment Requirements
- Go 1.22+
- Docker and Docker Compose
- Make
- protoc (for gRPC generation)
- golangci-lint

Please create this complete directory structure with all the specified files and their initial content. Focus on making it immediately runnable with `make run-local` after setup.