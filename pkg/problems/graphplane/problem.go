package graphplane

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneProblem struct {
	Graph  *Graph  `json:"graph"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func NewGraphPlaneProblem(numVertices int, edgeFill float64, width, height float64) problems.Problem {
	return &GraphPlaneProblem{NewRandomGraph(numVertices, edgeFill), width, height}
}

func (p *GraphPlaneProblem) Name() string {
	return "GraphPlane"
}

func (p *GraphPlaneProblem) RandomSolution() problems.Solution {
	return RandomGraphPlaneSolution(p.Graph, p.Width, p.Height)
}
