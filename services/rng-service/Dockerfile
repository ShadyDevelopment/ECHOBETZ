FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire service directory, including the 'proto' folder
COPY services/rng-service/ ./services/rng-service/

# Build from the directory containing main.go
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./services/rng-service/

FROM alpine:latest
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/rng-service"]