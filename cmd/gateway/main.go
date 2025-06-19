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