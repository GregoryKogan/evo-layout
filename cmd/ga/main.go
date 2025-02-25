package main

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/knapsack"
)

func main() {
	// numVertices := 12
	// edgeFill := 0.35
	// problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, 1.0, 1.0)
	problem := knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
		Dimensions:         2,
		ItemsNum:           100,
		InitialMaxValue:    100,
		InitialMaxResource: 100,
		InitialMaxAmount:   10,
		Constraints:        []int{3000},
	})

	targetFitness := 1_000_000_000.0
	params := sga.SGAParams{
		PopulationSize:       10000,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.5,
		MutationRate:         0.05,
	}
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, "knapsack.jsonl", params)
	ga.Run()
}
