package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rutujak24/grpc-client-serve/internal/kvstore"
	pb "github.com/rutujak24/grpc-client-serve/proto"
)

const (
	defaultPort = ":50051"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("KV_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	// Create a listener on TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	server := grpc.NewServer()

	// Register our service
	kvServer := kvstore.NewServer()
	pb.RegisterKeyValueStoreServer(server, kvServer)

	// Enable reflection for easier testing with tools like grpcurl
	reflection.Register(server)

	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gRPC server...")
		server.GracefulStop()
	}()

	log.Printf("Key-Value gRPC server starting on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}