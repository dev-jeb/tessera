package main

import (
	"fmt"
	h3_ptf "tessera/src"

	"github.com/uber/h3-go/v4"
)

// BasicExploration demonstrates simple pathfinding without completion
func BasicExploration() {
	fmt.Println("=== Basic Exploration Demo ===")
	cell := h3.Cell(0x8844d072a3fffff)

	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(10).
		WithMinSimilarity(0.1)

	path, _ := explorer.Explore(cell)

	fmt.Printf("Starting from: %s\n", cell.String())
	fmt.Printf("Path length: %d steps\n", len(path.Steps))
	fmt.Printf("Is complete: %t\n", path.IsComplete())

	fmt.Println("Path progression:")
	for i, step := range path.Steps {
		fmt.Printf("  Step %d: %.2f -> %s\n", i+1, step.Similarity, step.Index.String())
	}
	fmt.Println()
}

// PathCompletionStrategies demonstrates basic exploration with and without stop on return
func PathCompletionStrategies() {
	fmt.Println("=== Stop On Return Demo ===")
	cell := h3.Cell(0x8844d072a3fffff)

	// Without stop on return
	fmt.Println("--- Without Stop On Return ---")
	for _, maxSteps := range []int{5, 10, 15} {
		explorer := h3_ptf.NewExplorer().
			WithMaxSteps(maxSteps).
			WithMinSimilarity(0.1).
			WithStopOnReturn(false)

		path, _ := explorer.Explore(cell)
		stats := path.CompletionStats()

		fmt.Printf("  Max steps %d: Length=%d, Complete=%t",
			maxSteps, stats["path_length"], stats["is_complete"])

		if len(path.Steps) > 0 {
			fmt.Printf(", AvgSim=%.2f", stats["average_similarity"])
		}
		fmt.Println()
	}

	// With stop on return
	fmt.Println("\n--- With Stop On Return ---")
	for _, maxSteps := range []int{5, 10, 15} {
		explorer := h3_ptf.NewExplorer().
			WithMaxSteps(maxSteps).
			WithMinSimilarity(0.1).
			WithStopOnReturn(true)

		path, _ := explorer.Explore(cell)
		stats := path.CompletionStats()

		fmt.Printf("  Max steps %d: Length=%d, Complete=%t",
			maxSteps, stats["path_length"], stats["is_complete"])

		if stats["is_complete"].(bool) {
			fmt.Printf(", AvgSim=%.2f, ReturnSim=%.2f",
				stats["average_similarity"], stats["completion_similarity"])
		} else if len(path.Steps) > 0 {
			fmt.Printf(", AvgSim=%.2f", stats["average_similarity"])
		}
		fmt.Println()
	}
	fmt.Println()
}

// MultipleStartingPoints compares paths from different origins
func MultipleStartingPoints() {
	fmt.Println("=== Multiple Starting Points Demo ===")

	testCells := []h3.Cell{
		h3.Cell(0x8844d072a3fffff),
		h3.Cell(0x8844d070c9fffff),
		h3.Cell(0x8844d072b5fffff),
	}

	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(12).
		WithMinSimilarity(0.2).
		WithStopOnReturn(true)

	for i, startCell := range testCells {
		path, _ := explorer.Explore(startCell)
		stats := path.CompletionStats()

		fmt.Printf("=== Starting Point %d ===\n", i+1)
		path.Print()
		fmt.Printf("Complete: %t", stats["is_complete"])
		if stats["is_complete"].(bool) {
			fmt.Printf(", avg_sim: %.2f", stats["average_similarity"])
		}
		fmt.Println("\n")
	}
}

// SimilarityThresholdComparison shows how different thresholds affect exploration
func SimilarityThresholdComparison() {
	fmt.Println("=== Similarity Threshold Comparison ===")
	cell := h3.Cell(0x8844d072a3fffff)

	thresholds := []float32{0.0, 0.5, 0.8, 0.9}

	for _, threshold := range thresholds {
		explorer := h3_ptf.NewExplorer().
			WithMaxSteps(15).
			WithMinSimilarity(threshold).
			WithStopOnReturn(true)

		path, _ := explorer.Explore(cell)
		stats := path.CompletionStats()

		fmt.Printf("Threshold %.1f: %d steps, complete=%t",
			threshold, stats["path_length"], stats["is_complete"])

		if len(path.Steps) > 0 {
			fmt.Printf(", avg_sim=%.2f", stats["average_similarity"])
		}
		fmt.Println()
	}
	fmt.Println()
}

// DetailedPathAnalysis provides in-depth analysis of a single path
func DetailedPathAnalysis() {
	fmt.Println("=== Detailed Path Analysis ===")
	cell := h3.Cell(0x8844d072a3fffff)

	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(15).
		WithMinSimilarity(0.1).
		WithStopOnReturn(true)

	completePath, _ := explorer.Explore(cell)

	fmt.Printf("Starting cell: %s\n", cell.String())
	fmt.Printf("Path completion: %t\n", completePath.IsComplete())
	fmt.Printf("Total steps: %d\n", len(completePath.Steps))

	// Print the path using the new Print method
	fmt.Println()
	completePath.Print()

	// Calculate path statistics
	stats := completePath.CompletionStats()
	fmt.Printf("\nPath Statistics:\n")
	fmt.Printf("  Average similarity: %.3f\n", stats["average_similarity"])
	fmt.Printf("  Min similarity: %.3f\n", stats["min_similarity"])
	fmt.Printf("  Max similarity: %.3f\n", stats["max_similarity"])

	if completePath.IsComplete() {
		fmt.Printf("  Return similarity: %.3f\n", stats["completion_similarity"])
	}
	fmt.Println()
}

