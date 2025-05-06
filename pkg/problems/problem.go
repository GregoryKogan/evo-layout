package problems

import "time"

type Problem interface {
	Name() string
	RandomSolution() Solution
}

type AlgorithmicSolution struct {
	Solution `json:"solution"`
	TimeTook time.Duration `json:"took"`
}

type AlgorithmicProblem interface {
	Problem
	AlgorithmicSolution() AlgorithmicSolution
}

type CrossoverFunc func(parentA, parentB Solution) []Solution
type MutationFunc func(individual Solution) Solution

type Solution interface {
	// Multi-objective genetic algorithms (like NSGA-II, SPEA2) use Objectives() method
	// Single objective problems just return a slice with one value
	// Single-objective genetic algorithms (like SGA, SSGA) use Fitness() method
	// Multi objective problems compute single fitness value based on objectives
	// (average, mean, min, max, ... or any other formula can be used)
	Objectives() []float64
	Fitness() float64
}
