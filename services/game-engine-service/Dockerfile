# Stage 1: Builder for Go Executable
FROM golang:alpine AS builder

WORKDIR /app

# 1. Copy go.mod and go.sum from the project root (context)
COPY go.mod go.sum ./

# 2. Download dependencies
RUN go mod download

# 3. FIX: Copy Source Code Structure
# Replicate the structure for all services imported by the Game Engine.
COPY services/game-engine-service/ ./services/game-engine-service/
COPY services/rng-service/ ./services/rng-service/

# 4. Copy the game configuration file (assuming it's in a top-level config directory)
COPY config/aurora_star.json /app/config/aurora_star.json

# 5. FIX: Build the application from the correct path
RUN CGO_ENABLED=0 go build -o /usr/local/bin/game-engine-service ./services/game-engine-service/

# Stage 2: Minimal Production Image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

# Copy the compiled binary
COPY --from=builder /usr/local/bin/game-engine-service /usr/local/bin/

EXPOSE 50052

ENTRYPOINT ["/usr/local/bin/game-engine-service"]