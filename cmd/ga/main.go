package main

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func main() {
	numVertices := 13
	edgeFill := 0.35
	width, height := 1.0, 1.0
	problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, width, height)

	targetFitness := 1.0
	populationSize := 10000
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, populationSize)
	ga.Run()
}
