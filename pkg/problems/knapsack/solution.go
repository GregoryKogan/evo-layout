package knapsack

import (
	"math/rand/v2"
	"slices"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type KnapsackSolution struct {
	problemParams    KnapsackProblemParams
	items            []Item
	Bits             []bool    `json:"bits"`
	CachedObjectives []float64 `json:"objectives"`
	CachedFitness    float64   `json:"fitness"`
}

func RandomKnapsackSolution(problemParams KnapsackProblemParams, items []Item) problems.Solution {
	Bits := make([]bool, problemParams.ItemsNum)
	for i := range problemParams.ItemsNum {
		Bits[i] = rand.IntN(2) == 1
	}
	return &KnapsackSolution{problemParams: problemParams, items: items, Bits: Bits}
}

func CrossoverFunc() problems.CrossoverFunc {
	return func(parentA, parentB problems.Solution) []problems.Solution {
		a, aOk := parentA.(*KnapsackSolution)
		b, bOk := parentB.(*KnapsackSolution)
		if !aOk || !bOk {
			panic("invalid parents")
		}
		return a.crossover(b)
	}
}

func MutationFunc() problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*KnapsackSolution)
		if !ok {
			panic("invalid individual")
		}
		return s.mutate()
	}
}

func (s *KnapsackSolution) crossover(other problems.Solution) []problems.Solution {
	otherKSS, ok := other.(*KnapsackSolution)
	if !ok {
		return []problems.Solution{s}
	}

	child1Bits := make([]bool, s.problemParams.ItemsNum)
	child2Bits := make([]bool, s.problemParams.ItemsNum)

	// uniform crossover
	for i := range s.problemParams.ItemsNum {
		if rand.Float64() < 0.5 {
			child1Bits[i] = s.Bits[i]
			child2Bits[i] = otherKSS.Bits[i]
		} else {
			child1Bits[i] = otherKSS.Bits[i]
			child2Bits[i] = s.Bits[i]
		}
	}

	return []problems.Solution{
		&KnapsackSolution{problemParams: s.problemParams, items: s.items, Bits: child1Bits},
		&KnapsackSolution{problemParams: s.problemParams, items: s.items, Bits: child2Bits},
	}
}

func (s *KnapsackSolution) mutate() problems.Solution {
	mutantBits := make([]bool, s.problemParams.ItemsNum)
	copy(mutantBits, s.Bits)

	for range max(s.problemParams.ItemsNum/100, 1) {
		bitInd := rand.IntN(s.problemParams.ItemsNum)
		mutantBits[bitInd] = !s.Bits[bitInd]
	}

	return &KnapsackSolution{problemParams: s.problemParams, items: s.items, Bits: mutantBits}
}

func (s *KnapsackSolution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}

	totalValue := 0
	resources := make([]int, s.problemParams.Dimensions-1)
	for i := range s.problemParams.ItemsNum {
		if s.Bits[i] {
			totalValue += s.items[i].Value
			for ri := range s.problemParams.Dimensions - 1 {
				resources[ri] += s.items[i].Resources[ri]
			}
		}
	}

	objectives := make([]float64, s.problemParams.Dimensions)
	objectives[0] = 1.0 / float64(totalValue)
	for ri := range s.problemParams.Dimensions - 1 {
		if resources[ri] > s.problemParams.Constraints[ri] {
			objectives[ri+1] = 1.0
		}
	}

	s.CachedObjectives = objectives
	return s.CachedObjectives
}

func (s *KnapsackSolution) Fitness() float64 {
	s.CachedFitness = slices.Max(s.Objectives())
	return s.CachedFitness
}
