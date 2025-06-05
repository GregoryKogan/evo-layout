package nsga2

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

// Params holds configurable parameters for the NSGA-II algorithm.
type Params struct {
	PopulationSize int
	MutationFunc   problems.MutationFunc
	CrossoverFunc  problems.CrossoverFunc
	Verbose        bool
}
