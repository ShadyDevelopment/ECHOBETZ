package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// Import local engine proto
	pb_engine "github.com/ShadyDevelopment/ECHOBETZ/services/game-engine-service/proto"
	// Import remote RNG proto (Works now because it's a library package!)
	pb_rng "github.com/ShadyDevelopment/ECHOBETZ/services/rng-service/proto"
)

type engineServer struct {
	pb_engine.UnimplementedGameEngineServer
	rngClient pb_rng.RNGServiceClient
}

func (s *engineServer) Spin(ctx context.Context, req *pb_engine.SpinRequest) (*pb_engine.SpinResponse, error) {
	// Call RNG Service
	rngResp, err := s.rngClient.GetNumbers(ctx, &pb_rng.RNGRequest{Count: 5})
	if err != nil {
		log.Printf("Error calling RNG: %v", err)
		return nil, err
	}

	// Logic: Use the random numbers
	// FIX: Use 'Numbers' field (Capitalized)
	log.Printf("Got RNG numbers: %v", rngResp.Numbers)

	// Dummy response
	return &pb_engine.SpinResponse{
		Matrix: []string{"A", "B", "C"},
		TotalWin: 100,
	}, nil
}

func main() {
	// Connect to RNG Service
	conn, err := grpc.Dial("rng-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to RNG: %v", err)
	}
	defer conn.Close()
	rngClient := pb_rng.NewRNGServiceClient(conn)

	// Start Engine Server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb_engine.RegisterGameEngineServer(s, &engineServer{rngClient: rngClient})
	reflection.Register(s)

	log.Printf("Game Engine listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}