// CallbackExploration demonstrates custom exploration logic with callbacks
func CallbackExploration() {
	fmt.Println("=== Simple Exploration Without Callbacks ===")
	cell := h3.Cell(0x8844d072a3fffff)

	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(20).
		WithMinSimilarity(0.85).
		WithStopOnReturn(true)

	path, _ := explorer.Explore(cell)
	stats := path.CompletionStats()

	fmt.Printf("High similarity exploration (threshold 0.85):\n")
	fmt.Printf("  Path length: %d steps\n", stats["path_length"])
	fmt.Printf("  Completed: %t\n", stats["is_complete"])
	if len(path.Steps) > 0 {
		fmt.Printf("  Average similarity: %.2f\n", stats["average_similarity"])
	}
	fmt.Println()
}

// ExtremeCases tests edge cases and unusual exploration scenarios
func ExtremeCases() {
	fmt.Println("=== Extreme Cases & Edge Scenarios ===")
	cell := h3.Cell(0x8844d072a3fffff)

	// Case 1: Very high similarity threshold (might fail to find paths)
	fmt.Println("Case 1: Impossibly high similarity threshold (0.99)")
	explorer1 := h3_ptf.NewExplorer().
		WithMaxSteps(10).
		WithMinSimilarity(0.99).
		WithStopOnReturn(true)

	path1, _ := explorer1.Explore(cell)
	fmt.Printf("  Result: %d steps, complete=%t\n", len(path1.Steps), path1.IsComplete())

	// Case 2: Single step exploration
	fmt.Println("Case 2: Single step exploration")
	explorer2 := h3_ptf.NewExplorer().
		WithMaxSteps(1).
		WithMinSimilarity(0.0).
		WithStopOnReturn(false)

	path2, _ := explorer2.Explore(cell)
	fmt.Printf("  Result: %d steps, complete=%t\n", len(path2.Steps), path2.IsComplete())

	// Case 3: Stop on return disabled
	fmt.Println("Case 3: Stop on return disabled")
	explorer3 := h3_ptf.NewExplorer().
		WithMaxSteps(10).
		WithMinSimilarity(0.0).
		WithStopOnReturn(false)

	path3, _ := explorer3.Explore(cell)
	fmt.Printf("  Result: %d steps, complete=%t\n", len(path3.Steps), path3.IsComplete())

	// Case 4: Very long path with low threshold
	fmt.Println("Case 4: Marathon exploration (50 steps, low threshold)")
	explorer4 := h3_ptf.NewExplorer().
		WithMaxSteps(50).
		WithMinSimilarity(0.1).
		WithStopOnReturn(true)

	path4, _ := explorer4.Explore(cell)
	stats4 := path4.CompletionStats()
	fmt.Printf("  Result: %d steps, complete=%t", len(path4.Steps), path4.IsComplete())
	if len(path4.Steps) > 0 {
		fmt.Printf(", avg_sim=%.2f", stats4["average_similarity"])
	}
	fmt.Println()

	// Case 5: Different starting points with extreme settings
	fmt.Println("Case 5: Multiple extreme starting points")
	extremeCells := []h3.Cell{
		h3.Cell(0x8844d072a3fffff),
		h3.Cell(0x8844d070c9fffff),
		h3.Cell(0x8844d072b5fffff),
	}

	for i, startCell := range extremeCells {
		explorer := h3_ptf.NewExplorer().
			WithMaxSteps(5).
			WithMinSimilarity(0.95). // Very high threshold
			WithStopOnReturn(true)

		path, _ := explorer.Explore(startCell)
		fmt.Printf("  Extreme start %d: %d steps, complete=%t\n",
			i+1, len(path.Steps), path.IsComplete())
	}
	fmt.Println()
}

// EdgeIndexDemo demonstrates how to access edge indices from neighbors
func EdgeIndexDemo() {
	fmt.Println("=== Edge Index Demo ===")
	cell := h3.Cell(0x89318a4f163ffff)

	// Create a tile and get its floor
	tile := h3_ptf.SimplePFT(cell)
	floor := tile.Floor()

	fmt.Printf("Anchor: %s\n", cell.String())
	fmt.Printf("Neighbors with their edge indices:\n")

	for i, neighbor := range floor.Neighbors {
		fmt.Printf("  Edge %d: %s (sim: %.2f)\n",
			neighbor.EdgeIndex, neighbor.Tile.Index.String(), neighbor.Similarity)

		// Verify the edge index matches the array index
		if neighbor.EdgeIndex != i {
			fmt.Printf("    ⚠️  Warning: EdgeIndex (%d) doesn't match array index (%d)\n",
				neighbor.EdgeIndex, i)
		}
	}

	// Now explore and show which edge was chosen
	fmt.Printf("\nExploration path:\n")
	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(5).
		WithMinSimilarity(0.0).
		WithStopOnReturn(true)

	path, _ := explorer.Explore(cell)

	for i, step := range path.Steps {
		fmt.Printf("  Step %d: Chose edge %d to reach %s\n",
			i+1, step.EdgeIndex, step.Index.String())
	}
	fmt.Println()
}
