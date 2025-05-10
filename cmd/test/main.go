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
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/zdt"
)

type Task struct {
	Problem       problems.Problem
	MutationFunc  problems.MutationFunc
	CrossoverFunc problems.CrossoverFunc
}

func main() {
	// Define task instances.
	taskList := []Task{
		// graphplane.NewGraphPlaneProblem(12, 0.35, 100.0, 100.0),
		// knapsack.NewKnapsackProblem(knapsack.KnapsackProblemParams{
		// 	Dimensions:         10,
		// 	ItemsNum:           1000,
		// 	InitialMaxValue:    100,
		// 	InitialMaxResource: 100,
		// 	Constraints:        []int{30000, 30000, 30000, 30000, 30000, 30000, 30000, 30000, 30000},
		// }),
		// tsp.NewTSProblem(tsp.TSProblemParameters{CitiesNum: 100}),
		{zdt.NewZDT1Problem(30), zdt.ZDT1MutationFunc(), zdt.ZDT1CrossoverFunc()}, // 30-dimensional ZDT1
		{zdt.NewZDT2Problem(30), zdt.ZDT2MutationFunc(), zdt.ZDT2CrossoverFunc()},
		{zdt.NewZDT3Problem(30), zdt.ZDT3MutationFunc(), zdt.ZDT3CrossoverFunc()},
		{zdt.NewZDT4Problem(30), zdt.ZDT4MutationFunc(), zdt.ZDT4CrossoverFunc()},
		{zdt.NewZDT6Problem(30), zdt.ZDT6MutationFunc(), zdt.ZDT6CrossoverFunc()},
	}

	population := 100

	// Define algorithm constructors.
	algorithmsList := []struct {
		name    string
		runFunc func(task Task, logger algos.ProgressLoggerProvider)
	}{
		{
			"SGA",
			func(task Task, logger algos.ProgressLoggerProvider) {
				params := sga.Params{
					PopulationSize:  population,
					ElitePercentile: 0.1,
					MutationFunc:    task.MutationFunc,
					CrossoverFunc:   task.CrossoverFunc,
				}
				alg := sga.NewAlgorithm(task.Problem, params, 100, logger)
				alg.Run()
			},
		},
		{
			"SSGA",
			func(task Task, logger algos.ProgressLoggerProvider) {
				params := ssga.Params{
					PopulationSize: population,
					MutationFunc:   task.MutationFunc,
					CrossoverFunc:  task.CrossoverFunc,
				}
				alg := ssga.NewAlgorithm(task.Problem, params, 100, logger)
				alg.Run()
			},
		},
		{
			"NSGA2",
			func(task Task, logger algos.ProgressLoggerProvider) {
				// Configure NSGA-II parameters.
				params := nsga2.NSGA2Params{
					PopulationSize: population,
					MutationFunc:   task.MutationFunc,
					CrossoverFunc:  task.CrossoverFunc,
				}
				alg := nsga2.NewAlgorithm(task.Problem, params, 100, logger)
				alg.Run()
			},
		},
		{
			"SPEA2",
			func(task Task, logger algos.ProgressLoggerProvider) {
				params := spea2.Params{
					PopulationSize: population,
					ArchiveSize:    population,
					DensityKth:     int(math.Sqrt(float64(population + population))), // typical choice: sqrt(population+archive)
					MutationFunc:   task.MutationFunc,
					CrossoverFunc:  task.CrossoverFunc,
				}
				alg := spea2.NewAlgorithm(task.Problem, params, 100, logger)
				alg.Run()
			},
		},
	}

	// Create a directory for logs.
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	// Define timeout for every algorithm run.
	timeLimit := 1 * time.Minute

	// Loop over every problem and algorithm.
	for _, task := range taskList {
		for _, alg := range algorithmsList {
			// Create a dedicated log file for this (problem, algorithm) pair.
			logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", task.Problem.Name(), alg.name))
			progressLogger := algos.NewProgressLogger(logPath)
			progressLogger.InitLogging()
			progressLogger.LogProblem(task.Problem)

			fmt.Printf("Testing %s on %s\n", alg.name, task.Problem.Name())
			runAlgorithm(alg.name, task.Problem.Name(), func() {
				alg.runFunc(task, progressLogger)
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
