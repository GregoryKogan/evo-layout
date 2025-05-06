package operators

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func NormWeightedMutation() problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)
		// Intersection-based weights
		weights := intersectionWeights(s)
		for i := range m.VertPositions {
			// base probability plus weighted factor
			p := weights[i]/(weights[i]+1) + 0.1
			if rand.Float64() < p {
				dx := rand.NormFloat64() * s.Width
				dy := rand.NormFloat64() * s.Height
				m.VertPositions[i].X = clamp(m.VertPositions[i].X+dx, 0, s.Width)
				m.VertPositions[i].Y = clamp(m.VertPositions[i].Y+dy, 0, s.Height)
			}
		}
		return m
	}
}

// intersectionWeights returns per-vertex count of incident crossings.
func intersectionWeights(s *graphplane.GraphPlaneSolution) []float64 {
	w := make([]float64, len(s.VertPositions))
	for i := range s.Graph.Edges {
		e1 := s.Graph.Edges[i]
		p1, p2 := s.VertPositions[e1.From], s.VertPositions[e1.To]
		for j := i + 1; j < len(s.Graph.Edges); j++ {
			e2 := s.Graph.Edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			p3, p4 := s.VertPositions[e2.From], s.VertPositions[e2.To]
			if segmentsIntersect(p1, p2, p3, p4) {
				w[e1.From]++
				w[e1.To]++
				w[e2.From]++
				w[e2.To]++
			}
		}
	}
	return w
}

// sharesVertex checks if two edges share an endpoint.
func sharesVertex(e1, e2 graphplane.Edge) bool {
	return e1.From == e2.From || e1.From == e2.To || e1.To == e2.From || e1.To == e2.To
}

// segmentsIntersect tests segment intersection via CCW.
func segmentsIntersect(a, b, c, d graphplane.VertexPos) bool {
	ccw := func(u, v, w graphplane.VertexPos) bool {
		return (w.Y-u.Y)*(v.X-u.X) > (v.Y-u.Y)*(w.X-u.X)
	}
	return ccw(a, c, d) != ccw(b, c, d) && ccw(a, b, c) != ccw(a, b, d)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
