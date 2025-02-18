package problems

type Problem interface {
	Name() string
	RandomSolution() Solution
}

type Solution interface {
	Crossover(Solution) Solution
	Mutate(rate float64) Solution
	Fitness() float64
}
