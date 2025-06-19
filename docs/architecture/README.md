# MindGateway Architecture

## Overview

MindGateway uses a hybrid architecture combining a modular monolith with extracted microservices:

- **Gateway (Modular Monolith)**: Handles incoming requests, authentication, routing, and communication with workers
- **Auth Service (Microservice)**: Manages authentication and authorization
- **Worker Registry (Microservice)**: Handles worker registration, discovery, and health monitoring

## Components

### Gateway

The Gateway is the main entry point for client requests. It provides:

- OpenAI-compatible API endpoints
- Request validation and rate limiting
- Authentication and authorization (via Auth Service)
- Request routing to appropriate workers
- Queue management for requests
- Response streaming and handling

### Auth Service

The Auth Service is responsible for:

- User authentication
- JWT token generation and validation
- Role-based access control
- Permission management
- User and API key management

### Worker Registry

The Worker Registry handles:

- Worker registration and discovery
- Worker health monitoring
- Worker capability tracking
- Load balancing information
- Worker connection management

## Communication

- Gateway to Auth Service: gRPC
- Gateway to Worker Registry: gRPC
- Gateway to Workers: Bidirectional streaming gRPC
- Clients to Gateway: HTTP/REST

## Key Features

- **Outbound-only worker connections**: Workers connect to the gateway, eliminating the need for inbound ports
- **Flexible routing**: Intelligent request routing based on model availability, load, and priorities
- **Enterprise security**: Role-based access control, token validation, and comprehensive audit logging
- **Horizontal scaling**: All components can be independently scaled
- **Observability**: Comprehensive metrics, logging, and tracing