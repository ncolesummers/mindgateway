FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/auth-service

# Create final minimal image
FROM alpine:3.18

WORKDIR /app

# Copy binary and configs
COPY --from=builder /bin/auth-service /app/auth-service
COPY --from=builder /app/configs /app/configs

# Add necessary certificates and timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Default environment
ENV CONFIG_PATH=/app/configs
ENV ENVIRONMENT=prod

EXPOSE 9091

CMD ["/app/auth-service"]
