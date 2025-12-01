FROM golang:alpine AS builder
WORKDIR /app

# 1. Install git
RUN apk update && apk add --no-cache git

# 2. Copy root module definition
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy service structure
COPY services/ ./services/
COPY config/ ./config/

# 4. Add explicit local replace directives for internal packages
RUN go mod edit -replace=github.com/ShadyDevelopment/ECHOBETZ/services/rng-service/proto=./services/rng-service/proto
RUN go mod edit -replace=github.com/ShadyDevelopment/ECHOBETZ/services/game-engine-service/proto=./services/game-engine-service/proto

# 5. Run tidy to process replacements
RUN go mod tidy

# 6. Build the Game Engine service
RUN CGO_ENABLED=0 go build -o /usr/local/bin/game-engine-service ./services/game-engine-service/

FROM alpine:latest
# ... rest of DockerfileCOPY --from=builder /usr/local/bin/game-engine-service /usr/local/bin/
EXPOSE 50052
ENTRYPOINT ["/usr/local/bin/game-engine-service"]