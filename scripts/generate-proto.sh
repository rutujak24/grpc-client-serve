#!/bin/bash

# Generate Go code from proto files
# Make sure you have protoc, protoc-gen-go, and protoc-gen-go-grpc installed

set -e

PROTO_DIR="proto"
OUT_DIR="."

echo "Generating Go code from protobuf definitions..."

# Generate Go code for gRPC
protoc \
  --proto_path=${PROTO_DIR} \
  --go_out=${OUT_DIR} \
  --go_opt=paths=source_relative \
  --go-grpc_out=${OUT_DIR} \
  --go-grpc_opt=paths=source_relative \
  ${PROTO_DIR}/*.proto

echo "Proto generation completed successfully!"