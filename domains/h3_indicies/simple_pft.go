/**
Probably the simplest PFT that I can think of
*/

package h3_indicies

type H3IndexTile struct {
	Attributes []string
}

type H3IndexCodomain struct {
	Elements []H3IndexTile
}

func NewSimplePFT(domain H3IndexDomain) H3IndexCodomain {
	tiles := make([]H3IndexTile, len(domain.Elements))
	for i, cell := range domain.Elements {
		indexStr := cell.String()
		attributes := make([]string, len(indexStr))
		for j, char := range indexStr {
			attributes[j] = string(char)
		}
		tiles[i] = H3IndexTile{Attributes: attributes}
	}
	return H3IndexCodomain{
		Elements: tiles,
	}
}
