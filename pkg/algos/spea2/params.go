package spea2

type Params struct {
	PopulationSize  int     // μ: population size
	ArchiveSize     int     // size of the external archive
	CrossoverProb   float64 // probability of performing crossover
	MutationProb    float64 // probability of mutation per individual
	DensityKth      int     // k for k‑th nearest neighbor density estimation
	GenerationLimit int
}
