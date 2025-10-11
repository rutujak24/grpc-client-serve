# Generate Go code from proto files for Windows
# Make sure you have protoc, protoc-gen-go, and protoc-gen-go-grpc installed

$PROTO_DIR = "proto"
$OUT_DIR = "."

Write-Host "Generating Go code from protobuf definitions..." -ForegroundColor Green

# Generate Go code for gRPC
protoc `
  --proto_path=$PROTO_DIR `
  --go_out=$OUT_DIR `
  --go_opt=paths=source_relative `
  --go-grpc_out=$OUT_DIR `
  --go-grpc_opt=paths=source_relative `
  $PROTO_DIR/*.proto

Write-Host "Proto generation completed successfully!" -ForegroundColor Green