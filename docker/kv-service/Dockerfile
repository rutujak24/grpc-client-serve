# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the KV service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kv-service ./cmd/kv-service

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh kvuser

WORKDIR /home/kvuser

# Copy binary from builder stage
COPY --from=builder /app/kv-service .

# Change ownership
RUN chown kvuser:kvuser kv-service

# Switch to non-root user
USER kvuser

# Expose port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD nc -z localhost 50051 || exit 1

# Run the service
CMD ["./kv-service"]