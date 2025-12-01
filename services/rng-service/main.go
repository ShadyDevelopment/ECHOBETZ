package main

import (
    "context"
    "crypto/sha256"
    "fmt"
    "log"
    "math/rand/v2"
    "net"
    "strconv"
    "time"

    "google.golang.org/grpc"
    pb "path/to/rng/rng" // Assume generated Go code path
)

const (
    port = ":50051"
)

type server struct {
    pb.UnimplementedRNGServer
}

// GetRandomNumbers implements rng.RNGServer
func (s *server) GetRandomNumbers(ctx context.Context, in *pb.RNGRequest) (*pb.RNGResponse, error) {
    // 1. Generate a cryptographically secure seed
    seed := time.Now().UnixNano()
    // In a real system, we'd use crypto/rand for better seeding
    source := rand.NewPCG(uint64(seed), rand.NewSource(seed).Split().S)

    randomNumbers := make([]int64, in.Count)
    for i := 0; i < int(in.Count); i++ {
        // Generate a large positive 64-bit random number
        randomNumbers[i] = source.Uint64() 
    }

    log.Printf("RNG Call: Count=%d, Seed=%d, Output (first 3)=%v", in.Count, seed, randomNumbers[:3])

    return &pb.RNGResponse{
        Numbers: randomNumbers,
        Seed:    fmt.Sprintf("%d", seed), // Convert seed to string for transmission
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterRNGServer(s, &server{})
    log.Printf("RNG server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}