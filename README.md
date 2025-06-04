# Tessera - H3 Geospatial Explorer

Tessera is a Go library for exploring H3 geospatial data using similarity-based pathfinding. It creates "floors" of neighbors with similarity scores and allows you to traverse this spatial map by following the paths of highest similarity. **NEW**: Features a simple flag to stop exploration when returning to the starting point!

## Features

- **H3 Integration**: Built on top of Uber's H3 geospatial indexing system
- **Similarity Scoring**: Each neighbor gets a similarity score based on shared attributes
- **Explorer Engine**: Traverse your geospatial map by always taking the most similar path
- **Path Tracking**: Keep track of the complete journey with Steps and Paths
- **Stop On Return**: Simple flag to stop exploration when encountering the starting tile
- **Flexible Configuration**: Control exploration with max steps, minimum similarity thresholds

## Quick Start

```go
package main

import (
    "fmt"
    h3_ptf "tessera/src"
    "github.com/uber/h3-go/v4"
)

func main() {
    // Start with an H3 cell
    cell := h3.Cell(0x8844d072a3fffff)
    
    // Create an explorer that stops when returning to start
    explorer := h3_ptf.NewExplorer().
        WithMaxSteps(15).
        WithMinSimilarity(0.8).
        WithStopOnReturn(true)
    
    // Explore and possibly complete the path
    path, err := explorer.Explore(cell)
    if err != nil {
        panic(err)
    }
    
    // Check if the path completed
    fmt.Printf("Path completed: %t\n", path.IsComplete())
    fmt.Printf("Path length: %d steps\n", len(path.Steps))
    
    // Get completion statistics
    stats := path.CompletionStats()
    fmt.Printf("Average similarity: %.2f\n", stats["average_similarity"])
    if path.IsComplete() {
        fmt.Printf("Return similarity: %.2f\n", stats["completion_similarity"])
    }
}
```

## Explorer Configuration

The Explorer type supports several configuration options:

- `WithMaxSteps(n)`: Maximum number of steps to take
- `WithMinSimilarity(threshold)`: Minimum similarity score to continue
- `WithStopOnReturn(true/false)`: Stop exploration if we encounter the starting tile

## Stop On Return Behavior

When `StopOnReturn` is enabled:

- Explorer continues normal pathfinding, always choosing the highest similarity neighbor
- If the best neighbor happens to be the starting tile, exploration stops and the path is completed
- This creates natural loops when the similarity map leads back to the origin

When `StopOnReturn` is disabled:

- Explorer continues until max steps or no valid neighbors remain
- Starting tile is treated like any other tile

**Note**: The explorer prevents immediate backtracking to the tile you just came from, but allows revisiting tiles from earlier in the path.

## Advanced Usage

### Basic Stop On Return

```go
// Enable stopping on return to start
explorer := h3_ptf.NewExplorer().
    WithStopOnReturn(true)

// Disable stopping on return (default behavior)
explorer := h3_ptf.NewExplorer().
    WithStopOnReturn(false)
```

### Analyzing Completed Paths

```go
path, _ := explorer.Explore(cell)

// Print the path nicely
path.Print()

// Check completion status
if path.IsComplete() {
    fmt.Println("✅ Path completed successfully!")
    
    // Get detailed statistics
    stats := path.CompletionStats()
    fmt.Printf("Average similarity: %.3f\n", stats["average_similarity"])
    fmt.Printf("Return similarity: %.3f\n", stats["completion_similarity"])
    fmt.Printf("Path length: %d\n", stats["path_length"])
}
```

### Multiple Starting Points

```go
cells := []h3.Cell{
    h3.Cell(0x8844d072a3fffff),
    h3.Cell(0x8844d070c9fffff),
    h3.Cell(0x8844d072b5fffff),
}

explorer := h3_ptf.NewExplorer().WithMaxSteps(8).WithStopOnReturn(true)

for _, startCell := range cells {
    path, _ := explorer.Explore(startCell)
    fmt.Printf("Path from %s: %d steps, complete: %t\n", 
        startCell.String(), len(path.Steps), path.IsComplete())
}
```

## Running Examples

The project includes a comprehensive scenario system for exploring different use cases:

```bash
# Quick start - run the basic demo
go run main.go

# Interactive scenario menu
cd examples
go run .

# Run specific scenarios
go run . 1    # Basic Exploration
go run . 2    # Stop On Return Demo  
go run . 3    # Multiple Starting Points
go run . 4    # Similarity Threshold Comparison
go run . 5    # Detailed Path Analysis
go run . 6    # Simple High-Threshold Exploration
go run . 7    # Extreme Cases & Edge Scenarios
go run . 8    # Edge Index Demo

# Run all scenarios
go run . all
```

## Adding New Scenarios

To add a new exploration scenario, simply:

1. Add a new function to `examples/scenarios.go`:

```go
func YourNewScenario() {
    fmt.Println("=== Your New Scenario ===")
    // Your exploration logic here
}
```

2. Update `examples/main.go` to include it:

```go
scenarios := map[string]func(){
    // ... existing scenarios ...
    "8": YourNewScenario,
}
```

3. Add it to the menu display

This modular approach makes it easy to experiment with different exploration patterns and share them with others.

## Data Structures

- **Tile**: Represents an H3 cell with attributes for similarity calculation
- **Neighbor**: A neighboring tile with similarity score and edge index (0-5) indicating its position relative to the anchor
- **Floor**: Contains an anchor tile and its neighbors with similarity scores
- **Step**: A single move in the exploration with direction, similarity, destination, and neighbor index (0-5)
- **Path**: Complete journey from start to finish with all steps
- **Explorer**: The engine that traverses the similarity map

## Path Analysis Methods

- `path.IsComplete()`: Returns true if the path forms a closed loop
- `path.CompletionStats()`: Returns detailed statistics about the path including averages, min/max similarities, and completion status
- `path.Print()`: Displays the path progression in a nice, readable format with step-by-step details including which neighbor (0-5) was chosen

### Example Output

```
Exploration path from 8944d072827ffff (3 steps):
  Step 1: Chose edge 1 to reach 8944d07282fffff
  Step 2: Chose edge 3 to reach 8944d072823ffff
  Step 3: Chose edge 0 to reach 8944d072827ffff (return to start)
  ✅ Path completed successfully!
```

The neighbor index (0-5) indicates which of the 6 H3 neighbors was selected at each step, providing insight into the exploration patterns.
