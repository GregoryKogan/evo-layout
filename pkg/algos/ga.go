package algos

import (
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GeneticAlgorithm struct {
	ProgressLoggerProvider
	StartTimestamp time.Time
	Problem        problems.Problem
	Solution       problems.Solution
	TargetFitness  float64
}

type GeneticAlgorithmStep struct {
	Elapsed  time.Duration     `json:"elapsed"`
	Solution problems.Solution `json:"solution"`
}

func NewGeneticAlgorithm(problem problems.Problem, targetFitness float64, logger ProgressLoggerProvider) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		StartTimestamp:         time.Now(),
		ProgressLoggerProvider: logger,
		Problem:                problem,
		Solution:               problem.RandomSolution(),
		TargetFitness:          targetFitness,
	}
}
