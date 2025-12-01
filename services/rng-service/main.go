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

	// Import the local protobuf package
	pb "github.com/ShadyDevelopment/ECHOBETZ/services/rng-service/proto"
)

type rngServer struct {
	pb.UnimplementedRNGServiceServer
}

func (s *rngServer) GetNumbers(ctx context.Context, req *pb.RNGRequest) (*pb.RNGResponse, error) {
	// Use math/rand (v1) for simplicity and compatibility
	rand.Seed(time.Now().UnixNano())

	count := int(req.GetCount())
	if count <= 0 { count = 1 }

	numbers := make([]int32, count)
	for i := 0; i < count; i++ {
		// Generate random number (0 to 100 for simplicity)
		numbers[i] = int32(rand.Intn(100))
	}

	// FIX: Use 'Numbers' (Capitalized) matching the generated code
	return &pb.RNGResponse{Numbers: numbers, Seed: time.Now().UnixNano()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRNGServiceServer(s, &rngServer{})
	reflection.Register(s)

	log.Printf("RNG Service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}