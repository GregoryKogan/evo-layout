package problems

type Problem interface {
	Name() string
	RandomSolution() Solution
	MarshalJSON() ([]byte, error)
}

type Solution interface {
	Crossover(Solution) Solution
	Mutate() Solution
	Fitness() float64
	MarshalJSON() ([]byte, error)
}
