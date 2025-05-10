package mutation

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func FixedTensionVector(epsilon float64) problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		n := s.Graph.NumVertices
		tangled := s.TangledVertexes()
		if len(tangled) == 0 {
			return m
		}
		u := tangled[rand.IntN(len(tangled))]

		disp := graphplane.VertexPos{X: 0, Y: 0}

		k := math.Sqrt((s.Height*s.Width)/float64(n)) * 0.5

		// repulsive forces
		for v := range n {
			if u == v {
				continue
			}
			dx := s.VertPositions[u].X - s.VertPositions[v].X
			dy := s.VertPositions[u].Y - s.VertPositions[v].Y
			d := math.Hypot(dx, dy) + 1e-9
			force := (k * k) / (d * d)
			disp.X += dx * force
			disp.Y += dy * force
		}

		dx := disp.X
		dy := disp.Y
		displacement := math.Hypot(dx, dy)

		temp := math.Min(s.Height, s.Width) * epsilon
		if displacement > 0 {
			scale := math.Min(displacement, temp) / displacement
			dx *= scale
			dy *= scale
		}

		m.VertPositions[u].X = clamp(m.VertPositions[u].X+dx, 0, s.Width)
		m.VertPositions[u].Y = clamp(m.VertPositions[u].Y+dy, 0, s.Height)

		return m
	}
}
