package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rutujak24/grpc-client-serve/internal/api"
)

const (
	defaultAPIPort = ":8080"
	defaultKVAddr  = "localhost:50051"
)

func main() {
	// Get configuration from environment variables
	apiPort := os.Getenv("API_SERVICE_PORT")
	if apiPort == "" {
		apiPort = defaultAPIPort
	}

	kvServiceAddr := os.Getenv("KV_SERVICE_ADDR")
	if kvServiceAddr == "" {
		kvServiceAddr = defaultKVAddr
	}

	// Create API server
	apiServer, err := api.NewServer(kvServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create API server: %v", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         apiPort,
		Handler:      apiServer.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		log.Println("Shutting down API server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()

	log.Printf("API server starting on port %s", apiPort)
	log.Printf("Connecting to KV service at %s", kvServiceAddr)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	log.Println("API server stopped")
}