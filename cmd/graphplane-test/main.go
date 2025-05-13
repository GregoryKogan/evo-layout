package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

// 100 vertexes planar

// Uniform - 105.47
// Mirror - 137.90
// Norm - 48.33
// Percentage - 68.03
// TensionVector - 110.30
// FixedUniform - 110.63
// FixedNorm - 36.97 - Best
// FixedPercentage - 67.40
// FixedTensionVector - 85.87

// Pop:50-Gen:100 - 74.10
// Pop:100-Gen:100 - 46.40
// Pop:250-Gen:100 - 40.00
// Pop:500-Gen:100 - 30.73
// Pop:1000-Gen:100 - 36.80
// Pop:50-Gen:250 - 58.57
// Pop:100-Gen:250 - 39.10
// Pop:250-Gen:250 - 25.43
// Pop:500-Gen:250 - 11.17 (20.14) ~ 16.78
// Pop:1000-Gen:250 - 34.77
// Pop:50-Gen:500 - 40.93
// Pop:100-Gen:500 - 40.63
// Pop:250-Gen:500 - 8.73 (15.98) ~ 13.26 - Best
// Pop:500-Gen:500 - 9.60
// Pop:1000-Gen:500 - 13.60

// 100 vertexes

// Uniform - 1615.60
// Norm - 1448.58
// Mirror - 2146.62
// Percentage - 1414.20 - Best
// TensionVector - 1909.54
// FixedUniform - 1613.36
// FixedNorm - 1449.88
// FixedPercentage - 1445.82
// FixedTensionVector - 2044.44

// 50 vertexes

// Pop:50-Gen:100 - 271.50
// Pop:50-Gen:250 - 260.07
// Pop:50-Gen:500 - 252.40
// Pop:100-Gen:100 - 272.03
// Pop:100-Gen:250 - 249.97
// Pop:100-Gen:500 - 247.73
// Pop:250-Gen:100 - 256.87
// Pop:250-Gen:250 - 243.77
// Pop:250-Gen:500 - 244.23 - Best
// Pop:500-Gen:100 - 256.63
// Pop:500-Gen:250 - 243.37
// Pop:500-Gen:500 - 238.67
// Pop:1000-Gen:100 - 251.37
// Pop:1000-Gen:250 - 239.90

// Planar graph avg edges count
// 50v ~ 137 edges
// 100v ~ 285 edges
// 200v ~ 583 edges

type Test struct {
	Repeat int
	Name   string
}

func main() {
	tests := []Test{
		{Repeat: 1, Name: "FR-NSGA2"},
	}

	for _, test := range tests {
		total := 0.0
		for range test.Repeat {
			problem := graphplane.NewGraphPlaneProblem(100, 285)
			logger := initLogger(problem, test.Name)

			forceLayer := graphplane.NewForceDirectedSolver(problem.RandomSolution(), graphplane.FDSParams{Steps: 2000, Temp: 0.005, K: 1.0}, logger)
			algoSolution := forceLayer.Solve()

			gaParams := nsga2.NSGA2Params{
				PopulationSize: 250,
				MutationFunc:   mutation.Percentage(),
				CrossoverFunc:  crossover.Uniform(0.45),
			}
			alg := nsga2.NewAlgorithm(problem, gaParams, 500, logger)
			alg.Seed(algoSolution.Solution)
			alg.Run()
			gpSol, _ := alg.Solution.(*graphplane.GraphPlaneSolution)
			total += float64(gpSol.Intersections)
		}
		fmt.Printf("%s - %.2f\n", test.Name, total/float64(test.Repeat))
	}
}

func initLogger(problem problems.Problem, method string) algos.ProgressLoggerProvider {
	os.RemoveAll("logs")
	os.Mkdir("logs", 0755)

	logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), method))
	logger := algos.NewProgressLogger(logPath)
	logger.InitLogging()
	logger.LogProblem(problem)
	return logger
}
