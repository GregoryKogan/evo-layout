package nsga2

// NSGA2Params holds configurable parameters for the NSGA-II algorithm.
type NSGA2Params struct {
	PopulationSize  int
	CrossoverProb   float64
	MutationProb    float64
	GenerationLimit int
}
