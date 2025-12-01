# Stage 1: Build the Go executable
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./main.go

# Stage 2: Create a minimal production image
FROM alpine:latest
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/rng-service"]