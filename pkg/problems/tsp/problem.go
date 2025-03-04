package tsp

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type TSProblem struct {
	Params TSProblemParameters `json:"parameters"`
	Cities []City              `json:"cities"`
}

type TSProblemParameters struct {
	CitiesNum int `json:"cities_num"`
}

func (p *TSProblemParameters) validate() {
	if p.CitiesNum < 2 {
		panic("TSP must have at least 2 cities")
	}
}

func NewTSProblem(params TSProblemParameters) problems.AlgorithmicProblem {
	params.validate()

	cities := make([]City, 0, params.CitiesNum)
	for range params.CitiesNum {
		cities = append(cities, NewRandomCity())
	}
	return &TSProblem{Params: params, Cities: cities}
}

func (p *TSProblem) Name() string {
	return "TSP"
}

func (p *TSProblem) RandomSolution() problems.Solution {
	return RandomTSPSolution(p.Params, p.Cities)
}

func (p *TSProblem) AlgorithmicSolution() problems.Solution {
	cities := make([]int, p.Params.CitiesNum-1)
	for i := range p.Params.CitiesNum - 1 {
		cities[i] = i + 1
	}

	bestSolution := TSPSolution{problemParams: p.Params, cities: p.Cities, VisitingOrder: cities}
	for _, order := range permutations(cities) {
		solution := TSPSolution{problemParams: p.Params, cities: p.Cities, VisitingOrder: order}
		if solution.Fitness() > bestSolution.Fitness() {
			bestSolution = solution
		}
	}

	return &bestSolution
}

func permutations(arr []int) [][]int {
	// https://stackoverflow.com/a/30226442
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
