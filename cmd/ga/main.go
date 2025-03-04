package main

import (
	"fmt"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/tsp"
)

func main() {
	// graphplane
	// numVertices := 12
	// edgeFill := 0.35
	// problem := graphplane.NewGraphPlaneProblem(numVertices, edgeFill, 1.0, 1.0)

	// knapsack
	// problem := knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
	// 	Dimensions:         2,
	// 	ItemsNum:           10000,
	// 	InitialMaxValue:    100,
	// 	InitialMaxResource: 100,
	// 	// Constraints:        []int{300000, 300000, 300000, 300000, 300000, 300000, 300000, 300000, 300000},
	// 	Constraints: []int{300000},
	// })

	// TSP
	problem := tsp.NewTSProblem(tsp.TSProblemParameters{CitiesNum: 12})

	algoSolution := problem.AlgorithmicSolution()
	fmt.Println("Algorithm:", algoSolution.Fitness())

	lg := algos.NewProgressLogger("algo-solution.jsonl")
	lg.InitLogging()
	lg.Log(algoSolution)

	targetFitness := algoSolution.Fitness() * 0.99
	params := sga.SGAParams{
		PopulationSize:       100,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.5,
		MutationRate:         0.1,
	}
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, "tsp.jsonl", params)
	ga.Run()
}
