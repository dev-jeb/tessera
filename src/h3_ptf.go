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
	Tile        Tile    `json:"tile"`
	Similarity  float32 `json:"similarity"`
	HouseNumber int     `json:"edgeIndex"`
}
type Explorer struct {
	StepsTaken    int     `json:"stepsTaken"`
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
			Tile:        neighborTile,
			Similarity:  similarity,
			HouseNumber: i,
		}
	}
	return neighbors, nil
}

func NewExplorer() *Explorer {
	return &Explorer{
		StepsTaken:    0,
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
func (e *Explorer) Explore(startIndex h3.Cell) error {
	visited := make(map[h3.Cell]int)
	startTile := SimplePFT(startIndex)
	visited[startTile.Index] = 1
	fmt.Printf("%s\n", startTile.Index.String())

	currentTile := &startTile
	var previousTile *Tile = nil

	for e.StepsTaken < e.MaxSteps {
		neighbors, err := currentTile.Neighbors()
		if err != nil {
			return err
		}
		var bestNeighbor *Neighbor = nil

		for _, neighbor := range neighbors {
			if neighbor.Similarity < e.MinSimilarity {
				continue
			}
			if previousTile != nil && neighbor.Tile.Index == (*previousTile).Index {
				continue
			}
			if bestNeighbor == nil {
				bestNeighbor = &neighbor
				continue
			}
			if neighbor.Similarity <= (*bestNeighbor).Similarity {
				bestNeighbor = &neighbor
				continue
			}
		}

		if bestNeighbor == nil {
			fmt.Printf("No valid neighbors\n")
			break
		} else if bestNeighbor.Tile.Index == startTile.Index {
			fmt.Printf("Returned to start\n")
			break
		} else {
			fmt.Printf("%s\n", bestNeighbor.Tile.Index.String())
		}

		previousTile = currentTile
		currentTile = &bestNeighbor.Tile
		e.StepsTaken++
		visited[currentTile.Index]++
	}
	fmt.Printf("Visited %d unique tiles\n", len(visited))

	return nil
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

func (e Explorer) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
