package h3_ptf

import (
	"encoding/json"
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
}

type Floor struct {
	Anchor    Tile       `json:"anchor"`
	Neighbors []Neighbor `json:"neighbors"`
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
