package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rutujak24/grpc-client-serve/proto"
)

// Server represents the REST API server
type Server struct {
	kvClient pb.KeyValueStoreClient
	router   *mux.Router
}

// Request/Response types for JSON API
type StoreRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StoreResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RetrieveResponse struct {
	Found   bool   `json:"found"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NewServer creates a new REST API server
func NewServer(kvServiceAddr string) (*Server, error) {
	// Connect to KV service
	conn, err := grpc.Dial(kvServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewKeyValueStoreClient(conn)

	server := &Server{
		kvClient: client,
		router:   mux.NewRouter(),
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Key-value operations
	api.HandleFunc("/kv", s.handleStore).Methods("POST")
	api.HandleFunc("/kv/{key}", s.handleRetrieve).Methods("GET")
	api.HandleFunc("/kv/{key}", s.handleDelete).Methods("DELETE")

	// Health check
	api.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Add CORS middleware
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.loggingMiddleware)
}

// GetRouter returns the configured router
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// HTTP Handlers

func (s *Server) handleStore(w http.ResponseWriter, r *http.Request) {
	var req StoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Key == "" {
		s.writeError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call KV service
	resp, err := s.kvClient.Store(ctx, &pb.StoreRequest{
		Key:   req.Key,
		Value: req.Value,
	})
	if err != nil {
		log.Printf("Error calling KV service: %v", err)
		s.writeError(w, http.StatusInternalServerError, "Failed to communicate with KV service")
		return
	}

	response := StoreResponse{
		Success: resp.Success,
		Message: resp.Message,
	}

	s.writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		s.writeError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call KV service
	resp, err := s.kvClient.Retrieve(ctx, &pb.RetrieveRequest{
		Key: key,
	})
	if err != nil {
		log.Printf("Error calling KV service: %v", err)
		s.writeError(w, http.StatusInternalServerError, "Failed to communicate with KV service")
		return
	}

	response := RetrieveResponse{
		Found:   resp.Found,
		Value:   resp.Value,
		Message: resp.Message,
	}

	statusCode := http.StatusOK
	if !resp.Found {
		statusCode = http.StatusNotFound
	}

	s.writeJSON(w, statusCode, response)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		s.writeError(w, http.StatusBadRequest, "Key cannot be empty")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call KV service
	resp, err := s.kvClient.Delete(ctx, &pb.DeleteRequest{
		Key: key,
	})
	if err != nil {
		log.Printf("Error calling KV service: %v", err)
		s.writeError(w, http.StatusInternalServerError, "Failed to communicate with KV service")
		return
	}

	response := DeleteResponse{
		Success: resp.Success,
		Message: resp.Message,
	}

	statusCode := http.StatusOK
	if !resp.Success {
		statusCode = http.StatusNotFound
	}

	s.writeJSON(w, statusCode, response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "api-service",
	})
}

// Helper methods

func (s *Server) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) writeError(w http.ResponseWriter, statusCode int, message string) {
	s.writeJSON(w, statusCode, ErrorResponse{
		Error: message,
	})
}

// Middleware

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}