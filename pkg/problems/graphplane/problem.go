package graphplane

import (
	"fmt"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneProblem struct {
	Graph         *Graph
	Width, Height float64
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

func (p *GraphPlaneProblem) MarshalJSON() ([]byte, error) {
	graphJSON, err := p.Graph.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return []byte(`{"Name":"GraphPlane","Width":` + fmt.Sprint(p.Width) + `,"Height":` + fmt.Sprint(p.Height) + `,"Graph":` + string(graphJSON) + `}`), nil
}
