package h3_ptf

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/uber/h3-go/v4"
)

type Tile struct {
	Index      h3.Cell  `json:"index"`
	Attributes []string `json:"attributes"`
}

type Neighbor struct {
	Tile       Tile    `json:"tile"`
	Similarity float32 `json:"similarity"`
	EdgeIndex  int     `json:"edgeIndex"`
}

type Floor struct {
	Anchor    Tile       `json:"anchor"`
	Neighbors []Neighbor `json:"neighbors"`
}

type Step struct {
	DirectedEdge h3.DirectedEdge `json:"directedEdge"`
	EdgeIndex    int             `json:"edgeIndex"`
	Similarity   float32         `json:"similarity"`
	Index        h3.Cell         `json:"index"`
}

type Path struct {
	Anchor Tile   `json:"anchor"`
	Steps  []Step `json:"steps"`
}

type Explorer struct {
	MaxSteps      int     `json:"maxSteps"`
	MinSimilarity float32 `json:"minSimilarity"`
	StopOnReturn  bool    `json:"stopOnReturn"`
}

func (t Tile) Cardinality() (uint, error) {
	length := len(t.Attributes)
	if length == 0 {
		panic("Tile has no attributes")
	}
	return uint(length), nil
}

func (t Tile) Similarity(other Tile) (float32, error) {
	anchorCardinality, err := t.Cardinality()
	if err != nil {
		panic(err)
	}
	otherCardinality, err := other.Cardinality()
	if err != nil {
		panic(err)
	}
	if anchorCardinality != otherCardinality {
		panic("Tiles have different cardinalities")
	}
	weight := float32(1) / float32(anchorCardinality)
	matches := 0
	for i := range t.Attributes {
		if t.Attributes[i] == other.Attributes[i] {
			matches++
		}
	}
	return float32(math.Floor(float64(matches)*float64(weight)*100) / 100), nil
}

func SimplePFT(index h3.Cell) Tile {
	attributes := make([]string, len(index.String()))
	for i, char := range index.String() {
		attributes[i] = string(char)
	}
	return Tile{
		Index:      index,
		Attributes: attributes,
	}
}

func (t Tile) Neighbors() ([]Neighbor, error) {
	directedEdges, err := t.Index.DirectedEdges()
	if err != nil {
		panic(err)
	}
	neighbors := make([]Neighbor, len(directedEdges))
	for i := range directedEdges {
		edge := directedEdges[i]
		neighbor, err := edge.Destination()
		if err != nil {
			panic(err)
		}
		neighborTile := SimplePFT(neighbor)
		similarity, err := t.Similarity(neighborTile)
		if err != nil {
			panic(err)
		}
		neighbors[i] = Neighbor{
			Tile:       neighborTile,
			Similarity: similarity,
			EdgeIndex:  i,
		}
	}
	return neighbors, nil
}

func (t Tile) Floor() Floor {
	neighbors, err := t.Neighbors()
	if err != nil {
		panic(err)
	}
	return Floor{
		Anchor:    t,
		Neighbors: neighbors,
	}
}

func (t Tile) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (n Neighbor) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (f Floor) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s Step) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p Path) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (e Explorer) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// NewExplorer creates a new explorer with default settings
func NewExplorer() *Explorer {
	return &Explorer{
		MaxSteps:      100,
		MinSimilarity: 0.0,
		StopOnReturn:  false,
	}
}

// WithMaxSteps sets the maximum number of steps the explorer can take
func (e *Explorer) WithMaxSteps(maxSteps int) *Explorer {
	e.MaxSteps = maxSteps
	return e
}

// WithMinSimilarity sets the minimum similarity threshold for continuing exploration
func (e *Explorer) WithMinSimilarity(minSimilarity float32) *Explorer {
	e.MinSimilarity = minSimilarity
	return e
}

// WithStopOnReturn enables stopping exploration if we encounter the starting tile
func (e *Explorer) WithStopOnReturn(stop bool) *Explorer {
	e.StopOnReturn = stop
	return e
}

