package nsga2

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

// NSGA2Params holds configurable parameters for the NSGA-II algorithm.
type NSGA2Params struct {
	PopulationSize int
	MutationFunc   problems.MutationFunc
	CrossoverFunc  problems.CrossoverFunc
	Verbose        bool
}
