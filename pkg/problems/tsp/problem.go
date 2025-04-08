package tsp

import (
	"log"
	"math"
	"time"

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

// Branch and bound
func (p *TSProblem) AlgorithmicSolution() problems.AlgorithmicSolution {
	startTime := time.Now()
	n := p.Params.CitiesNum

	log.Printf("Starting TSP Branch and Bound for %d cities", n)

	// bestCost: the total distance of the best tour found so far.
	bestCost := math.Inf(1)
	// bestPath holds the complete tour (including the starting city at index 0)
	var bestPath []int

	// Precompute the minimum outgoing distance (minEdge) for each city.
	minEdge := make([]float64, n)
	for i := 0; i < n; i++ {
		min := math.Inf(1)
		for j := 0; j < n; j++ {
			if i != j {
				d := p.Cities[i].Distance(p.Cities[j])
				if d < min {
					min = d
				}
			}
		}
		minEdge[i] = min
	}

	// Count number of DFS nodes visited for logging.
	var nodesVisited int64

	// DFS helper function implementing branch and bound.
	// current: index of the current city.
	// count: number of cities visited so far.
	// currCost: cumulative distance so far.
	// path: sequence of visited cities (including the starting city at index 0).
	// visited: boolean slice tracking visited cities.
	var dfs func(current, count int, currCost float64, path []int, visited []bool)
	dfs = func(current, count int, currCost float64, path []int, visited []bool) {
		nodesVisited++
		// Log progress every 10_000_000 nodes visited.
		if nodesVisited%10_000_000 == 0 {
			log.Printf("Visited nodes: %d, current best cost: %.5f", nodesVisited, bestCost)
		}

		// When all cities are visited, complete the cycle by returning to the starting city.
		if count == n {
			totalCost := currCost + p.Cities[current].Distance(p.Cities[0])
			if totalCost < bestCost {
				bestCost = totalCost
				// Copy current complete path as the best path.
				bestPath = make([]int, len(path))
				copy(bestPath, path)
				log.Printf("New best path found with cost: %.5f after %d nodes", bestCost, nodesVisited)
			}
			return
		}

		// Compute a lower bound for this partial solution:
		// lb = current cost + (minimum cost edge from current city) +
		//      sum(minimum outgoing edge for each unvisited city)
		lb := currCost + minEdge[current]
		for i := range n {
			if !visited[i] {
				lb += minEdge[i]
			}
		}

		// Prune this branch if the lower bound is no better than the best complete tour.
		if lb >= bestCost {
			return
		}

		// Continue to add cities to the current path.
		for i := 1; i < n; i++ {
			if !visited[i] {
				visited[i] = true
				dfs(i, count+1, currCost+p.Cities[current].Distance(p.Cities[i]), append(path, i), visited)
				visited[i] = false
			}
		}
	}

	// Initialize visited slice with city 0 as the start.
	visited := make([]bool, n)
	visited[0] = true
	dfs(0, 1, 0.0, []int{0}, visited)

	// Log final statistics.
	elapsed := time.Since(startTime)
	log.Printf("DFS completed. Total nodes visited: %d, Best cost: %.5f, Time elapsed: %v", nodesVisited, bestCost, elapsed)

	// bestPath holds the full tour (starting at city 0), but TSPSolution's VisitingOrder excludes the starting city.
	visitingOrder := bestPath[1:]
	bestSolution := TSPSolution{
		problemParams: p.Params,
		cities:        p.Cities,
		VisitingOrder: visitingOrder,
		// CachedFitness is defined as 1/totalDistance.
		CachedFitness: 1.0 / bestCost,
	}

	return problems.AlgorithmicSolution{
		Solution: &bestSolution,
		TimeTook: elapsed,
	}
}

// Brute force
func (p *TSProblem) BruteForceSolution() problems.AlgorithmicSolution {
	start := time.Now()

	cities := make([]int, p.Params.CitiesNum-1)
	for i := range p.Params.CitiesNum - 1 {
		cities[i] = i + 1
	}

	permutations := func(arr []int) [][]int {
		// https://stackoverflow.com/a/30226442
		var helper func([]int, int)
		res := [][]int{}

		helper = func(arr []int, n int) {
			if n == 1 {
				tmp := make([]int, len(arr))
				copy(tmp, arr)
				res = append(res, tmp)
			} else {
				for i := range n {
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

	bestSolution := TSPSolution{problemParams: p.Params, cities: p.Cities, VisitingOrder: cities}
	for _, order := range permutations(cities) {
		solution := TSPSolution{problemParams: p.Params, cities: p.Cities, VisitingOrder: order}
		if solution.Fitness() > bestSolution.Fitness() {
			bestSolution = solution
		}
	}

	return problems.AlgorithmicSolution{
		Solution: &bestSolution,
		TimeTook: time.Since(start),
	}
}
