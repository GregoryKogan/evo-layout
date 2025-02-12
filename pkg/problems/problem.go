package problems

type Problem interface {
	Name() string
	RandomSolution() Solution
}

type Solution interface {
	Crossover(Solution) Solution
	Mutate() Solution
	Fitness() float64
}
