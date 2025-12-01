package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	// You must have the RNG Protobuf files copied into this directory
	// (or imported with a path) for the RNG types (RNGRequest, RNGResponse, RNGServiceClient)
	// to be available in this 'main' package scope.

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// The gRPC client connection to the RNG Service (unprefixed)
var rngServiceClient RngServiceClient

// engineServer implements the GameEngineServer interface defined in engine_grpc.pb.go.
type engineServer struct {
	// FIX 1: Corrected casing of the embedded struct to match engine_grpc.pb.go
	UnimplementedGameEngineServer 
}

// Spin handles a request to spin a reel and returns the result using the RNG service.
func (s *engineServer) Spin(ctx context.Context, req *SpinRequest) (*SpinResponse, error) {
	log.Printf("Received Spin request from User: %s for Game: %s", req.GetUserId(), req.GetGameId())

	// 1. Request random numbers from the RNG Service
	rngRequest := &RNGRequest{Count: 3} 
	
	rngResponse, err := rngServiceClient.GetRandomInts(ctx, rngRequest)
	if err != nil {
		log.Printf("Could not connect to RNG service: %v", err)
		// Use a temporary fix for error formatting, as grpc.Errorf is deprecated
		return nil, grpc.ErrClientConnClosing
	}

	// 2. Simple logic to map random numbers to reel positions (0-9)
	// FIX 2: The field name in the RNGResponse struct is capitalized 'Ints'.
	results := make([]int32, len(rngResponse.Ints))
	for i, r := range rngResponse.Ints {
		results[i] = int32(r % 10)
	}
	
	// 3. Determine win
	win := false
	if len(results) >= 3 && results[0] == results[1] && results[1] == results[2] {
		win = true
	}

	return &SpinResponse{
		ReelResults: results,
		Win:         win,
		Payout:      100,
	}, nil
}

func main() {
	// --- 1. Set up RNG Service Client (Connect to Dependency) ---
	rngServiceAddr := os.Getenv("RNG_SERVICE_ADDR")
	if rngServiceAddr == "" {
		rngServiceAddr = "rng-service:50051" 
	}

	conn, err := grpc.DialContext(context.Background(), rngServiceAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), 
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect to RNG Service at %s: %v", rngServiceAddr, err)
	}
	defer conn.Close()

	// FIX 3: Corrected casing of the client constructor
	rngServiceClient = NewRngServiceClient(conn) 
	log.Printf("Successfully connected to RNG Service at %s", rngServiceAddr)

	// --- 2. Start Game Engine Server ---
	port := os.Getenv("GAME_ENGINE_PORT")
	if port == "" {
		port = "50052" 
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// FIX 4: Corrected casing of the registration function
	RegisterGameEngineServer(s, &engineServer{})

	reflection.Register(s)

	log.Printf("Game Engine Service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}