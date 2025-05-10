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
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

// 100 vertexes, 0.1 fill, 1min limit, population 500, algo SSGA, uniform(0.45) crossover
// Tension Vector - 11396 intersections (eps=0.01)
// Norm           - 12954 intersections
// Percentage     - 13735 intersections
// Uniform        - 13921 intersections
// Mirror         - 14011 intersections

// 100 vertexes, 0.1 fill, 1min limit, population 500, algo SSGA, TV(eps=0.01) mutation
// uniform(0.30) crossover - 14015 intersections
// uniform(0.35) crossover - 13531 intersections
// uniform(0.40) crossover - 12349 intersections
// uniform(0.45) crossover - 12340 intersections (best)
// uniform(0.50) crossover - 12912 intersections

// Planar 100 vertexes, 5min limit, population 500, TV(eps=0.01) mutation, uniform(0.45) crossover
// SPEA2 - 1505 intersections
// NSGA2 - 1794 intersections
// SSGA  - 1878 intersections
// SGA   - 3697 intersections
// Planar 100 vertexes Force (2000 iterations, 0.002 init temp) - 108 intersections
// NSGA2 100 seconds -> Force 2000 iters - 109 intersections
// Force 2000 iters -> NSGA2 1 min - 106 intersections (133 -> 106)

// Planar 50 vertexes, 5min limit, population 500, TV(eps=0.01) mutation, uniform(0.45) crossover
// NSGA2 - 169 intersections
// SPEA2 - 295 intersections
// SSGA  - 312 intersections
// SGA   - 622 intersections
// Planar 50 vertexes Force (2000 iterations, 0.002 init temp) - 65 intersections
// NSGA2 100 seconds -> Force 2000 iters - 25 intersections
// Force 2000 iters -> NSGA2 1 min - 8 intersections (54 -> 8)

func main() {
	problem := graphplane.NewPlanarGraphPlaneProblem(50)
	population := 500

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
					MutationFunc:         mutation.TensionVector(0.01),
					CrossoverFunc:        crossover.Uniform(0.45),
				}
				alg := sga.NewAlgorithm(problem, params, math.MaxInt, logger)
				alg.Run()
			},
		},
		{
			"SSGA",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := ssga.Params{
					PopulationSize: population,
					MutationFunc:   mutation.TensionVector(0.01),
					CrossoverFunc:  crossover.Uniform(0.45),
				}
				alg := ssga.NewAlgorithm(problem, params, math.MaxInt, logger)
				alg.Run()
			},
		},
		{
			"NSGA2",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				// Configure NSGA-II parameters.
				params := nsga2.NSGA2Params{
					PopulationSize: population,
					MutationFunc:   mutation.TensionVector(0.01),
					CrossoverFunc:  crossover.Uniform(0.45),
				}
				alg := nsga2.NewAlgorithm(problem, params, math.MaxInt, logger)
				alg.Run()
			},
		},
		{
			"SPEA2",
			func(problem problems.Problem, logger algos.ProgressLoggerProvider) {
				params := spea2.Params{
					PopulationSize: population,
					ArchiveSize:    population,
					DensityKth:     int(math.Sqrt(float64(population + population))), // typical choice: sqrt(population+archive)
					MutationFunc:   mutation.TensionVector(0.01),
					CrossoverFunc:  crossover.Uniform(0.45),
				}
				alg := spea2.NewAlgorithm(problem, params, math.MaxInt, logger)
				alg.Run()
			},
		},
	}
	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	// Define timeout for every algorithm run.
	timeLimit := 5 * time.Minute

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
