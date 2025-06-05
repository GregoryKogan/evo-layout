package algos

import (
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GeneticAlgorithm struct {
	ProgressLoggerProvider
	StartTimestamp  time.Time
	GenerationLimit int
	Problem         problems.Problem
	Solution        problems.Solution
}

type GAStep struct {
	Elapsed     time.Duration     `json:"elapsed"`
	Step        int               `json:"step"`
	Solution    problems.Solution `json:"solution"`
	ParetoFront [][]float64       `json:"pareto_front"`
}

func NewGeneticAlgorithm(
	problem problems.Problem,
	generationLimit int,
	logger ProgressLoggerProvider,
) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		Problem:                problem,
		Solution:               problem.RandomSolution(),
		StartTimestamp:         time.Now(),
		GenerationLimit:        generationLimit,
		ProgressLoggerProvider: logger,
	}
}

func (ga *GeneticAlgorithm) GetSolution() problems.Solution {
	return ga.Solution
}