// Explore traverses the similarity map starting from the given index
func (e *Explorer) Explore(startIndex h3.Cell) (Path, error) {
	startTile := SimplePFT(startIndex)
	path := Path{
		Anchor: startTile,
		Steps:  make([]Step, 0),
	}

	currentTile := startTile
	var previousTile *Tile = nil // Track the previous tile to prevent immediate backtracking

	for len(path.Steps) < e.MaxSteps {
		floor := currentTile.Floor()
		bestNeighbor := floor.FindBestNeighbor()

		// Stop if no neighbors or similarity too low
		if bestNeighbor == nil || bestNeighbor.Similarity < e.MinSimilarity {
			break
		}

		// Check if we should stop on return to start
		if e.StopOnReturn && bestNeighbor.Tile.Index == startIndex {
			// Create final step back to start and stop
			directedEdges, err := currentTile.Index.DirectedEdges()
			if err != nil {
				return path, err
			}

			step := Step{
				DirectedEdge: directedEdges[bestNeighbor.EdgeIndex],
				EdgeIndex:    bestNeighbor.EdgeIndex,
				Similarity:   bestNeighbor.Similarity,
				Index:        startIndex,
			}

			path.Steps = append(path.Steps, step)
			break // Stop exploration - we've returned to start
		}

		// Skip if this neighbor is the tile we just came from (prevent immediate backtracking)
		if previousTile != nil && bestNeighbor.Tile.Index == previousTile.Index {
			// Try to find next best neighbor that isn't the previous tile
			var nextBest *Neighbor
			for i := range floor.Neighbors {
				neighbor := &floor.Neighbors[i]
				if neighbor.Tile.Index != previousTile.Index && neighbor.Similarity >= e.MinSimilarity {
					if nextBest == nil || neighbor.Similarity > nextBest.Similarity {
						nextBest = neighbor
					}
				}
			}
			if nextBest == nil {
				break // No valid neighbors that aren't the previous tile
			}
			bestNeighbor = nextBest
		}

		// Get the directed edge for this neighbor
		directedEdges, err := currentTile.Index.DirectedEdges()
		if err != nil {
			return path, err
		}

		// Create step using the neighbor's edge index
		step := Step{
			DirectedEdge: directedEdges[bestNeighbor.EdgeIndex],
			EdgeIndex:    bestNeighbor.EdgeIndex,
			Similarity:   bestNeighbor.Similarity,
			Index:        bestNeighbor.Tile.Index,
		}

		path.Steps = append(path.Steps, step)

		// Update for next iteration
		previousTile = &currentTile
		currentTile = bestNeighbor.Tile
	}

	return path, nil
}

// FindBestNeighbor returns the neighbor with the highest similarity score
func (f Floor) FindBestNeighbor() *Neighbor {
	if len(f.Neighbors) == 0 {
		return nil
	}

	bestNeighbor := &f.Neighbors[0]
	for i := 1; i < len(f.Neighbors); i++ {
		if f.Neighbors[i].Similarity > bestNeighbor.Similarity {
			bestNeighbor = &f.Neighbors[i]
		}
	}
	return bestNeighbor
}

// IsComplete checks if the path returns to its starting point
func (p Path) IsComplete() bool {
	if len(p.Steps) == 0 {
		return false
	}
	lastStep := p.Steps[len(p.Steps)-1]
	return lastStep.Index == p.Anchor.Index
}

// CompletionStats provides statistics about the completed path
func (p Path) CompletionStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["is_complete"] = p.IsComplete()
	stats["path_length"] = len(p.Steps)

	if len(p.Steps) > 0 {
		totalSimilarity := float32(0)
		minSimilarity := p.Steps[0].Similarity
		maxSimilarity := p.Steps[0].Similarity

		for _, step := range p.Steps {
			totalSimilarity += step.Similarity
			if step.Similarity < minSimilarity {
				minSimilarity = step.Similarity
			}
			if step.Similarity > maxSimilarity {
				maxSimilarity = step.Similarity
			}
		}

		stats["average_similarity"] = totalSimilarity / float32(len(p.Steps))
		stats["min_similarity"] = minSimilarity
		stats["max_similarity"] = maxSimilarity

		if p.IsComplete() {
			stats["completion_similarity"] = p.Steps[len(p.Steps)-1].Similarity
		}
	}

	return stats
}

// Print displays the path progression in a nice, readable format
func (p Path) Print() {
	if len(p.Steps) == 0 {
		fmt.Printf("Path from %s: No steps taken\n", p.Anchor.Index.String())
		return
	}

	fmt.Printf("Exploration path from %s (%d steps):\n", p.Anchor.Index.String(), len(p.Steps))

	for i, step := range p.Steps {
		isReturn := step.Index == p.Anchor.Index
		marker := ""
		if isReturn {
			marker = " (return to start)"
		}

		fmt.Printf("  Step %d: Chose edge %d to reach %s%s\n",
			i+1, step.EdgeIndex, step.Index.String(), marker)
	}

	if p.IsComplete() {
		fmt.Printf("  ✅ Path completed successfully!\n")
	}
}
