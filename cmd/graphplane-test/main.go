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
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 10 * time.Minute

	problem := graphplane.NewGraphPlaneProblem(100, 0.1, 1.0, 1.0)
	population := 1000
	mutationFunc := operators.NormWeightedMutation()
	crossoverFunc := operators.UniformCrossover(0.5)

	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), "NSGA2"))
	logger := algos.NewProgressLogger(logPath)
	logger.InitLogging()
	logger.LogProblem(problem)

	params := nsga2.NSGA2Params{
		PopulationSize:  population,
		GenerationLimit: math.MaxInt,
		MutationFunc:    mutationFunc,
		CrossoverFunc:   crossoverFunc,
	}
	alg := nsga2.NewAlgorithm(problem, timeLimit, params, logger)

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
