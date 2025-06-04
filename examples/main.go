package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	scenarios := map[string]func(){
		"1": BasicExploration,
		"2": PathCompletionStrategies,
		"3": MultipleStartingPoints,
		"4": SimilarityThresholdComparison,
		"5": DetailedPathAnalysis,
		"6": CallbackExploration,
		"7": ExtremeCases,
		"8": EdgeIndexDemo,
	}

	if len(os.Args) > 1 {
		// Run specific scenario if argument provided
		scenario := os.Args[1]
		if fn, exists := scenarios[scenario]; exists {
			fn()
			return
		}

		// Try to run all scenarios matching a pattern
		if scenario == "all" {
			fmt.Println("🗺️  Running All Tessera Exploration Scenarios")
			for i := 1; i <= len(scenarios); i++ {
				scenarios[strconv.Itoa(i)]()
			}
			return
		}

		fmt.Printf("Unknown scenario: %s\n", scenario)
	}

	// Interactive menu
	fmt.Println("🗺️  Tessera H3 Explorer - Scenario Demos")
	fmt.Println("=========================================")
	fmt.Println("1. Basic Exploration")
	fmt.Println("2. Path Completion Strategies")
	fmt.Println("3. Multiple Starting Points")
	fmt.Println("4. Similarity Threshold Comparison")
	fmt.Println("5. Detailed Path Analysis")
	fmt.Println("6. Callback-Controlled Exploration")
	fmt.Println("7. Extreme Cases & Edge Scenarios")
	fmt.Println("8. Edge Index Demo")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run . [scenario_number]  # Run specific scenario")
	fmt.Println("  go run . all               # Run all scenarios")
	fmt.Println("  go run .                   # Show this menu")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run . 2                 # Run path completion demo")
	fmt.Println("  go run . 8                 # Run edge index demo")
}
