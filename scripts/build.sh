#!/bin/bash

# Build script for the decomposed key-value store services
set -e

echo "Building Decomposed Key-Value Store Services..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Create bin directory
mkdir -p bin

# Build flags
BUILD_FLAGS="-ldflags=-w -s"
if [ "$1" = "debug" ]; then
    BUILD_FLAGS=""
    print_status "Building in debug mode..."
else
    print_status "Building in release mode..."
fi

# Build KV service
print_status "Building KV service..."
go build $BUILD_FLAGS -o bin/kv-service ./cmd/kv-service
if [ $? -eq 0 ]; then
    print_status "✓ KV service built successfully"
else
    print_error "✗ Failed to build KV service"
    exit 1
fi

# Build API service
print_status "Building API service..."
go build $BUILD_FLAGS -o bin/api-service ./cmd/api-service
if [ $? -eq 0 ]; then
    print_status "✓ API service built successfully"
else
    print_error "✗ Failed to build API service"
    exit 1
fi

print_status "All services built successfully!"
print_status "Binaries are available in the ./bin/ directory:"
ls -la bin/

echo ""
print_status "To run the services:"
echo "  KV Service:  ./bin/kv-service"
echo "  API Service: ./bin/api-service"