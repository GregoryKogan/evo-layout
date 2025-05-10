package mutation

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func FixedPercentage() problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		tangled := s.TangledVertexes()
		if len(tangled) == 0 {
			return m
		}
		i := tangled[rand.IntN(len(tangled))]
		if rand.Float64() < 0.5 {
			m.VertPositions[i].X = clamp(m.VertPositions[i].X*(0.8+rand.Float64()*0.4), 0, s.Width)
		} else {
			m.VertPositions[i].Y = clamp(m.VertPositions[i].Y*(0.8+rand.Float64()*0.4), 0, s.Height)
		}

		return m
	}
}
