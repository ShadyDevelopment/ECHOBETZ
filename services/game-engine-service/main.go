package main

import (
	"context"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	pb_engine "path/to/game-engine-service/engine" // Assuming your generated protobuf paths
	pb_rng "path/to/rng/rng" // Assuming your generated protobuf paths
)

const (
	port        = ":50052"
	rngServerAddr = "localhost:50051" // Address of the running RNG service
	reelCount   = 5
)

// Server implements the gRPC interface for the Game Engine
type server struct {
	pb_engine.UnimplementedGameEngineServer
}

// Spin handles the core spin request from the Integration Gateway
func (s *server) Spin(ctx context.Context, in *pb_engine.SpinRequest) (*pb_engine.SpinResponse, error) {
	log.Printf("Received Spin request for Game: %s, Bet: %d", in.GameCode, in.BetAmount)

	// --- 1. Internal Call to RNG Service ---
	rngConn, err := grpc.Dial(rngServerAddr, grpc.WithInsecure()) // Note: Use WithTransportCredentials for production TLS/gRPC
	if err != nil {
		return nil, logError("Failed to connect to RNG Service", err)
	}
	defer rngConn.Close()

	rngClient := pb_rng.NewRNGClient(rngConn)
	rngRequest := &pb_rng.RNGRequest{Count: reelCount}

	rngResponse, err := rngClient.GetRandomNumbers(ctx, rngRequest)
	if err != nil {
		return nil, logError("Failed to get random numbers", err)
	}
	
	// Ensure stop numbers are positive integers for modulo operations
	stopNumbers := make([]int64, reelCount)
    for i, num := range rngResponse.Numbers {
        // Ensure result is positive (gRPC int64)
        stopNumbers[i] = num & 0x7FFFFFFFFFFFFFFF 
    }
	
	// --- 2. Perform Spin and Evaluation ---
	// Assume GameConfig is loaded via LoadGameConfig in main()
	result := PerformSpin(stopNumbers, int(in.BetAmount))
	
	// --- 3. Return outcome to Gateway ---
	
	// Convert result matrix back to a flat array or protobuf-friendly structure for transport
	var flatMatrix []string
	for _, row := range result.Matrix {
		flatMatrix = append(flatMatrix, row...)
	}

	return &pb_engine.SpinResponse{
		Matrix:    flatMatrix,
		TotalWin:  int64(result.TotalWin),
		WinDetails: result.WinLines,
		RngSeed:   rngResponse.Seed,
	}, nil
}

func logError(msg string, err error) error {
	log.Printf("%s: %v", msg, err)
	return err
}

func main() {
	// 1. Load Game Configuration
	// Path needs adjustment based on where you run the service
	if err := LoadGameConfig("../config/aurora_star.json"); err != nil {
		log.Fatalf("Failed to load AURORA_STAR config: %v", err)
	}
	
	// 2. Start gRPC Server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// NOTE: You need to define the pb_engine proto file and generate code first!
	// pb_engine.RegisterGameEngineServer(s, &server{}) 
	log.Printf("Game Engine server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}