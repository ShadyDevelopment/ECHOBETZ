package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import the Protobuf packages using your defined module path
	// The path for the Game Engine's server implementation
	pb_engine "github.com/ShadyDevelopment/ECHOBETZ/services/game-engine-service"
	// The path for the RNG Client connection
	pb_rng "github.com/ShadyDevelopment/ECHOBETZ/services/rng-service" 
)

const (
	// The port the Game Engine will listen on
	port = ":50052"
	// The address for the RNG Service (using Docker Compose service name)
	rngServiceAddr = "rng-service:50051"
	// Path to the game configuration file (relative to the container's WORKDIR)
	configPath = "/app/config/aurora_star.json" 
)

// Define the server struct that will implement the gRPC methods
type server struct {
	pb_engine.UnimplementedGameEngineServiceServer
	rngClient pb_rng.RNGServiceClient
	// Store the game configuration (e.g., loaded from aurora_star.json)
	gameConfig GameConfiguration 
}

// --- Configuration Structs (Minimal Representation) ---

// GameConfiguration holds the reel strips and paytable (simplified)
type GameConfiguration struct {
	ReelStrips [][]string // 5 reels, each with a long list of symbols
	Paytable   map[string]float64
}

// loadConfig simulates loading configuration data
func loadConfig(path string) GameConfiguration {
	// In a real system, you would read the JSON file here.
	// For this minimal demo, we hardcode simple, short reels.
	log.Printf("Loading configuration from %s...", path)
	
	// Define short, dummy reel strips for quick testing
	reel1 := []string{"S_WILD", "S_HIGH_A", "S_LOW_E", "S_HIGH_A", "S_LOW_D"}
	reel2 := []string{"S_HIGH_A", "S_SCATTER", "S_LOW_D", "S_HIGH_A", "S_WILD"}
	reel3 := []string{"S_LOW_E", "S_WILD", "S_SCATTER", "S_LOW_D", "S_LOW_E"}
	reel4 := []string{"S_HIGH_A", "S_LOW_D", "S_WILD", "S_HIGH_A", "S_SCATTER"}
	reel5 := []string{"S_SCATTER", "S_LOW_E", "S_HIGH_A", "S_WILD", "S_LOW_D"}

	return GameConfiguration{
		ReelStrips: [][]string{reel1, reel2, reel3, reel4, reel5},
		Paytable: map[string]float64{
			"S_WILD":    100.0,
			"S_SCATTER": 50.0,
			"S_HIGH_A":  20.0,
			"S_LOW_D":   10.0,
			"S_LOW_E":   5.0,
		},
	}
}

// --- gRPC Implementation ---

// PerformSpin implements the GameEngineService's PerformSpin method.
func (s *server) PerformSpin(ctx context.Context, req *pb_engine.SpinRequest) (*pb_engine.SpinResponse, error) {
	log.Printf("Received spin request for player: %s, bet: %.2f", req.PlayerId, req.BetAmount)

	// 1. Call RNG Service to get random stop positions (Numbers)
	rngReq := &pb_rng.RNGRequest{
		Count: 5, // We need 5 numbers, one for each reel
		Min:   0,
		// Max is based on the longest reel strip length to ensure a valid index
		Max:   int32(len(s.gameConfig.ReelStrips[0])), 
	}

	rngRes, err := s.rngClient.GetNumbers(ctx, rngReq)
	if err != nil {
		log.Printf("ERROR: Failed to call RNG Service: %v", err)
		return nil, fmt.Errorf("rng service failure: %w", err)
	}

	log.Printf("Received RNG results: %v", rngRes.Numbers)
	
	// 2. Resolve Spin and Calculate Win
	
	// 2.1. Determine the final 3x5 symbol matrix
	matrix, totalWin := s.resolveSpin(rngRes.Numbers)

	// 3. Prepare Response
	
	// Convert the 3x5 matrix (currently a slice of slices of strings) into the Protobuf format
	protoMatrix := make([]*pb_engine.SpinResponse_Row, len(matrix))
	for i, row := range matrix {
		protoMatrix[i] = &pb_engine.SpinResponse_Row{
			Symbols: row,
		}
	}

	response := &pb_engine.SpinResponse{
		Matrix:    protoMatrix,
		TotalWin:  float32(totalWin),
		// In a real system, WinLines, FeatureTriggers, and AuditData would be populated here.
		AuditData: "RNG Seed Used: " + fmt.Sprint(rngRes.Seed),
	}

	log.Printf("Spin resolved. Win: %.2f", totalWin)
	return response, nil
}

// resolveSpin calculates the final matrix and the resulting win amount
func (s *server) resolveSpin(stopPositions []int32) ([][]string, float64) {
	// The visible 3x5 matrix (3 rows x 5 reels)
	matrix := make([][]string, 3) 
	for i := range matrix {
		matrix[i] = make([]string, 5)
	}
	
	totalWin := 0.0

	// 1. Map stop positions to symbols
	for r := 0; r < 5; r++ { // Iterate through 5 reels
		reelLen := len(s.gameConfig.ReelStrips[r])
		
		// The RNG stop position defines the index of the TOP symbol (Row 0)
		startPos := int(stopPositions[r]) 

		for i := 0; i < 3; i++ { // Iterate through 3 visible rows
			// Use modulo arithmetic to wrap the reel strip
			symbolIndex := (startPos + i) % reelLen
			symbol := s.gameConfig.ReelStrips[r][symbolIndex]
			matrix[i][r] = symbol
		}
	}
	
	// 2. Simplified Payline Check (only checks for 5-of-a-kind on the center row)
	centerRow := matrix[1]
	
	// Check if all symbols on the center row are the same (excluding WILD for simplicity)
	if centerRow[0] != "" && centerRow[0] == centerRow[1] && centerRow[0] == centerRow[2] && centerRow[0] == centerRow[3] && centerRow[0] == centerRow[4] {
		winAmount, exists := s.gameConfig.Paytable[centerRow[0]]
		if exists {
			// Total win is win amount * bet multiplier (we assume 1x for simplicity)
			totalWin = winAmount 
		}
	}

	return matrix, totalWin
}

// --- Main Setup ---

func main() {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
	// 1. Establish gRPC connection to RNG Service
	log.Printf("Attempting to connect to RNG Service at %s...", rngServiceAddr)

	// Use insecure connection for internal Docker network communication
	conn, err := grpc.Dial(rngServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: Could not connect to RNG Service: %v", err)
	}
	defer conn.Close()

	rngClient := pb_rng.NewRNGServiceClient(conn)
	log.Println("Successfully connected to RNG Service.")

	// 2. Load Game Configuration
	config := loadConfig(configPath)
	log.Printf("Game configuration loaded. Total Reels: %d", len(config.ReelStrips))

	// 3. Start gRPC Server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("FATAL: Failed to listen: %v", err)
	}

	// Create new gRPC server instance
	s := grpc.NewServer()
	
	// Register the Game Engine server implementation
	pb_engine.RegisterGameEngineServiceServer(s, &server{
		rngClient: rngClient,
		gameConfig: config,
	})

	log.Printf("Game Engine Service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("FATAL: Failed to serve: %v", err)
	}
}