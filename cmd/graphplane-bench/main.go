package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	problem := graphplane.NewGraphPlaneProblem(100, 0.1, 1.0, 1.0)
	population := 500

	// Define algorithm constructors.
	algorithmsList := []struct {
		name    string
		runFunc func(problem problems.Problem, logger algos.ProgressLoggerProvider)
	}{
		{
			"SGA-35",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       population,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationFunc:         operators.NormWeightedMutation(),
					CrossoverFunc:        operators.UniformCrossover(0.35),
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SGA-40",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       population,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationFunc:         operators.NormWeightedMutation(),
					CrossoverFunc:        operators.UniformCrossover(0.40),
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SGA-45",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       population,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationFunc:         operators.NormWeightedMutation(),
					CrossoverFunc:        operators.UniformCrossover(0.45),
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SGA-50",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       population,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationFunc:         operators.NormWeightedMutation(),
					CrossoverFunc:        operators.UniformCrossover(0.5),
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
	}
	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	// Loop over ever algorithm.
	for _, alg := range algorithmsList {
		// Create a dedicated log file for this (problem, algorithm) pair.
		logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), alg.name))
		progressLogger := algos.NewProgressLogger(logPath)
		progressLogger.InitLogging()
		progressLogger.LogProblem(problem)

		fmt.Printf("Testing %s on %s\n", alg.name, problem.Name())
		runAlgorithm(alg.name, problem.Name(), func() {
			alg.runFunc(problem, progressLogger)
		}, timeLimit)
	}
}

// runAlgorithm wraps an algorithm run in a goroutine and enforces a timeout.
func runAlgorithm(algName string, problemName string, algRun func(), timeout time.Duration) {
	fmt.Printf("Running %s on %s with timeout %v\n", algName, problemName, timeout)
	done := make(chan bool)
	go func() {
		algRun()
		done <- true
	}()
	select {
	case <-done:
		fmt.Printf("%s on %s completed.\n", algName, problemName)
	case <-time.After(timeout):
		fmt.Printf("%s on %s timed out after %v.\n", algName, problemName, timeout)
	}
}
