package main

import (
	"context"
	"github.com/zahidhasanpapon/iam-bridge/internal/di"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize application using wire
	app, err := di.WireInitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	//// Initialize configuration
	//cfg, err := config.Load()
	//if err != nil {
	//	log.Fatalf("failed to load configuration: %v", err)
	//}
	//
	//// Initialize logger
	//l := logger.NewLogger(cfg.LogLevel)
	//
	//// Create new server instance
	//srv := server.NewServer(cfg, &l)

	// Start server in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

}
