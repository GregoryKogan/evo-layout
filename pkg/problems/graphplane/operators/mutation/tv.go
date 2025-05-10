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

		n := s.Graph.NumVertices
		u := rand.IntN(n)

		var neighbors []int
		for _, e := range m.Graph.Edges {
			if e.From == u {
				neighbors = append(neighbors, e.To)
			} else if e.To == u {
				neighbors = append(neighbors, e.From)
			}
		}

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

		// attractive forces along edges
		for _, v := range neighbors {
			dx := m.VertPositions[u].X - m.VertPositions[v].X
			dy := m.VertPositions[u].Y - m.VertPositions[v].Y
			d := math.Hypot(dx, dy) + 1e-9
			force := (d * d) / k
			dxNorm := dx / d
			dyNorm := dy / d
			disp.X -= dxNorm * force
			disp.Y -= dyNorm * force
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
