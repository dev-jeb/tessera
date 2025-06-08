package main

import (
	h3_ptf "tessera/src"

	"github.com/uber/h3-go/v4"
)

func main() {
	explorer := h3_ptf.NewExplorer().
		WithMaxSteps(10).
		WithMinSimilarity(0).
		WithStopOnReturn(true)

	explorer.Explore(h3.Cell(0x8924a56a483ffff))
}
