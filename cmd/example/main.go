package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/sga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

func main() {
	vertexes := 200
	timeLimit := 300 * time.Second

	problem := graphplane.NewPlanarGraphPlaneProblem(vertexes)

	logger := initLogger(problem, "SGA")

	params := sga.Params{
		PopulationSize:       500,
		ElitePercentile:      0.1,
		MatingPoolPercentile: 0.5,
		MutationFunc:         mutation.Uniform(),
		CrossoverFunc:        crossover.Uniform(0.5),
	}
	alg := sga.NewAlgorithm(problem, params, math.MaxInt, logger)

	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	alg.Run(ctx)
	cancel()
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
