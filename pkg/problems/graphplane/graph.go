package graphplane

import (
	"fmt"
	"math/rand/v2"
)

type Graph struct {
	NumVertices int
	NumEdges    int
	Edges       []Edge
}

type Edge struct {
	From, To int
}

func NewRandomGraph(numVertices int, edgeFill float64) *Graph {
	if edgeFill < 0 || edgeFill > 1 {
		panic("edgeFill must be in the range [0, 1]")
	}
	maxEdges := numVertices * (numVertices - 1) / 2
	numEdges := int(edgeFill * float64(maxEdges))

	edges := make([]Edge, 0, numEdges)
	edgeSet := make(map[string]bool)
	for len(edges) < numEdges {
		u := rand.IntN(numVertices)
		v := rand.IntN(numVertices)
		if u == v {
			continue // skip self-loops
		}
		// Normalize edge representation.
		small, large := u, v
		if u > v {
			small, large = v, u
		}
		key := fmt.Sprintf("%d-%d", small, large)
		if edgeSet[key] {
			continue // skip parallel edges
		}
		edgeSet[key] = true
		edges = append(edges, Edge{From: small, To: large})
	}

	return &Graph{NumVertices: numVertices, NumEdges: numEdges, Edges: edges}
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	return []byte(`{"NumVertices":` + fmt.Sprint(g.NumVertices) + `,"NumEdges":` + fmt.Sprint(g.NumEdges) + `,"Edges":` + fmt.Sprint(g.Edges) + `}`), nil
}
