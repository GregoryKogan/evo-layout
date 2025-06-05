package nsga2

import (
	"context"
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
type Algorithm struct {
	algos.GeneticAlgorithm // embedded common fields (start time, timeout, logger, problem, etc.)
	params                 Params
	generation             int
	population             []Individual
}

// NewAlgorithm creates a new NSGA-II instance.
func NewAlgorithm(problem problems.Problem, params Params, generationLimit int, logger algos.ProgressLoggerProvider) *Algorithm {
	return &Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, generationLimit, logger),
		params:           params,
		generation:       0,
	}
}

func (alg *Algorithm) Seed(seedSolution problems.Solution) {
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.params.PopulationSize {
		alg.population[i] = Individual{Solution: alg.params.MutationFunc(seedSolution)}
	}
	alg.population[0] = Individual{Solution: seedSolution}
	alg.Solution = seedSolution
}

func (alg *Algorithm) SetPopulation(pop []problems.Solution) {
	if len(pop) != alg.params.PopulationSize {
		panic("Wrong population size")
	}
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.params.PopulationSize {
		alg.population[i] = Individual{Solution: pop[i]}
	}
}

func (alg *Algorithm) GetSteps() int {
	return alg.generation
}

// Run executes the NSGA-II process until timeout.
func (alg *Algorithm) Run(ctx context.Context) {
	if len(alg.population) < alg.params.PopulationSize {
		alg.initPopulation()
	}
	for alg.generation < alg.GenerationLimit {
		select {
		case <-ctx.Done():
			return
		default:
		}

		alg.generation++

		// Generate offspring population by selection, crossover and mutation.
		offspring := alg.makeOffspring()

		// Combine populations.
		combined := append(alg.population, offspring...)
		fronts := fastNonDominatedSort(combined)
		nextPopulation := make([]Individual, 0, alg.params.PopulationSize)
		for _, front := range fronts {
			computeCrowdingDistance(front)
			// If adding the full front would exceed population, sort by crowding distance.
			if len(nextPopulation)+len(front) > alg.params.PopulationSize {
				sort.Slice(front, func(i, j int) bool {
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

		// Log current generation data: record generation number and the Pareto front
		if len(fronts) > 0 {
			var pareto [][]float64
			for _, ind := range fronts[0] {
				pareto = append(pareto, ind.Solution.Objectives())
				if ind.Solution.Fitness() < alg.Solution.Fitness() {
					alg.Solution = ind.Solution
				}
			}

			if !alg.params.Verbose {
				pareto = nil
			}

			if alg.ProgressLoggerProvider != nil {
				alg.LogStep(algos.GAStep{
					Elapsed:     time.Since(alg.StartTimestamp),
					Step:        alg.generation,
					ParetoFront: pareto,
					Solution:    alg.Solution,
				})
			}
		}
	}
}

// initPopulation creates the initial population randomly.
func (alg *Algorithm) initPopulation() {
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.params.PopulationSize {
		alg.population[i] = Individual{Solution: alg.Problem.RandomSolution()}
	}
}

// makeOffspring performs selection, crossover and mutation to create offspring population.
func (alg *Algorithm) makeOffspring() []Individual {
	offspring := make([]Individual, 0, alg.params.PopulationSize)
	for len(offspring) < alg.params.PopulationSize {
		parent1 := tournamentSelection(alg.population)
		parent2 := tournamentSelection(alg.population)

		children := alg.params.CrossoverFunc(parent1.Solution, parent2.Solution)

		for _, child := range children {
			child = alg.params.MutationFunc(child)
			offspring = append(offspring, Individual{Solution: child})
			if len(offspring) >= alg.params.PopulationSize {
				break
			}
		}
	}
	return offspring
}

// tournamentSelection picks one individual using binary tournament selection.
func tournamentSelection(pop []Individual) Individual {
	i := rand.IntN(len(pop))
	j := rand.IntN(len(pop))
	ind1, ind2 := pop[i], pop[j]
	// Compare by rank (lower is better) then crowding distance.
	if ind1.Rank < ind2.Rank {
		return ind1
	} else if ind1.Rank > ind2.Rank {
		return ind2
	}
	if ind1.CrowdingDistance > ind2.CrowdingDistance {
		return ind1
	}
	return ind2
}

// fastNonDominatedSort performs fast non-dominated sort.
func fastNonDominatedSort(pop []Individual) [][]Individual {
	fronts := [][]Individual{}
	n := len(pop)
	domCount := make([]int, n)
	dominatedSet := make([][]int, n)
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
	}
	currentFront := []int{}
	for i := range n {
		if domCount[i] == 0 {
			pop[i].Rank = 0
			currentFront = append(currentFront, i)
		}
	}
	curRank := 0
	for len(currentFront) > 0 {
		var front []Individual
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

// computeCrowdingDistance calculates crowding distances for a front.
func computeCrowdingDistance(front []Individual) {
	l := len(front)
	if l == 0 {
		return
	}
	for i := range front {
		front[i].CrowdingDistance = 0
	}
	numObjs := len(front[0].Solution.Objectives())
	for m := range numObjs {
		sort.Slice(front, func(i, j int) bool {
			return front[i].Solution.Objectives()[m] < front[j].Solution.Objectives()[m]
		})
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

// dominates returns true if solution a dominates solution b (minimization).
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
