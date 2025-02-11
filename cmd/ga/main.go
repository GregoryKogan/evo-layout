package main

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func main() {
	numVertices := 10
	edgeFill := 0.5
	width, height := 1.0, 1.0
	problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, width, height)

	targetFitness := 0.9
	populationSize := 1000
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, populationSize)
	ga.Run()
}
