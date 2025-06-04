package main

import (
	"fmt"
	h3_ptf "tessera/src"

	"github.com/uber/h3-go/v4"
)

func main() {
	fmt.Println("🗺️  Tessera H3 Explorer - Quick Start Demo")
	fmt.Println("==========================================")

	// Basic path completion example
	cell := h3.Cell(0x8944d072827ffff)
	fmt.Printf("Starting from H3 cell: %s\n\n", cell.String())

	// Create an explorer that stops when returning to start
	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(10).
		WithMinSimilarity(0).
		WithStopOnReturn(true)

	// Explore and possibly complete the path
	path, err := explorer.Explore(cell)
	if err != nil {
		panic(err)
	}

	// Print the path nicely
	path.Print()

	// Show statistics
	fmt.Printf("\n📊 Path Statistics:\n")
	stats := path.CompletionStats()
	fmt.Printf("   Average similarity: %.2f\n", stats["average_similarity"])
	fmt.Printf("   Min similarity: %.2f\n", stats["min_similarity"])
	fmt.Printf("   Max similarity: %.2f\n", stats["max_similarity"])
	if path.IsComplete() {
		fmt.Printf("   Return similarity: %.2f\n", stats["completion_similarity"])
	}

	fmt.Println("\n🔍 Want to see more exploration scenarios?")
	fmt.Println("   cd examples && go run .")
	fmt.Println("   cd examples && go run . 2  # Stop on return demo")
	fmt.Println("   cd examples && go run . 5  # Detailed analysis")
	fmt.Println("   cd examples && go run . all # All scenarios")
}
