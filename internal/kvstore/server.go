package kvstore

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/rutujak24/grpc-client-serve/proto"
)

// Server implements the KeyValueStore gRPC service
type Server struct {
	pb.UnimplementedKeyValueStoreServer
	store map[string]string
	mutex sync.RWMutex
}

// NewServer creates a new KeyValueStore server instance
func NewServer() *Server {
	return &Server{
		store: make(map[string]string),
	}
}

// Store implements the Store RPC method
func (s *Server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	if req.Key == "" {
		return &pb.StoreResponse{
			Success: false,
			Message: "Key cannot be empty",
		}, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[req.Key] = req.Value

	return &pb.StoreResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully stored key '%s'", req.Key),
	}, nil
}

// Retrieve implements the Retrieve RPC method
func (s *Server) Retrieve(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	if req.Key == "" {
		return &pb.RetrieveResponse{
			Found:   false,
			Value:   "",
			Message: "Key cannot be empty",
		}, nil
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.store[req.Key]
	if !exists {
		return &pb.RetrieveResponse{
			Found:   false,
			Value:   "",
			Message: fmt.Sprintf("Key '%s' not found", req.Key),
		}, nil
	}

	return &pb.RetrieveResponse{
		Found:   true,
		Value:   value,
		Message: fmt.Sprintf("Successfully retrieved key '%s'", req.Key),
	}, nil
}

// Delete implements the Delete RPC method
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.Key == "" {
		return &pb.DeleteResponse{
			Success: false,
			Message: "Key cannot be empty",
		}, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.store[req.Key]
	if !exists {
		return &pb.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("Key '%s' not found", req.Key),
		}, nil
	}

	delete(s.store, req.Key)

	return &pb.DeleteResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully deleted key '%s'", req.Key),
	}, nil
}