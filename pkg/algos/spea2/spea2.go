package spea2

import (
	"math"
	"slices"
	"sort"
	"time"

	"math/rand/v2"

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

// Algorithm implements the SPEA2 EA with refactored truncation and logging.
type Algorithm struct {
	algos.GeneticAlgorithm // embeds start time, timeout, problem, logger
	params                 Params
	population             []Individual
	archive                []Individual
	generation             int
}

// Step is emitted each generation for logging Pareto front.
type Step struct {
	algos.GeneticAlgorithmStep
	Generation  int         `json:"generation"`
	ParetoFront [][]float64 `json:"pareto_front"`
}

// NewAlgorithm constructs a SPEA2Algorithm.
func NewAlgorithm(
	problem problems.Problem,
	timeout time.Duration,
	params Params,
	logger algos.ProgressLoggerProvider,
) *Algorithm {
	ga := algos.NewGeneticAlgorithm(problem, timeout, logger)
	return &Algorithm{
		GeneticAlgorithm: *ga,
		params:           params,
		generation:       0,
	}
}

// Run executes SPEA2 until timeout or generation limit, logging each generation.
func (alg *Algorithm) Run() {
	alg.initPopulation()
	alg.archive = nil

	for time.Since(alg.StartTimestamp) < alg.Timeout && alg.generation < alg.params.GenerationLimit {
		alg.generation++

		combined := slices.Concat(alg.population, alg.archive)
		alg.assignFitness(combined)
		alg.updateArchive(combined)
		alg.logParetoFront()
		alg.reproduce()
	}
}

// initPopulation initializes the population with random solutions.
func (alg *Algorithm) initPopulation() {
	alg.population = make([]Individual, alg.params.PopulationSize)
	for i := range alg.population {
		alg.population[i] = Individual{sol: alg.Problem.RandomSolution()}
	}
}

// assignFitness computes strength, raw fitness, density, and combined fitness.
func (alg *Algorithm) assignFitness(all []Individual) {
	// strength
	for i := range all {
		all[i].strength = 0
		for j := range all {
			if dominates(all[i].sol, all[j].sol) {
				all[i].strength++
			}
		}
	}
	// raw fit
	for i := range all {
		all[i].rawFit = 0
		for j := range all {
			if dominates(all[j].sol, all[i].sol) {
				all[i].rawFit += all[j].strength
			}
		}
	}
	// density
	for i := range all {
		all[i].density = computeKthDist(all, i, alg.params.DensityKth)
		all[i].density = 1.0 / (all[i].density + 2.0)
	}
	// combined fitness
	for i := range all {
		all[i].fitness = all[i].rawFit + all[i].density
	}
}

// updateArchive performs nondominated selection, filling, and iterative k-NN truncation.
func (alg *Algorithm) updateArchive(combined []Individual) {
	// 1) select nondominated
	nd := filter(combined, func(ind Individual) bool {
		return ind.rawFit < 1
	})
	// 2) fill if too few
	if len(nd) < alg.params.ArchiveSize {
		nd = append(nd, selectDominated(combined, alg.params.ArchiveSize-len(nd))...)
	}
	// 3) truncate if too many
	alg.archive = truncateToSize(nd, alg.params.ArchiveSize, alg.params.DensityKth)
}

// logParetoFront logs the current archive’s Pareto front.
func (alg *Algorithm) logParetoFront() {
	var pareto [][]float64
	for _, ind := range alg.archive {
		if ind.rawFit < 1 {
			pareto = append(pareto, ind.sol.Objectives())
		}
	}
	alg.LogStep(Step{
		GeneticAlgorithmStep: algos.GeneticAlgorithmStep{Elapsed: time.Since(alg.StartTimestamp)},
		Generation:           alg.generation,
		ParetoFront:          pareto,
	})
}

// reproduce creates the next population via binary tournament, crossover, and mutation.
func (alg *Algorithm) reproduce() {
	var nextP []Individual
	for len(nextP) < alg.params.PopulationSize {
		p1 := alg.tournamentSelect()
		p2 := alg.tournamentSelect()
		var kids []problems.Solution
		if rand.Float64() < alg.params.CrossoverProb {
			kids = p1.sol.Crossover(p2.sol)
		} else {
			kids = []problems.Solution{p1.sol, p2.sol}
		}
		for _, k := range kids {
			k = k.Mutate(alg.params.MutationProb)
			nextP = append(nextP, Individual{sol: k})
			if len(nextP) >= alg.params.PopulationSize {
				break
			}
		}
	}
	alg.population = nextP
}

// filter returns elements where keep returns true.
func filter(slice []Individual, keep func(Individual) bool) []Individual {
	var out []Individual
	for _, ind := range slice {
		if keep(ind) {
			out = append(out, ind)
		}
	}
	return out
}

// selectDominated picks the best 'count' dominated individuals.
func selectDominated(all []Individual, count int) []Individual {
	var dom []Individual
	for _, ind := range all {
		if ind.rawFit >= 1 {
			dom = append(dom, ind)
		}
	}
	sort.Slice(dom, func(i, j int) bool {
		return dom[i].fitness < dom[j].fitness
	})
	if count > len(dom) {
		count = len(dom)
	}
	return dom[:count]
}

// truncateToSize iteratively removes the most crowded individual until size is met.
func truncateToSize(cands []Individual, size, k int) []Individual {
	archive := slices.Clone(cands)
	for len(archive) > size {
		idx := mostCrowdedIndex(archive)
		archive = slices.Delete(archive, idx, idx+1)
	}
	return archive
}

// mostCrowdedIndex returns the index with closest neighbor
func mostCrowdedIndex(archive []Individual) int {
	minDist, removeIdx := math.Inf(1), 0
	for i := range archive {
		fi := archive[i].sol.Objectives()
		for j := range archive {
			if i == j {
				continue
			}
			dist := euclidean(fi, archive[j].sol.Objectives())
			if dist < minDist {
				minDist, removeIdx = dist, i
			}
		}
	}
	return removeIdx
}

// computeKthDist returns the distance to the k-th nearest neighbor of archive[i].
func computeKthDist(archive []Individual, i, k int) float64 {
	n := len(archive)
	fi := archive[i].sol.Objectives()
	dists := make([]float64, 0, n-1)
	for j := range archive {
		if i == j {
			continue
		}
		dists = append(dists, euclidean(fi, archive[j].sol.Objectives()))
	}
	sort.Float64s(dists)
	return dists[k]
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

// euclidean computes the Euclidean distance between two objective vectors.
func euclidean(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}
