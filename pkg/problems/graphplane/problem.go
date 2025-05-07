package graphplane

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneProblem struct {
	name   string
	Graph  *Graph  `json:"graph"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func NewGraphPlaneProblem(numVertices int, edgeFill float64) problems.Problem {
	return &GraphPlaneProblem{"GraphPlane", NewRandomGraph(numVertices, edgeFill), 1.0, 1.0}
}

func NewPlanarGraphPlaneProblem(numVertices int) problems.Problem {
	return &GraphPlaneProblem{"PlanarGraphPlane", NewRandomPlanarGraph(numVertices), 1.0, 1.0}
}

func (p *GraphPlaneProblem) Name() string {
	return p.name
}

func (p *GraphPlaneProblem) RandomSolution() problems.Solution {
	return RandomGraphPlaneSolution(p.Graph, p.Width, p.Height)
}
