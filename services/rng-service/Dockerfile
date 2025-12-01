# Stage 1: Builder for Go Executable
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# 1. Copy go.mod and go.sum from the project root (context)
COPY go.mod go.sum ./

# 2. Download dependencies
RUN go mod download

# 3. Copy Source Code Structure: 
# We must copy the entire directory structure to satisfy the 'replace' directive.

# Copy the RNG service source code
COPY services/rng-service/ ./services/rng-service/

# 4. Build the application: 
# Target the main package inside the replicated directory structure.
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./services/rng-service/

# Stage 2: Minimal Production Image
FROM alpine:latest
# Copy the compiled binary
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/

# Expose the gRPC port
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/rng-service"]