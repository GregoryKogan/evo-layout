package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 70 * time.Second

	problem := graphplane.NewGraphPlaneProblem(25, 0.35, 1.0, 1.0)
	population := 1000
	mutationProb := 0.1
	crossoverProb := 0.9

	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), "SGA+Force"))
	logger := algos.NewProgressLogger(logPath)
	logger.InitLogging()
	logger.LogProblem(problem)

	logPath2 := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), "Force"))
	logger2 := algos.NewProgressLogger(logPath2)
	logger2.InitLogging()
	logger2.LogProblem(problem)

	// SGA
	// params := sga.Params{
	// 	PopulationSize:       population,
	// 	ElitePercentile:      0.1,
	// 	MatingPoolPercentile: 0.5,
	// 	MutationProb:         mutationProb,
	// 	CrossoverProb:        crossoverProb,
	// }
	// alg := sga.NewAlgorithm(problem, timeLimit, params, logger)

	// SSGA
	// params := ssga.Params{
	// 	PopulationSize: population,
	// 	MutationProb:   mutationProb,
	// 	CrossoverProb:  crossoverProb,
	// }
	// alg := ssga.NewAlgorithm(problem, timeLimit, params, logger)

	// NSGA2
	params := nsga2.NSGA2Params{
		PopulationSize:  population,
		GenerationLimit: math.MaxInt,
		MutationProb:    mutationProb,
		CrossoverProb:   crossoverProb,
	}
	alg := nsga2.NewAlgorithm(problem, timeLimit, params, logger)

	// SPEA2
	// params := spea2.Params{
	// 	PopulationSize:  population,
	// 	ArchiveSize:     population,
	// 	DensityKth:      int(math.Sqrt(float64(population + population))), // typical choice: sqrt(population+archive)
	// 	GenerationLimit: math.MaxInt,
	// 	MutationProb:    mutationProb,
	// 	CrossoverProb:   crossoverProb,
	// }
	// alg := spea2.NewAlgorithm(problem, timeLimit, params, logger)

	solution := runAlgorithm(problem.Name(), alg, timeLimit)

	forceLayer := graphplane.NewForceDirectedLayer(solution, 2000, 0.002, logger)
	forceLayer.Solve()

	forceLayer2 := graphplane.NewForceDirectedLayer(problem.RandomSolution(), 2000, 0.002, logger2)
	forceLayer2.Solve()
}

// runAlgorithm wraps an algorithm run in a goroutine and enforces a timeout.
func runAlgorithm(problemName string, alg *nsga2.Algorithm, timeout time.Duration) problems.Solution {
	fmt.Printf("Running %s with timeout %v\n", problemName, timeout)
	done := make(chan bool)
	go func() {
		alg.Run()
		done <- true
	}()
	select {
	case <-done:
		fmt.Printf("%s completed.\n", problemName)
	case <-time.After(timeout):
		fmt.Printf("%s timed out after %v.\n", problemName, timeout)
	}
	return alg.Solution
}
