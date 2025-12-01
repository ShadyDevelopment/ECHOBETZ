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

// rngServer implements the RngServiceServer interface.
type rngServer struct {
	// FIX 1: Corrected casing of the embedded struct to match rng_grpc.pb.go
	UnimplementedRngServiceServer 
}

// GetRandomInts generates a list of random 64-bit integers.
func (s *rngServer) GetRandomInts(ctx context.Context, req *RNGRequest) (*RNGResponse, error) {
	// Note: rand.NewSource is deprecated in newer Go versions, but this keeps compatibility.
	source := rand.New(rand.NewSource(time.Now().UnixNano()))

	count := int(req.GetCount())
	if count <= 0 || count > 1000 {
		count = 1 
	}

	ints := make([]int64, count)
	for i := 0; i < count; i++ {
		ints[i] = source.Int63()
	}

	// FIX 2: The field name in the RNGResponse struct is capitalized 'Ints'.
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
	// FIX 3: Corrected casing of the registration function to match rng_grpc.pb.go
	RegisterRngServiceServer(s, &rngServer{})

	reflection.Register(s)

	log.Printf("RNG Service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}