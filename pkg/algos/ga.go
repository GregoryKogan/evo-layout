package algos

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GeneticAlgorithm struct {
	ProgressLoggerProvider
	Problem       problems.Problem
	Solution      problems.Solution
	TargetFitness float64
}

func NewGeneticAlgorithm(problem problems.Problem, targetFitness float64, logFilepath string) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		ProgressLoggerProvider: NewProgressLogger(logFilepath),
		Problem:                problem,
		Solution:               problem.RandomSolution(),
		TargetFitness:          targetFitness,
	}
}
