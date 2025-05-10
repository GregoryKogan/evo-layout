package spea2

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

type Params struct {
	PopulationSize int // μ: population size
	ArchiveSize    int // size of the external archive
	DensityKth     int // k for k‑th nearest neighbor density estimation
	MutationFunc   problems.MutationFunc
	CrossoverFunc  problems.CrossoverFunc
}
