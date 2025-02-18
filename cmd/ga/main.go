package main

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func main() {
	numVertices := 20
	edgeFill := 0.35
	width, height := 1.0, 1.0
	problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, width, height)

	targetFitness := 1.0
	params := sga.SGAParams{
		PopulationSize:       1000,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.9,
		MutationRate:         0.3,
	}
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, "", params)
	ga.Run()
}
