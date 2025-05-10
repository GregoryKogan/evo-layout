package mutation

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func FixedNorm() problems.MutationFunc {
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
		dx := rand.NormFloat64() * s.Width / 10
		dy := rand.NormFloat64() * s.Height / 10
		m.VertPositions[i].X = clamp(m.VertPositions[i].X+dx, 0, s.Width)
		m.VertPositions[i].Y = clamp(m.VertPositions[i].Y+dy, 0, s.Height)

		return m
	}
}
