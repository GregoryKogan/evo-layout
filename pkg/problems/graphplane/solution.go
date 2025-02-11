package graphplane

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneSolution struct {
	graph         *Graph
	width, height float64
	vertPositions []VertexPos
}

type VertexPos struct {
	X, Y float64
}

func RandomGraphPlaneSolution(g *Graph, width, height float64) problems.Solution {
	numVertices := g.NumVertices
	vertPositions := make([]VertexPos, numVertices)
	for i := range vertPositions {
		vertPositions[i].X = rand.Float64() * width
		vertPositions[i].Y = rand.Float64() * height
	}
	return &GraphPlaneSolution{g, width, height, vertPositions}
}

func (s *GraphPlaneSolution) Crossover(other problems.Solution) problems.Solution {
	otherGPS, ok := other.(*GraphPlaneSolution)
	if !ok {
		return s
	}
	child := &GraphPlaneSolution{
		vertPositions: make([]VertexPos, s.graph.NumVertices),
		width:         s.width,
		height:        s.height,
		graph:         s.graph, // assume same graph
	}
	// For each vertex choose allele from one of the parents.
	for i := 0; i < s.graph.NumVertices; i++ {
		if rand.Float64() < 0.5 {
			child.vertPositions[i] = s.vertPositions[i]
		} else {
			child.vertPositions[i] = otherGPS.vertPositions[i]
		}
	}
	return child
}

func (s *GraphPlaneSolution) Mutate() problems.Solution {
	mutant := &GraphPlaneSolution{
		vertPositions: make([]VertexPos, s.graph.NumVertices),
		width:         s.width,
		height:        s.height,
		graph:         s.graph,
	}
	copy(mutant.vertPositions, s.vertPositions)
	// pick a random vertex and change its position by a small delta
	i := rand.Intn(s.graph.NumVertices)
	deltaX := (rand.Float64() - 0.5) * s.width * 0.05
	deltaY := (rand.Float64() - 0.5) * s.height * 0.05
	mutant.vertPositions[i].X = math.Max(0, math.Min(s.width, mutant.vertPositions[i].X+deltaX))
	mutant.vertPositions[i].Y = math.Max(0, math.Min(s.height, mutant.vertPositions[i].Y+deltaY))
	return mutant
}

func (s *GraphPlaneSolution) Fitness() float64 {
	// count intersections between all pairs of non-adjacent edges
	count := 0
	edges := s.graph.Edges
	// helper function for ccw test
	ccw := func(ax, ay, bx, by, cx, cy float64) bool {
		return (cy-ay)*(bx-ax) > (by-ay)*(cx-ax)
	}
	segmentsIntersect := func(ax, ay, bx, by, cx, cy, dx, dy float64) bool {
		return ccw(ax, ay, cx, cy, dx, dy) != ccw(bx, by, cx, cy, dx, dy) &&
			ccw(ax, ay, bx, by, cx, cy) != ccw(ax, ay, bx, by, dx, dy)
	}
	n := len(edges)
	for i := 0; i < n; i++ {
		a := edges[i].From
		b := edges[i].To
		ax, ay := s.vertPositions[a].X, s.vertPositions[a].Y
		bx, by := s.vertPositions[b].X, s.vertPositions[b].Y
		for j := i + 1; j < n; j++ {
			// skip if edges share a vertex
			if edges[i].From == edges[j].From ||
				edges[i].From == edges[j].To ||
				edges[i].To == edges[j].From ||
				edges[i].To == edges[j].To {
				continue
			}
			c := edges[j].From
			d := edges[j].To
			cx, cy := s.vertPositions[c].X, s.vertPositions[c].Y
			dx, dy := s.vertPositions[d].X, s.vertPositions[d].Y
			if segmentsIntersect(ax, ay, bx, by, cx, cy, dx, dy) {
				count++
			}
		}
	}
	// Higher fitness for fewer intersections.
	return 1.0 / (1.0 + float64(count))
}

func (s *GraphPlaneSolution) MarshalJSON() ([]byte, error) {
	vertPosJSON, err := json.Marshal(s.vertPositions)
	if err != nil {
		return nil, err
	}
	return []byte(`{"Vertices":` + string(vertPosJSON) + `}`), nil
}
