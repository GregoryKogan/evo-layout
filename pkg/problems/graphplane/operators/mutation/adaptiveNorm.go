package mutation

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func AdaptiveNorm(maxSteps int) problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}

		tangled := s.TangledVertexes()
		if len(tangled) > 0 {
			return FixedNorm()(s)
		}

		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		for range maxSteps {
			i := rand.IntN(len(m.VertPositions))

			oldX := m.VertPositions[i].X
			oldY := m.VertPositions[i].Y

			dx := rand.NormFloat64() * s.Width / 50
			dy := rand.NormFloat64() * s.Height / 50
			m.VertPositions[i].X = clamp(m.VertPositions[i].X+dx, 0, s.Width)
			m.VertPositions[i].Y = clamp(m.VertPositions[i].Y+dy, 0, s.Height)

			if m.CountIntersections() == 0 {
				return m
			}

			m.VertPositions[i].X = oldX
			m.VertPositions[i].Y = oldY
		}

		return m
	}
}
