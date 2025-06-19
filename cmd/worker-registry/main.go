package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ncolesummers/mindgateway/internal/registry"
	"github.com/ncolesummers/mindgateway/internal/shared/config"
	"github.com/ncolesummers/mindgateway/internal/shared/logging"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger(cfg.LogLevel)

	// Create registry service
	svc, err := registry.NewService(
		registry.WithConfig(cfg),
		registry.WithLogger(logger),
	)
	if err != nil {
		logger.Fatalf("Failed to create registry service: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	svc.RegisterWithServer(grpcServer)

	// Start server
	go func() {
		if err := svc.Start(grpcServer); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down registry service...")
	
	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	
	if err := svc.Shutdown(ctx); err != nil {
		logger.Errorf("Service forced to shutdown: %v", err)
	}
	
	logger.Info("Registry service exited")
}