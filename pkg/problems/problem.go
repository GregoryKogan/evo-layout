package problems

type Problem interface {
	Name() string
	RandomSolution() Solution
}

type AlgorithmicProblem interface {
	Problem
	AlgorithmicSolution() Solution
}

type Solution interface {
	Crossover(Solution) []Solution
	Mutate(rate float64) Solution
	Fitness() float64
}
