package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// The gRPC client connection to the RNG Service
var rngServiceClient RNGServiceClient

// engineServer implements the GameEngineServer interface defined in engine_grpc.pb.go.
type engineServer struct {
	// Must be embedded to satisfy the interface.
	UnimplementedGameEngineServer
}

// Spin handles a request to spin a reel and returns the result using the RNG service.
func (s *engineServer) Spin(ctx context.Context, req *SpinRequest) (*SpinResponse, error) {
	log.Printf("Received Spin request from User: %s for Game: %s", req.GetUserId(), req.GetGameId())

	// 1. Request random numbers from the RNG Service
	rngRequest := &RNGRequest{Count: 3} // Request 3 random numbers for a simple 3-reel slot
	
	// Note: We use the client types directly, without the module prefix.
	rngResponse, err := rngServiceClient.GetRandomInts(ctx, rngRequest)
	if err != nil {
		log.Printf("Could not connect to RNG service: %v", err)
		return nil, grpc.Errorf(grpc.CodeInternal, "failed to get RNG: %v", err)
	}

	// 2. Simple logic to map random numbers to reel positions (0-9)
	results := make([]int32, len(rngResponse.Ints))
	for i, r := range rngResponse.Ints {
		// Map the large random number to a simple reel position (0-9)
		results[i] = int32(r % 10)
	}
	
	// 3. Determine if it was a win (simplistic: all numbers match)
	win := false
	if len(results) >= 3 && results[0] == results[1] && results[1] == results[2] {
		win = true
	}

	return &SpinResponse{
		ReelResults: results,
		Win:         win,
		Payout:      100, // Example payout
	}, nil
}

func main() {
	// --- 1. Set up RNG Service Client (Connect to Dependency) ---
	rngServiceAddr := os.Getenv("RNG_SERVICE_ADDR")
	if rngServiceAddr == "" {
		// Use the Docker Compose service name and default port
		rngServiceAddr = "rng-service:50051" 
	}

	// Set up connection to RNG Service
	conn, err := grpc.DialContext(context.Background(), rngServiceAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), 
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect to RNG Service at %s: %v", rngServiceAddr, err)
	}
	defer conn.Close()

	// Initialize the client using the now-unprefixed client constructor
	rngServiceClient = NewRNGServiceClient(conn)
	log.Printf("Successfully connected to RNG Service at %s", rngServiceAddr)

	// --- 2. Start Game Engine Server ---
	port := os.Getenv("GAME_ENGINE_PORT")
	if port == "" {
		port = "50052" // Default port for Game Engine service
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// Register the service using the now-unprefixed Register function
	RegisterGameEngineServer(s, &engineServer{})

	// Register reflection service
	reflection.Register(s)

	log.Printf("Game Engine Service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}