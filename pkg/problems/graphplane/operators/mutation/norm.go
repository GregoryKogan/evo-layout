package mutation

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func Norm(k float64) problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		i := rand.IntN(len(m.VertPositions))
		dx := rand.NormFloat64() * s.Width * k
		dy := rand.NormFloat64() * s.Height * k
		m.VertPositions[i].X = clamp(m.VertPositions[i].X+dx, 0, s.Width)
		m.VertPositions[i].Y = clamp(m.VertPositions[i].Y+dy, 0, s.Height)

		return m
	}
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
