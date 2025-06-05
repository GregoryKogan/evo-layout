package ssga

import (
	"context"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type Algorithm struct {
	algos.GeneticAlgorithm
	params     Params
	generation int
	population []problems.Solution
}

func NewAlgorithm(
	problem problems.Problem,
	params Params,
	generationLimit int,
	logger algos.ProgressLoggerProvider,
) *Algorithm {
	return &Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, generationLimit, logger),
		params:           params,
		generation:       0,
	}
}

func (alg *Algorithm) Run(ctx context.Context) {
	alg.InitPopulation()
	fitness := 0.0
	for alg.generation < alg.GenerationLimit {
		select {
		case <-ctx.Done():
			return
		default:
		}
		alg.Evolve()
		alg.generation++
		bestFitness := alg.Solution.Fitness()
		if fitness != bestFitness {
			fitness = bestFitness
			if alg.ProgressLoggerProvider != nil {
				alg.LogStep(algos.GAStep{Elapsed: time.Since(alg.StartTimestamp), Solution: alg.Solution, Step: alg.generation})
			}
		}
	}
}

func (alg *Algorithm) InitPopulation() {
	pop := make([]problems.Solution, alg.params.PopulationSize)
	for i := range pop {
		pop[i] = alg.Problem.RandomSolution()
	}
	alg.population = pop
}

func (alg *Algorithm) Seed(seedSolution problems.Solution) {
	alg.population = make([]problems.Solution, alg.params.PopulationSize)
	for i := range alg.params.PopulationSize {
		alg.population[i] = alg.params.MutationFunc(seedSolution)
	}
	alg.population[0] = seedSolution
	alg.Solution = seedSolution
}

func (alg *Algorithm) GetPopulation() []problems.Solution {
	pop := make([]problems.Solution, len(alg.population))
	copy(pop, alg.population)
	return pop
}

func (alg *Algorithm) GetSteps() int {
	return alg.generation
}

func (alg *Algorithm) Evolve() {
	alg.evaluateGeneration()

	replaced := false
	for !replaced {
		p1Ind := alg.tournamentSelect()
		p2Ind := alg.tournamentSelect()
		if p1Ind == p2Ind {
			continue
		}
		parent1 := alg.population[p1Ind]
		parent2 := alg.population[p2Ind]

		children := alg.params.CrossoverFunc(parent1, parent2)

		for i := range children {
			children[i] = alg.params.MutationFunc(children[i])
			alg.population[alg.params.PopulationSize-i-1] = children[i]
		}
		replaced = true
	}
}

func (alg *Algorithm) evaluateGeneration() {
	// descending fitness
	sort.Slice(alg.population, func(i, j int) bool {
		return alg.population[i].Fitness() < alg.population[j].Fitness()
	})
	alg.Solution = alg.population[0]
}

func (alg *Algorithm) tournamentSelect() int {
	ind1 := rand.IntN(alg.params.PopulationSize)
	ind2 := rand.IntN(alg.params.PopulationSize)
	if ind1 == ind2 {
		return ind1
	}
	if alg.population[ind1].Fitness() < alg.population[ind2].Fitness() {
		return ind1
	}
	return ind2
}
