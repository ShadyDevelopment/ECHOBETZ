# Stage 1: Builder for Go Executable
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# 1. Copy go.mod and go.sum from the project root (context)
COPY go.mod go.sum ./

# 2. Download dependencies (to leverage Docker cache)
# This uses the files copied to /app/go.mod
RUN go mod download

# 3. Copy Source Code Structure: 
# We must copy the entire directory structure to satisfy the 'replace' directives in go.mod.

# Copy the Game Engine source code
COPY services/game-engine-service/ ./services/game-engine-service/

# Copy the RNG service source code (needed for the Game Engine's gRPC client)
COPY services/rng-service/ ./services/rng-service/

# 4. Copy the game configuration file
COPY config/aurora_star.json /app/config/aurora_star.json

# 5. Build the application: 
# Target the main package inside the replicated directory structure.
RUN CGO_ENABLED=0 go build -o /usr/local/bin/game-engine-service ./services/game-engine-service/

# Stage 2: Minimal Production Image
FROM alpine:latest
# Install CA certificates for secure connections (best practice)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary
COPY --from=builder /usr/local/bin/game-engine-service /usr/local/bin/

# Expose the gRPC port
EXPOSE 50052

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/game-engine-service"]