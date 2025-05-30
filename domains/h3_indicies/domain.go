package h3_indicies

import (
	"github.com/uber/h3-go/v4"
)

type H3IndexDomain struct {
	Elements []h3.Cell
}

func (d H3IndexDomain) Cardinality() int {
	return len(d.Elements)
}

func NewH3IndexDomain() H3IndexDomain {
	return H3IndexDomain{
		Elements: generateH3Indices(),
	}
}

func generateH3Indices() []h3.Cell {
	baseCells, err := h3.Res0Cells()
	if err != nil {
		panic(err)
	}
	return baseCells
}
