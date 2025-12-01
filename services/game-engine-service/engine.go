package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Simplified structures matching config/aurora_star.json
type GameConfig struct {
	GameCode string `json:"game_code"`
	Grid     struct {
		Rows  int `json:"rows"`
		Reels int `json:"reels"`
	} `json:"grid"`
	Paylines [][]int `json:"paylines"`
	// ReelStrips and other fields are loaded here
	ReelStrips [][]string `json:"reel_strips"`
}

var loadedConfig GameConfig // Global or cached config

func LoadGameConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &loadedConfig)
}

// SpinResult holds the outcome of a game round
type SpinResult struct {
	Matrix [][]string `json:"matrix"`
	TotalWin int `json:"total_win"`
	WinLines []string `json:"win_lines"` // Simplified win detail
}

// PerformSpin simulates the spin and win evaluation.
// rngOutputs are the stop indices received from the RNG service.
func PerformSpin(rngOutputs []int64, betAmount int) SpinResult {
	if len(rngOutputs) != loadedConfig.Grid.Reels {
		log.Fatal("RNG output count mismatch with reel count")
	}

	resultMatrix := make([][]string, loadedConfig.Grid.Reels)
	
	// 1. Determine Stop Positions and Reel Matrix
	for i := 0; i < loadedConfig.Grid.Reels; i++ {
		strip := loadedConfig.ReelStrips[i]
		stopIndex := int(rngOutputs[i]) % len(strip)
		
		// Extract the visible window (3 symbols)
		resultMatrix[i] = make([]string, loadedConfig.Grid.Rows)
		for j := 0; j < loadedConfig.Grid.Rows; j++ {
			// Calculate index with wrap-around logic
			symbolIndex := (stopIndex + j) % len(strip)
			resultMatrix[i][j] = strip[symbolIndex]
		}
	}
	
	// Transpose the matrix for easier evaluation (Reels x Rows -> Rows x Reels)
	finalMatrix := make([][]string, loadedConfig.Grid.Rows)
	for r := 0; r < loadedConfig.Grid.Rows; r++ {
		finalMatrix[r] = make([]string, loadedConfig.Grid.Reels)
		for c := 0; c < loadedConfig.Grid.Reels; c++ {
			finalMatrix[r][c] = resultMatrix[c][r]
		}
	}


	// 2. Win Evaluation (Highly simplified, full logic is complex)
	totalWin := 0
	winLines := []string{}
	// For Aurora Star (20 Paylines), iterate through paylines to check matches
	// ... actual win calculation based on paytable and line matches ...
	
	// Placeholder: Award a simple win if the middle symbol on reel 3 is the top symbol
	if finalMatrix[1][2] == "S_HIGH_A" {
	    totalWin = betAmount * 5
	    winLines = append(winLines, "CENTER_MATCH")
	}

	log.Printf("Spin resolved. Matrix: %v, Win: %d", finalMatrix, totalWin)
	
	return SpinResult{
		Matrix:   finalMatrix,
		TotalWin: totalWin,
		WinLines: winLines,
	}
}