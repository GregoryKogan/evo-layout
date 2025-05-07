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
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	problem := graphplane.NewPlanarGraphPlaneProblem(50)
	population := 500
	mutationFunc := mutation.Uniform()
	crossoverFunc := crossover.Uniform(0.45)

	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), "NSGA2+Force"))
	logger := algos.NewProgressLogger(logPath)
	logger.InitLogging()
	logger.LogProblem(problem)

	forceLayer := graphplane.NewForceDirectedLayer(problem.RandomSolution(), 2000, 0.002, logger)
	algoSolution := forceLayer.Solve()

	params := nsga2.NSGA2Params{
		PopulationSize:  population,
		GenerationLimit: math.MaxInt,
		MutationFunc:    mutationFunc,
		CrossoverFunc:   crossoverFunc,
	}
	alg := nsga2.NewAlgorithm(problem, timeLimit, params, logger)
	alg.Seed(algoSolution.Solution)
	runAlgorithm(problem.Name(), alg, timeLimit)
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
