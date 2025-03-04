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

	lg := algos.NewProgressLogger(fmt.Sprintf("%s.jsonl", problem.Name()))
	lg.InitLogging()
	lg.LogProblem(problem)

	algoSolution := problem.AlgorithmicSolution()
	fmt.Println("Algorithm:", algoSolution.Fitness())
	lg.Log(algoSolution)

	targetFitness := algoSolution.Fitness() * 0.99
	// targetFitness := 1.0
	params := sga.SGAParams{
		PopulationSize:       1000,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.5,
		MutationRate:         0.01,
	}
	ga := sga.NewSimpleGeneticAlgorithm(problem, targetFitness, params, lg)
	ga.Run()
}
