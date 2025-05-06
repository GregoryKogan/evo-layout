package algos

import (
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GeneticAlgorithm struct {
	ProgressLoggerProvider
	StartTimestamp time.Time
	Timeout        time.Duration
	Problem        problems.Problem
	Solution       problems.Solution
}

type GeneticAlgorithmStep struct {
	Elapsed     time.Duration     `json:"elapsed"`
	Solution    problems.Solution `json:"solution"`
	ParetoFront [][]float64       `json:"pareto_front"`
}

func NewGeneticAlgorithm(
	problem problems.Problem,
	timeout time.Duration,
	logger ProgressLoggerProvider,
) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		Problem:                problem,
		Solution:               problem.RandomSolution(),
		StartTimestamp:         time.Now(),
		Timeout:                timeout,
		ProgressLoggerProvider: logger,
	}
}
