package spea2

import (
	"math"
	"math/rand/v2"
	"sort"
	"time"

	"slices"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// Individual wraps a solution plus SPEA2 metadata.
type Individual struct {
	sol      problems.Solution
	strength float64
	rawFit   float64
	density  float64
	fitness  float64
}

// Algorithm implements the SPEA2 EA with logging.
type Algorithm struct {
	algos.GeneticAlgorithm // embeds StartTimestamp, Timeout, Problem, Logger
	params                 Params
	population             []Individual
	archive                []Individual
	generation             int
}

// Step is emitted each generation and contains elapsed time, generation index,
// and the current Pareto front (archive members with rawFit < 1).
type Step struct {
	algos.GeneticAlgorithmStep
	Generation  int         `json:"generation"`
	ParetoFront [][]float64 `json:"pareto_front"`
}

// NewAlgorithm constructs a Algorithm.
func NewAlgorithm(problem problems.Problem, timeout time.Duration, params Params, logger algos.ProgressLoggerProvider) *Algorithm {
	ga := algos.NewGeneticAlgorithm(problem, timeout, logger)
	return &Algorithm{
		GeneticAlgorithm: *ga,
		params:           params,
		generation:       0,
	}
}

// Run executes SPEA2 until Timeout, logging one Step per generation.
func (alg *Algorithm) Run() {
	alg.initPopulation()
	alg.archive = nil

	for time.Since(alg.StartTimestamp) < alg.Timeout && alg.generation < alg.params.GenerationLimit {
		alg.generation++

		// 1) Combine pop + archive for fitness assignment
		combined := slices.Concat(alg.population, alg.archive)

		// 2) Fitness assignment
		alg.assignFitness(combined)

		// 3) Environmental selection → new archive
		alg.updateArchive(combined)

		// 4) Log current archive Pareto front
		var pareto [][]float64
		for _, ind := range alg.archive {
			if ind.rawFit < 1 {
				pareto = append(pareto, ind.sol.Objectives())
			}
		}
		alg.LogStep(Step{
			GeneticAlgorithmStep: algos.GeneticAlgorithmStep{
				Elapsed: time.Since(alg.StartTimestamp),
			},
			Generation:  alg.generation,
			ParetoFront: pareto,
		})

		// 5) Reproduction → new population
		alg.reproduce()
	}
}

// initPopulation creates the initial random population.
func (alg *Algorithm) initPopulation() {
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.population {
		alg.population[i] = Individual{sol: alg.Problem.RandomSolution()}
	}
}

// assignFitness computes strength, raw fitness, density, and final fitness.
func (alg *Algorithm) assignFitness(all []Individual) {
	n := len(all)
	// Strength: count dominated
	for i := range all {
		all[i].strength = 0
		for j := range all {
			if dominates(all[i].sol, all[j].sol) {
				all[i].strength++
			}
		}
	}
	// Raw fitness: sum strengths of those dominating i
	for i := range all {
		all[i].rawFit = 0
		for j := range all {
			if dominates(all[j].sol, all[i].sol) {
				all[i].rawFit += all[j].strength
			}
		}
	}
	// Density: distance to k-th neighbor in objective space
	for i := range all {
		fi := all[i].sol.Objectives()
		dists := make([]float64, 0, n-1)
		for j := range all {
			if i == j {
				continue
			}
			fj := all[j].sol.Objectives()
			dists = append(dists, euclidean(fi, fj))
		}
		sort.Float64s(dists)
		k := alg.params.DensityKth
		if k >= len(dists) {
			k = len(dists) - 1
		}
		all[i].density = 1.0 / (dists[k] + 2.0)
	}
	// Final fitness = raw + density
	for i := range all {
		all[i].fitness = all[i].rawFit + all[i].density
	}
}

// updateArchive builds the next archive based on fitness.
func (alg *Algorithm) updateArchive(combined []Individual) {
	// 1) keep all with rawFit < 1
	var nextA []Individual
	for _, ind := range combined {
		if ind.rawFit < 1 {
			nextA = append(nextA, ind)
		}
	}
	// 2) truncate if too large
	if len(nextA) > alg.params.ArchiveSize {
		sort.Slice(nextA, func(i, j int) bool {
			return nextA[i].density > nextA[j].density
		})
		nextA = nextA[:alg.params.ArchiveSize]
	}
	// 3) fill if too small
	if len(nextA) < alg.params.ArchiveSize {
		var dominated []Individual
		for _, ind := range combined {
			if ind.rawFit >= 1 {
				dominated = append(dominated, ind)
			}
		}
		sort.Slice(dominated, func(i, j int) bool {
			return dominated[i].fitness < dominated[j].fitness
		})
		need := alg.params.ArchiveSize - len(nextA)
		nextA = append(nextA, dominated[:need]...)
	}
	alg.archive = nextA
}

// reproduce fills the next population via tournament, crossover, and mutation.
func (alg *Algorithm) reproduce() {
	var nextP []Individual
	for len(nextP) < alg.params.PopulationSize {
		p1 := alg.tournamentSelect()
		p2 := alg.tournamentSelect()
		var children []problems.Solution
		if rand.Float64() < alg.params.CrossoverProb {
			children = p1.sol.Crossover(p2.sol)
		} else {
			children = []problems.Solution{p1.sol, p2.sol}
		}
		for _, c := range children {
			c = c.Mutate(alg.params.MutationProb)
			nextP = append(nextP, Individual{sol: c})
			if len(nextP) >= alg.params.PopulationSize {
				break
			}
		}
	}
	alg.population = nextP
}

// tournamentSelect chooses one archive member by binary tournament on fitness.
func (alg *Algorithm) tournamentSelect() Individual {
	i := rand.IntN(len(alg.archive))
	j := rand.IntN(len(alg.archive))
	if alg.archive[i].fitness < alg.archive[j].fitness {
		return alg.archive[i]
	}
	return alg.archive[j]
}

// dominates checks Pareto dominance (minimize objectives).
func dominates(a, b problems.Solution) bool {
	ao, bo := a.Objectives(), b.Objectives()
	strictly := false
	for i := range ao {
		if ao[i] > bo[i] {
			return false
		}
		if ao[i] < bo[i] {
			strictly = true
		}
	}
	return strictly
}

// euclidean computes Euclidean distance between objective vectors.
func euclidean(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}
