package knapsack

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type KnapsackSolution struct {
	problemParams KnapsackProblemParams
	items         []Item
	Amounts       []int `json:"Amounts"`
}

func RandomKnapsackSolution(problemParams KnapsackProblemParams, items []Item) problems.Solution {
	Amounts := make([]int, problemParams.ItemsNum)
	for i := range problemParams.ItemsNum {
		Amounts[i] = rand.IntN(problemParams.InitialMaxAmount)
	}
	return &KnapsackSolution{problemParams: problemParams, items: items, Amounts: Amounts}
}

func (s *KnapsackSolution) Crossover(other problems.Solution) problems.Solution {
	otherKSS, ok := other.(*KnapsackSolution)
	if !ok {
		return s
	}

	childAmounts := make([]int, s.problemParams.ItemsNum)
	for i := range s.problemParams.ItemsNum {
		if rand.IntN(2) == 1 {
			childAmounts[i] = s.Amounts[i]
		} else {
			childAmounts[i] = otherKSS.Amounts[i]
		}
	}

	return &KnapsackSolution{problemParams: s.problemParams, items: s.items, Amounts: childAmounts}
}

func (s *KnapsackSolution) Mutate(rate float64) problems.Solution {
	mutantAmounts := make([]int, s.problemParams.ItemsNum)

	for i := range s.problemParams.ItemsNum {
		if rand.Float64() < rate {
			delta := math.Round(float64(s.Amounts[i]) * 10.0 * rand.NormFloat64() * rate)
			mutantAmounts[i] = max(s.Amounts[i]+int(delta), 0)
		}
	}

	return &KnapsackSolution{problemParams: s.problemParams, items: s.items, Amounts: mutantAmounts}
}

func (s *KnapsackSolution) Fitness() float64 {
	totalValue := 0
	resources := make([]int, s.problemParams.Dimensions-1)

	for i := range s.problemParams.ItemsNum {
		totalValue += s.Amounts[i] * s.items[i].Value
		for ri := range s.problemParams.Dimensions - 1 {
			resources[ri] += s.Amounts[i] * s.items[i].Resources[ri]
		}
	}

	fitness := float64(totalValue)
	for ri := range s.problemParams.Dimensions - 1 {
		if resources[ri] > s.problemParams.Constraints[ri] {
			fitness /= 10.0
		}
	}

	return fitness
}
