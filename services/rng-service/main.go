package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// rngServer implements the RNGServiceServer interface defined in rng_grpc.pb.go.
type rngServer struct {
	// Must be embedded to satisfy the interface.
	UnimplementedRNGServiceServer
}

// GetRandomInts generates a list of random 64-bit integers.
func (s *rngServer) GetRandomInts(ctx context.Context, req *RNGRequest) (*RNGResponse, error) {
	// Initialize a new pseudo-random source based on the current time
	// Note: We are using the newer math/rand/v2 which is imported implicitly 
	// or the older rand.NewSource if running Go < 1.20
	source := rand.New(rand.NewSource(time.Now().UnixNano()))

	count := int(req.GetCount())
	if count <= 0 || count > 1000 {
		count = 1 // Default to 1 if count is invalid
	}

	ints := make([]int64, count)
	for i := 0; i < count; i++ {
		// Generate a random int64
		ints[i] = source.Int63()
	}

	return &RNGResponse{Ints: ints}, nil
}

func main() {
	port := os.Getenv("RNG_PORT")
	if port == "" {
		port = "50051" // Default port for RNG service
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// Register the service using the now-unprefixed Register function
	RegisterRNGServiceServer(s, &rngServer{})

	// Register reflection service on gRPC server for testing tools
	reflection.Register(s)

	log.Printf("RNG Service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}