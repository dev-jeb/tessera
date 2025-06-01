package main

import (
	"fmt"
	h3_ptf "tessera/src"

	"github.com/uber/h3-go/v4"
)

func main() {
	cell := h3.Cell(0x8844d072a3fffff)
	anchor := h3_ptf.SimplePFT(cell)
	floor := anchor.Floor()
	floorJSON, err := floor.ToJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(floorJSON)
}
