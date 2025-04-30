package tsp

import (
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type TSPSolution struct {
	problemParams TSProblemParameters
	cities        []City
	VisitingOrder []int   `json:"order"`
	CachedFitness float64 `json:"fitness"`
}

func RandomTSPSolution(problemParams TSProblemParameters, cities []City) problems.Solution {
	order := rand.Perm(problemParams.CitiesNum - 1)
	for i := range problemParams.CitiesNum - 1 {
		order[i]++
	}
	return &TSPSolution{problemParams: problemParams, cities: cities, VisitingOrder: order}
}

func (s *TSPSolution) Mutate() problems.Solution {
	newOrder := make([]int, s.problemParams.CitiesNum-1)
	copy(newOrder, s.VisitingOrder)

	i := rand.IntN(s.problemParams.CitiesNum - 1)
	j := rand.IntN(s.problemParams.CitiesNum - 1)
	newOrder[i], newOrder[j] = newOrder[j], newOrder[i]

	return &TSPSolution{problemParams: s.problemParams, cities: s.cities, VisitingOrder: newOrder}
}

func (s *TSPSolution) Crossover(other problems.Solution) []problems.Solution {
	otherTSS, ok := other.(*TSPSolution)
	if !ok || s.problemParams.CitiesNum != otherTSS.problemParams.CitiesNum {
		return []problems.Solution{s}
	}

	order1 := make([]int, s.problemParams.CitiesNum-1)
	order2 := make([]int, s.problemParams.CitiesNum-1)
	set1 := make(map[int]bool, s.problemParams.CitiesNum-1)
	set2 := make(map[int]bool, s.problemParams.CitiesNum-1)
	for i := range s.problemParams.CitiesNum - 1 {
		if rand.IntN(2) == 1 {
			order1[i] = s.VisitingOrder[i]
			set1[s.VisitingOrder[i]] = true
		} else {
			order2[i] = otherTSS.VisitingOrder[i]
			set2[otherTSS.VisitingOrder[i]] = true
		}
	}

	sIndex := 0
	oIndex := 0
	for i := range s.problemParams.CitiesNum - 1 {
		if order1[i] == 0 {
			for _, ok := set1[otherTSS.VisitingOrder[oIndex]]; ok; {
				oIndex++
				_, ok = set1[otherTSS.VisitingOrder[oIndex]]
			}
			order1[i] = otherTSS.VisitingOrder[oIndex]
			set1[otherTSS.VisitingOrder[oIndex]] = true
			oIndex++
		} else {
			for _, ok := set2[s.VisitingOrder[sIndex]]; ok; {
				sIndex++
				_, ok = set2[s.VisitingOrder[sIndex]]
			}
			order2[i] = s.VisitingOrder[sIndex]
			set2[s.VisitingOrder[sIndex]] = true
			sIndex++
		}
	}

	return []problems.Solution{
		&TSPSolution{problemParams: s.problemParams, cities: s.cities, VisitingOrder: order1},
		&TSPSolution{problemParams: s.problemParams, cities: s.cities, VisitingOrder: order2},
	}
}

func (s *TSPSolution) Fitness() float64 {
	if s.CachedFitness != 0 {
		return s.CachedFitness
	}

	totalLength := 0.0
	curCity := s.cities[0]

	for _, i := range s.VisitingOrder {
		nextCity := s.cities[i]
		totalLength += curCity.Distance(nextCity)
		curCity = nextCity
	}
	totalLength += curCity.Distance(s.cities[0])

	s.CachedFitness = totalLength
	return s.CachedFitness
}

func (s *TSPSolution) Objectives() []float64 {
	return []float64{s.Fitness()}
}
