# Decomposed Key-Value Store

A simple decomposed key-value store implementation with two services communicating over gRPC:

1. **KV Service**: A gRPC service that implements the key-value storage backend
2. **API Service**: A REST API service that provides JSON HTTP endpoints and communicates with the KV service

## Architecture

```
┌─────────────────┐    HTTP/JSON    ┌─────────────────┐    gRPC    ┌─────────────────┐
│                 │ ◄─────────────  │                 │ ◄───────── │                 │
│   HTTP Client   │                 │   API Service   │            │   KV Service    │
│                 │ ───────────────►│                 │ ──────────►│                 │
└─────────────────┘                 └─────────────────┘            └─────────────────┘
                                           :8080                         :50051
```

## Features

- **Three Core Operations**: Store, Retrieve, and Delete key-value pairs
- **Service Decomposition**: Separate services for API layer and storage layer
- **gRPC Communication**: High-performance inter-service communication
- **REST API Interface**: JSON-based HTTP API for easy integration
- **Docker Support**: Containerized services with docker-compose orchestration
- **Health Checks**: Built-in health monitoring for both services
- **Graceful Shutdown**: Proper signal handling and cleanup

## Quick Start

### Prerequisites

- Go 1.21+ (for local development)
- Docker and Docker Compose
- Protocol Buffers compiler (protoc) - if regenerating proto files

### Running with Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/rutujak24/grpc-client-serve.git
   cd grpc-client-serve
   ```

2. Start both services:
   ```bash
   docker-compose up -d
   ```

3. Verify services are running:
   ```bash
   docker-compose ps
   ```

4. Test the API:
   ```bash
   # Health check
   curl http://localhost:8080/api/v1/health
   
   # Store a key-value pair
   curl -X POST http://localhost:8080/api/v1/kv \
     -H "Content-Type: application/json" \
     -d '{"key": "hello", "value": "world"}'
   
   # Retrieve a value
   curl http://localhost:8080/api/v1/kv/hello
   
   # Delete a key
   curl -X DELETE http://localhost:8080/api/v1/kv/hello
   ```

### Running Locally (Development)

1. Generate protobuf code (if needed):
   ```bash
   # Windows
   .\scripts\generate-proto.ps1
   
   # Unix/Linux/Mac
   ./scripts/generate-proto.sh
   ```

2. Start the KV service:
   ```bash
   go run ./cmd/kv-service
   ```

3. Start the API service (in another terminal):
   ```bash
   go run ./cmd/api-service
   ```

## API Endpoints

### Base URL: `http://localhost:8080/api/v1`

### Store Key-Value Pair
- **Endpoint**: `POST /kv`
- **Request Body**:
  ```json
  {
    "key": "your-key",
    "value": "your-value"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "message": "Successfully stored key 'your-key'"
  }
  ```

### Retrieve Value by Key
- **Endpoint**: `GET /kv/{key}`
- **Response** (if found):
  ```json
  {
    "found": true,
    "value": "your-value",
    "message": "Successfully retrieved key 'your-key'"
  }
  ```
- **Response** (if not found):
  ```json
  {
    "found": false,
    "value": "",
    "message": "Key 'your-key' not found"
  }
  ```

### Delete Key
- **Endpoint**: `DELETE /kv/{key}`
- **Response**:
  ```json
  {
    "success": true,
    "message": "Successfully deleted key 'your-key'"
  }
  ```

### Health Check
- **Endpoint**: `GET /health`
- **Response**:
  ```json
  {
    "status": "healthy",
    "timestamp": 1697040000,
    "service": "api-service"
  }
  ```

## Testing

### Manual Testing with curl

1. **Store a value**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/kv \
     -H "Content-Type: application/json" \
     -d '{"key": "name", "value": "John Doe"}'
   ```

2. **Retrieve the value**:
   ```bash
   curl http://localhost:8080/api/v1/kv/name
   ```

3. **Try to retrieve a non-existent key**:
   ```bash
   curl http://localhost:8080/api/v1/kv/nonexistent
   ```

4. **Delete the key**:
   ```bash
   curl -X DELETE http://localhost:8080/api/v1/kv/name
   ```

5. **Verify deletion**:
   ```bash
   curl http://localhost:8080/api/v1/kv/name
   ```

### Test Scenarios

#### Scenario 1: Basic CRUD Operations
```bash
# Store
curl -X POST http://localhost:8080/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "value": "Alice"}'

# Read
curl http://localhost:8080/api/v1/kv/user:123

# Update (same as store)
curl -X POST http://localhost:8080/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "value": "Alice Smith"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/kv/user:123
```

#### Scenario 2: Error Handling
```bash
# Empty key
curl -X POST http://localhost:8080/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "", "value": "test"}'

# Invalid JSON
curl -X POST http://localhost:8080/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "test"'

# Delete non-existent key
curl -X DELETE http://localhost:8080/api/v1/kv/does-not-exist
```

#### Scenario 3: Load Testing
```bash
# Store multiple values
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/kv \
    -H "Content-Type: application/json" \
    -d "{\"key\": \"test$i\", \"value\": \"value$i\"}" &
done
wait

# Retrieve all values
for i in {1..10}; do
  curl http://localhost:8080/api/v1/kv/test$i &
done
wait
```

### Automated Testing (Using jq for JSON parsing)

If you have `jq` installed, you can create more sophisticated tests:

```bash
#!/bin/bash

# Store a value and check response
response=$(curl -s -X POST http://localhost:8080/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "test", "value": "automated"}')

success=$(echo $response | jq -r '.success')
if [ "$success" = "true" ]; then
  echo "✓ Store operation successful"
else
  echo "✗ Store operation failed"
fi

# Retrieve and verify
response=$(curl -s http://localhost:8080/api/v1/kv/test)
found=$(echo $response | jq -r '.found')
value=$(echo $response | jq -r '.value')

if [ "$found" = "true" ] && [ "$value" = "automated" ]; then
  echo "✓ Retrieve operation successful"
else
  echo "✗ Retrieve operation failed"
fi
```

## Development

### Project Structure
```
├── cmd/
│   ├── api-service/        # REST API service main
│   └── kv-service/         # gRPC KV service main
├── internal/
│   ├── api/                # REST API implementation
│   └── kvstore/            # KV storage implementation
├── proto/                  # Protocol buffer definitions
├── docker/                 # Docker configurations
├── scripts/               # Build and utility scripts
├── docker-compose.yml     # Production deployment
└── docker-compose.dev.yml # Development deployment
```

### Building from Source

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Build both services:
   ```bash
   # KV Service
   go build -o bin/kv-service ./cmd/kv-service
   
   # API Service
   go build -o bin/api-service ./cmd/api-service
   ```

3. Run the services:
   ```bash
   # Terminal 1: KV Service
   ./bin/kv-service
   
   # Terminal 2: API Service
   ./bin/api-service
   ```

#### KV Service
- `KV_SERVICE_PORT`: Port to listen on (default: `:50051`)

#### API Service
- `API_SERVICE_PORT`: Port to listen on (default: `:8080`)
- `KV_SERVICE_ADDR`: Address of KV service (default: `localhost:50051`)

### Docker Commands

```bash
# Build individual services
docker build -f docker/kv-service/Dockerfile -t kv-service .
docker build -f docker/api-service/Dockerfile -t api-service .

# Run individual containers
docker run -p 50051:50051 kv-service
docker run -p 8080:8080 -e KV_SERVICE_ADDR=host.docker.internal:50051 api-service

# Production deployment
docker-compose up -d

# Development deployment with hot reload
docker-compose -f docker-compose.dev.yml up

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```



