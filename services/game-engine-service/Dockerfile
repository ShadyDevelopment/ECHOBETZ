FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# FIX: Copy ALL services to ensure cross-service imports (like RNG proto) work.
COPY services/ ./services/
COPY config/ ./config/

RUN CGO_ENABLED=0 go build -o /usr/local/bin/game-engine-service ./services/game-engine-service/

FROM alpine:latest
COPY --from=builder /usr/local/bin/game-engine-service /usr/local/bin/
EXPOSE 50052
ENTRYPOINT ["/usr/local/bin/game-engine-service"]