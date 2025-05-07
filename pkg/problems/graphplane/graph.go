package graphplane

import (
	"fmt"
	"math/rand/v2"

	"github.com/fogleman/delaunay"
)

type Graph struct {
	NumVertices                    int    `json:"numVertices"`
	NumEdges                       int    `json:"numEdges"`
	Edges                          []Edge `json:"edges"`
	cachedMaxPossibleIntersections int
}

type Edge struct {
	From int `json:"from"` // From < To
	To   int `json:"to"`
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

// NewRandomPlanarGraph builds a planar graph using Delaunay triangulation.
func NewRandomPlanarGraph(numVertices int) *Graph {
	// 1) Sample random points in unit square
	pts := make([]delaunay.Point, numVertices)
	for i := range pts {
		pts[i] = delaunay.Point{X: rand.Float64(), Y: rand.Float64()}
	}

	// 2) Compute Delaunay triangulation (planar maximal graph)
	tri, err := delaunay.Triangulate(pts)
	if err != nil {
		panic(err)
	}

	// 3) Extract unique undirected edges
	edgeMap := make(map[[2]int]struct{})
	for ti := 0; ti < len(tri.Triangles); ti += 3 {
		for k := 0; k < 3; k++ {
			a, b := tri.Triangles[ti+k], tri.Triangles[ti+(k+1)%3]
			if a > b {
				a, b = b, a
			}
			edgeMap[[2]int{a, b}] = struct{}{}
		}
	}

	edges := make([]Edge, 0, len(edgeMap))
	for e := range edgeMap {
		edges = append(edges, Edge{From: e[0], To: e[1]})
	}

	// 6) Build and return Graph
	return &Graph{
		NumVertices: numVertices,
		NumEdges:    len(edges),
		Edges:       edges,
	}
}

// MaxPossibleIntersections returns the count of edge pairs
// that could cross (i.e., pairs not sharing a vertex).
func (g *Graph) MaxPossibleIntersections() int {
	if g.cachedMaxPossibleIntersections != 0 {
		return g.cachedMaxPossibleIntersections
	}
	cnt := 0
	for i := range g.Edges {
		e1 := g.Edges[i]
		for j := i + 1; j < len(g.Edges); j++ {
			e2 := g.Edges[j]
			if !sharesVertex(e1, e2) {
				cnt++
			}
		}
	}
	g.cachedMaxPossibleIntersections = cnt
	return cnt
}
