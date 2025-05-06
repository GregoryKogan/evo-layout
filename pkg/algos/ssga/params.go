package ssga

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

type Params struct {
	PopulationSize int
	MutationFunc   problems.MutationFunc
	CrossoverFunc  problems.CrossoverFunc
}
