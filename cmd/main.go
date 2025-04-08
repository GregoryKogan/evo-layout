package main

import (
	"fmt"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/ssga"
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
	problem := tsp.NewTSProblem(tsp.TSProblemParameters{CitiesNum: 100})

	lg := algos.NewProgressLogger(fmt.Sprintf("%s.jsonl", problem.Name()))
	lg.InitLogging()
	lg.LogProblem(problem)

	targetFitness := 1.0
	params := ssga.Params{
		PopulationSize: 3000,
		MutationRate:   0.02,
	}
	ga := ssga.NewAlgorithm(problem, targetFitness, params, lg)
	ga.Run()
}
