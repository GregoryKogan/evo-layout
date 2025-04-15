package nsga2

import (
	"math"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// Individual wraps a candidate solution along with NSGA-II specific metadata.
type Individual struct {
	Solution         problems.Solution
	Rank             int
	CrowdingDistance float64
}

// NSGA2Algorithm implements the NSGA-II multiobjective evolutionary algorithm.
type NSGA2Algorithm struct {
	algos.GeneticAlgorithm // embedded common fields (start, timeout, logger, problem, etc.)
	params                 NSGA2Params
	generation             int
	population             []Individual
}

// NewAlgorithm creates a new NSGA-II instance.
func NewAlgorithm(problem problems.Problem, params NSGA2Params, logger algos.ProgressLoggerProvider) *NSGA2Algorithm {
	return &NSGA2Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, 0, logger), // timeout not used here; you can choose to add it.
		params:           params,
		generation:       0,
	}
}

// Run executes the NSGA-II process
func (alg *NSGA2Algorithm) Run() {
	alg.initPopulation()
	for time.Since(alg.StartTimestamp) < alg.Timeout {
		alg.generation++

		// Generate offspring population by selection, crossover and mutation.
		offspring := alg.makeOffspring()

		// Combine populations.
		combined := append(alg.population, offspring...)
		fronts := fastNonDominatedSort(combined)
		nextPopulation := make([]Individual, 0, alg.params.PopulationSize)
		for _, front := range fronts {
			computeCrowdingDistance(front)
			// If adding the full front would exceed population, sort by distance and pick the best individuals.
			if len(nextPopulation)+len(front) > alg.params.PopulationSize {
				sort.Slice(front, func(i, j int) bool {
					// Higher crowding distance is preferred.
					return front[i].CrowdingDistance > front[j].CrowdingDistance
				})
				remaining := alg.params.PopulationSize - len(nextPopulation)
				nextPopulation = append(nextPopulation, front[:remaining]...)
				break
			} else {
				nextPopulation = append(nextPopulation, front...)
			}
		}
		alg.population = nextPopulation

		// Log best solution from the first front.
		if len(fronts) > 0 && len(fronts[0]) > 0 {
			alg.LogStep(struct {
				Generation int               `json:"generation"`
				Best       problems.Solution `json:"best_solution"`
			}{
				Generation: alg.generation,
				Best:       fronts[0][0].Solution,
			})
		}
	}
}

// initPopulation creates the initial population randomly.
func (alg *NSGA2Algorithm) initPopulation() {
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.params.PopulationSize {
		alg.population[i] = Individual{Solution: alg.Problem.RandomSolution()}
	}
}

// makeOffspring performs selection, crossover and mutation to create offspring population.
func (alg *NSGA2Algorithm) makeOffspring() []Individual {
	offspring := make([]Individual, 0, alg.params.PopulationSize)
	// Create offspring equal to population size.
	for len(offspring) < alg.params.PopulationSize {
		// Select two parents using NSGA-II tournament selection.
		parent1 := tournamentSelection(alg.population)
		parent2 := tournamentSelection(alg.population)

		// Crossover with probability.
		var children []problems.Solution
		if rand.Float64() < alg.params.CrossoverProb {
			children = parent1.Solution.Crossover(parent2.Solution)
		} else {
			children = []problems.Solution{parent1.Solution, parent2.Solution}
		}
		// Mutate each child.
		for _, child := range children {
			child = child.Mutate(alg.params.MutationProb)
			offspring = append(offspring, Individual{Solution: child})
			if len(offspring) >= alg.params.PopulationSize {
				break
			}
		}
	}
	return offspring
}

// tournamentSelection picks one individual from the population using binary tournament selection.
func tournamentSelection(pop []Individual) Individual {
	i := rand.IntN(len(pop))
	j := rand.IntN(len(pop))
	ind1, ind2 := pop[i], pop[j]
	// Compare based on rank first (lower is better) then by crowding distance.
	if ind1.Rank < ind2.Rank {
		return ind1
	} else if ind1.Rank > ind2.Rank {
		return ind2
	}
	// If the same rank, choose the one with a larger crowding distance.
	if ind1.CrowdingDistance > ind2.CrowdingDistance {
		return ind1
	}
	return ind2
}

// fastNonDominatedSort performs the fast non-dominated sort on the population.
// Returns a slice of fronts, where each front is a slice of Individuals.
func fastNonDominatedSort(pop []Individual) [][]Individual {
	fronts := [][]Individual{}
	// For each individual, initialize the domination count and set of individuals it dominates.
	n := len(pop)
	domCount := make([]int, n)
	dominatedSet := make([][]int, n)
	rank := make([]int, n)

	for i := range n {
		dominatedSet[i] = []int{}
		for j := range n {
			if i == j {
				continue
			}
			if dominates(pop[i].Solution, pop[j].Solution) {
				dominatedSet[i] = append(dominatedSet[i], j)
			} else if dominates(pop[j].Solution, pop[i].Solution) {
				domCount[i]++
			}
		}
		if domCount[i] == 0 {
			rank[i] = 0
			pop[i].Rank = 0
		}
	}

	// Collect the first front.
	currentFront := []int{}
	for i := 0; i < n; i++ {
		if domCount[i] == 0 {
			currentFront = append(currentFront, i)
		}
	}
	var curRank int
	for len(currentFront) > 0 {
		front := []Individual{}
		nextFront := []int{}
		for _, i := range currentFront {
			pop[i].Rank = curRank
			front = append(front, pop[i])
			for _, j := range dominatedSet[i] {
				domCount[j]--
				if domCount[j] == 0 {
					nextFront = append(nextFront, j)
				}
			}
		}
		fronts = append(fronts, front)
		curRank++
		currentFront = nextFront
	}
	return fronts
}

// computeCrowdingDistance calculates the crowding distance for each individual in the front.
func computeCrowdingDistance(front []Individual) {
	l := len(front)
	if l == 0 {
		return
	}
	// Initialize distances to zero.
	for i := range front {
		front[i].CrowdingDistance = 0
	}
	// Get objectives from the first solution as representative.
	firstObjs := front[0].Solution.Objectives()
	numObjs := len(firstObjs)
	// For each objective, sort the front and update distances.
	for m := range numObjs {
		sort.Slice(front, func(i, j int) bool {
			return front[i].Solution.Objectives()[m] < front[j].Solution.Objectives()[m]
		})
		// Set boundary points to infinite distance.
		front[0].CrowdingDistance = math.Inf(1)
		front[l-1].CrowdingDistance = math.Inf(1)
		objMin := front[0].Solution.Objectives()[m]
		objMax := front[l-1].Solution.Objectives()[m]
		if objMax-objMin == 0 {
			continue
		}
		for i := 1; i < l-1; i++ {
			prev := front[i-1].Solution.Objectives()[m]
			next := front[i+1].Solution.Objectives()[m]
			front[i].CrowdingDistance += (next - prev) / (objMax - objMin)
		}
	}
}

// dominates returns true if solution a dominates solution b.
// Note that in minimization, a dominates b if all objectives of a are less than or equal to those of b, with at least one strictly less.
func dominates(a, b problems.Solution) bool {
	aObjs := a.Objectives()
	bObjs := b.Objectives()
	less := false
	for i := range aObjs {
		if aObjs[i] > bObjs[i] {
			return false
		}
		if aObjs[i] < bObjs[i] {
			less = true
		}
	}
	return less
}
