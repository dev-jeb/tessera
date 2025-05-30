/*
Package tessera provides a formalization of procedural tile functions (PTF)
over a finite domain. I have chosen to use Go to formalize the definitions
because my co worker has been trying to get me to write some Go for
a while. GBU.
**/

package core

// The domain of a PTMF is a finite set of elements.
type Domain struct {
	Elements []any
}

// The cardinality of a domain is the number of elements it has.
func (d Domain) Cardinality() int {
	return len(d.Elements)
}

// A tile is a structured object with a finite set of ordered attributes.
type Tile struct {
	Attributes []any
}

// The cardinality of a tile is the number of attributes it has.
func (t Tile) Cardinality() int {
	return len(t.Attributes)
}

// two tiles are said to be equal if they have the same cardinality and the same
// attributes in the same order.
func (t Tile) Equals(other Tile) bool {
	if t.Cardinality() != other.Cardinality() {
		return false
	}
	for i := range t.Attributes {
		if t.Attributes[i] != other.Attributes[i] {
			return false
		}
	}
	return true
}

// Now we have the ability to fomalize the Codomain. The codomain of a PTF is a
// finite set of tiles.
type Codomain struct {
	Elements []Tile
}

// A PTF is a deterministic function that takes an element from the domain and
// returns a tile from the codomain.
type ProceduralTileFunction func(any) Tile

// IsDeterministic verifies that the function produces the same output for the
// same input by testing it multiple times with the same input.
func (f ProceduralTileFunction) IsDeterministic(input any, numTests int) bool {
	if numTests < 2 {
		return true
	}

	firstResult := f(input)
	for i := 1; i < numTests; i++ {
		if !firstResult.Equals(f(input)) {
			return false
		}
	}
	return true
}

// A PTF is said to be ONTO if for every tile in the codomain there is an element in the domain that maps to it.
func (f ProceduralTileFunction) IsOnto(domain Domain, codomain Codomain) bool {
	for _, tile := range codomain.Elements {
		found := false
		for _, element := range domain.Elements {
			if f(element).Equals(tile) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// A PTF is said to be ONE-TO-ONE if for every tile in the codomain there is one and only one element in the domain that maps to it.
func (f ProceduralTileFunction) IsOneToOne(domain Domain, codomain Codomain) bool {
	// the function must be onto to be one-to-one
	if !f.IsOnto(domain, codomain) {
		return false
	}

	// now we need to check that each tile in the codomain is mapped to by exactly one element in the domain
	for _, tile := range codomain.Elements {
		count := 0
		for _, element := range domain.Elements {
			if f(element).Equals(tile) {
				count++
				break
			}
		}
		if count != 1 {
			return false
		}
	}
	return true
}

// A PTF is said to be BIJECTIVE if it is both one-to-one and onto.
func (f ProceduralTileFunction) IsBijective(domain Domain, codomain Codomain) bool {
	return f.IsOneToOne(domain, codomain) && f.IsOnto(domain, codomain)
}
