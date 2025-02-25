package main

import (
	"fmt"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/knapsack"
)

func main() {
	// numVertices := 12
	// edgeFill := 0.35
	// problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, 1.0, 1.0)
	problem := knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
		Dimensions:         2,
		ItemsNum:           10000,
		InitialMaxValue:    100,
		InitialMaxResource: 100,
		// Constraints:        []int{300000, 300000, 300000, 300000, 300000, 300000, 300000, 300000, 300000},
		Constraints: []int{300000},
	})

	algoSolution := problem.AlgorithmicSolution()
	fmt.Println(algoSolution.Fitness())

	lg := algos.NewProgressLogger("algo-solution.jsonl")
	lg.InitLogging()
	lg.Log(algoSolution)

	targetFitness := algoSolution.Fitness() * 0.95
	params := sga.SGAParams{
		PopulationSize:       1000,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.5,
		MutationRate:         0.0001,
	}
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, "knapsack.jsonl", params)
	ga.Run()
}
