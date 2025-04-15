package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/ssga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/zdt"
)

func main() {
	// Define timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	// Define problem instances.
	problemsList := []problems.Problem{
		// graphplane.NewGraphPlaneProblem(12, 0.35, 100.0, 100.0),
		// knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
		// 	Dimensions:         10,
		// 	ItemsNum:           1000,
		// 	InitialMaxValue:    100,
		// 	InitialMaxResource: 100,
		// 	Constraints:        []int{30000, 30000, 30000, 30000, 30000, 30000, 30000, 30000, 30000},
		// }),
		// tsp.NewTSProblem(tsp.TSProblemParameters{CitiesNum: 100}),
		zdt.NewZDT1Problem(30), // 30-dimensional ZDT1
		zdt.NewZDT2Problem(30),
		zdt.NewZDT3Problem(30),
		zdt.NewZDT4Problem(30),
		zdt.NewZDT6Problem(30),
	}

	// Define algorithm constructors.
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
					MutationRate:         0.02,
				}
				alg := sga.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
		{
			"SSGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := ssga.Params{
					PopulationSize: 1000,
					MutationRate:   0.02,
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
					PopulationSize: 100,
					CrossoverProb:  0.9,
					MutationProb:   0.1,
				}
				alg := nsga2.NewAlgorithm(problem, timeLimit, params, logger)
				alg.Run()
			},
		},
	}

	// Create a directory for logs.
	os.RemoveAll("logs")
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
