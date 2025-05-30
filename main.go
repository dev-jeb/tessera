package main

import (
	"fmt"
	"tessera/domains/h3_indicies"
)

func main() {
	domain := h3_indicies.NewH3IndexDomain()
	for i := 0; i < 10; i++ {
		fmt.Printf("Element %d: %v\n", i, domain.Elements[i])
	}
	codomain := h3_indicies.NewSimplePFT(domain)

	fmt.Printf("Created %d tiles\n", len(codomain.Elements))

	// Show first 3 tiles
	for i := 0; i < 3 && i < len(codomain.Elements); i++ {
		tile := codomain.Elements[i]
		fmt.Printf("\nTile %d:\n", i)
		fmt.Printf("  Attributes: %v\n", tile.Attributes)
	}
}
