# Stage 1: Builder for Go Executable
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# 1. Copy go.mod and go.sum from the project root (context)
COPY go.mod go.sum ./

# 2. Download dependencies
RUN go mod download

# 3. FIX: Copy Source Code Structure
# Replicate the services/rng-service structure required by the 'replace' directive.
COPY services/rng-service/ ./services/rng-service/

# 4. FIX: Build the application from the correct path
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./services/rng-service/

# Stage 2: Minimal Production Image
FROM alpine:latest
# Install CA certificates (good practice)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/

EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/rng-service"]