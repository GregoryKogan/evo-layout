package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/spea2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/ssga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	problem := graphplane.NewGraphPlaneProblem(100, 0.1, 1.0, 1.0)
	population := 1000
	mutationProb := 0.1
	crossoverProb := 0.9

	// Define algorithm constructors.
	algorithmsList := []struct {
		name    string
		runFunc func(problem problems.Problem, logger algos.ProgressLoggerProvider)
	}{
		{
			"SGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:       population,
					ElitePercentile:      0.1,
					MatingPoolPercentile: 0.5,
					MutationProb:         mutationProb,
					CrossoverProb:        crossoverProb,
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SSGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := ssga.Params{
					PopulationSize: population,
					MutationProb:   mutationProb,
					CrossoverProb:  crossoverProb,
				}
				alg := ssga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"NSGA2",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				// Configure NSGA-II parameters.
				params := nsga2.NSGA2Params{
					PopulationSize:  population,
					GenerationLimit: math.MaxInt,
					MutationProb:    mutationProb,
					CrossoverProb:   crossoverProb,
				}
				alg := nsga2.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SPEA2",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := spea2.Params{
					PopulationSize:  population,
					ArchiveSize:     population,
					DensityKth:      int(math.Sqrt(float64(population + population))), // typical choice: sqrt(population+archive)
					GenerationLimit: math.MaxInt,
					MutationProb:    mutationProb,
					CrossoverProb:   crossoverProb,
				}
				alg := spea2.NewAlgorithm(problem, timeLimit, params, logger)
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
