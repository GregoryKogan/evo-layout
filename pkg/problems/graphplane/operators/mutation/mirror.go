package mutation

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func Mirror() problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		i := rand.IntN(len(m.VertPositions))
		if rand.Float64() < 0.5 {
			m.VertPositions[i].X = s.Width - m.VertPositions[i].X
		} else {
			m.VertPositions[i].Y = s.Height - m.VertPositions[i].Y
		}
		return m
	}
}
