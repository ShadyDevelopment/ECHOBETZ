FROM golang:alpine AS builder
WORKDIR /app

# 1. Copy root module definition (go.mod, go.sum)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy the ENTIRE 'services' directory structure.
# This ensures that relative packages (like /services/rng-service/proto) 
# are available locally within the module.
COPY services/ ./services/

# 3. Build the RNG service
# The 'go build' command should now correctly find the local proto package.
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./services/rng-service/

FROM alpine:latest
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/rng-service"]