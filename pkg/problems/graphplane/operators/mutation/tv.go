package mutation

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func TensionVector(epsilon float64) problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*graphplane.GraphPlaneSolution)
		if !ok {
			panic("invalid individual")
		}
		m := &graphplane.GraphPlaneSolution{Graph: s.Graph, Width: s.Width, Height: s.Height}
		m.VertPositions = make([]graphplane.VertexPos, len(s.VertPositions))
		copy(m.VertPositions, s.VertPositions)

		u := rand.IntN(len(m.VertPositions))

		var neighbors []int
		for _, e := range m.Graph.Edges {
			if e.From == u {
				neighbors = append(neighbors, e.To)
			} else if e.To == u {
				neighbors = append(neighbors, e.From)
			}
		}

		disp := graphplane.VertexPos{X: 0, Y: 0}

		// attractive forces along edges
		springK := math.Min(s.Width, s.Height) / math.Sqrt(float64(s.Graph.NumVertices))
		for _, v := range neighbors {
			dx := m.VertPositions[u].X - m.VertPositions[v].X
			dy := m.VertPositions[u].Y - m.VertPositions[v].Y
			d := math.Hypot(dx, dy) + 1e-9
			force := (d * d) / springK
			dxNorm := dx / d
			dyNorm := dy / d
			disp.X -= dxNorm * force
			disp.Y -= dyNorm * force
		}

		px := m.VertPositions[u].X + disp.X*epsilon
		py := m.VertPositions[u].Y + disp.Y*epsilon
		m.VertPositions[u].X = clamp(px, 0, s.Width)
		m.VertPositions[u].Y = clamp(py, 0, s.Height)

		return m
	}
}
