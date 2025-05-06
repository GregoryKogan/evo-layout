package sga

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

type Params struct {
	PopulationSize       int
	ElitePercentile      float64
	MatingPoolPercentile float64
	MutationFunc         problems.MutationFunc
	CrossoverFunc        problems.CrossoverFunc
}
