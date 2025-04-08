package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/ssga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/knapsack"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/tsp"
)

func main() {
	// Define 5 minute timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	// Define problem instances.
	problemsList := []problems.Problem{
		graphplane.NewGraphPlaneProblem(12, 0.35, 100.0, 100.0),
		knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
			Dimensions:         2,
			ItemsNum:           1000,
			InitialMaxValue:    100,
			InitialMaxResource: 100,
			Constraints:        []int{300000},
		}),
		tsp.NewTSProblem(tsp.TSProblemParameters{CitiesNum: 50}),
	}

	// Define algorithm constructors.
	// Each algorithm is configured with example parameters and a target fitness of 1.0.
	targetFitness := 1.0
	algorithmsList := []struct {
		name    string
		runFunc func(problem problems.Problem, logger algos.ProgressLoggerProvider)
	}{
		{
			"SGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       1000,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationRate:         0.01,
				}
				alg := sga.NewAlgorithm(problem, targetFitness, params, logger)
				alg.Run()
			},
		},
		{
			"SSGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := ssga.Params{
					PopulationSize: 1000,
					MutationRate:   0.01,
				}
				alg := ssga.NewAlgorithm(problem, targetFitness, params, logger)
				alg.Run()
			},
		},
	}

	// Create a directory for logs.
	os.Mkdir("logs", 0755)

	// Loop over every problem and algorithm.
	for _, prob := range problemsList {
		for _, alg := range algorithmsList {
			// Create a dedicated log file for this (problem, algorithm) pair.
			logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", prob.Name(), alg.name))
			progressLogger := algos.NewProgressLogger(logPath)
			progressLogger.InitLogging()
			progressLogger.LogProblem(prob)

			fmt.Printf("Testing %s on %s\n", alg.name, prob.Name())
			runAlgorithm(alg.name, prob.Name(), func() {
				alg.runFunc(prob, progressLogger)
			}, timeLimit)
		}
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
