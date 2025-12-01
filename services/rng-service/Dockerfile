FROM golang:alpine AS builder
WORKDIR /app

# 1. Install git
RUN apk update && apk add --no-cache git

# 2. Copy root module definition
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy the entire services directory structure.
COPY services/ ./services/

# 4. FIX: Add explicit local replace directives for ALL internal packages the service uses.
# RNG Service needs its own proto and likely the Game Engine's proto as well.
RUN go mod edit -replace=github.com/ShadyDevelopment/ECHOBETZ/services/rng-service/proto=./services/rng-service/proto
RUN go mod edit -replace=github.com/ShadyDevelopment/ECHOBETZ/services/game-engine-service/proto=./services/game-engine-service/proto

# 5. Run tidy to process replacements
RUN go mod tidy

# 6. Build the RNG service
RUN CGO_ENABLED=0 go build -o /usr/local/bin/rng-service ./services/rng-service/

FROM alpine:latest
COPY --from=builder /usr/local/bin/rng-service /usr/local/bin/
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/rng-service"]