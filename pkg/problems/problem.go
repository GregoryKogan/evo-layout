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

type Solution interface {
	Crossover(Solution) []Solution
	Mutate(rate float64) Solution
	Fitness() float64
}
