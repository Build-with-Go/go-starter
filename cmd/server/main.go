package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Build-with-Go/go-starter/internal/config"
	"github.com/Build-with-Go/go-starter/internal/logger"
	"github.com/Build-with-Go/go-starter/internal/server"
)

func main() {
	// Parse command line flags
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	loggerConfig := &logger.LoggerConfig{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}

	appLogger, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create server
	srv := server.New(cfg, appLogger)

	// Setup graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		appLogger.Logger.Info().Msg("Starting HTTP server")

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			appLogger.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	<-quit
	appLogger.Logger.Info().Msg("Shutting down server...")

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	appLogger.Shutdown("Server stopped")
}
