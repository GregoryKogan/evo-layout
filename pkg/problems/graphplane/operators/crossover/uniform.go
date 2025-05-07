package crossover

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func Uniform(swapProb float64) problems.CrossoverFunc {
	return func(parentA, parentB problems.Solution) []problems.Solution {
		a, aOk := parentA.(*graphplane.GraphPlaneSolution)
		b, bOk := parentB.(*graphplane.GraphPlaneSolution)
		if !aOk || !bOk || len(a.VertPositions) != len(b.VertPositions) {
			panic("invalid parents")
		}

		c1 := &graphplane.GraphPlaneSolution{Graph: a.Graph, Width: a.Width, Height: a.Height}
		c2 := &graphplane.GraphPlaneSolution{Graph: a.Graph, Width: a.Width, Height: a.Height}
		c1.VertPositions = make([]graphplane.VertexPos, len(a.VertPositions))
		c2.VertPositions = make([]graphplane.VertexPos, len(a.VertPositions))
		for i := range a.VertPositions {
			if rand.Float64() < swapProb {
				c1.VertPositions[i] = b.VertPositions[i]
				c2.VertPositions[i] = a.VertPositions[i]
			} else {
				c1.VertPositions[i] = a.VertPositions[i]
				c2.VertPositions[i] = b.VertPositions[i]
			}
		}
		return []problems.Solution{c1, c2}
	}
}
