package sga

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
	params         Params
	generation     int
	population     []problems.Solution
	eliteSize      int
	matingPoolSize int
}

func NewAlgorithm(problem problems.Problem, params Params, generationLimit int, logger algos.ProgressLoggerProvider) *Algorithm {
	return &Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, generationLimit, logger),
		params:           params,
		generation:       0,
		eliteSize:        int(float64(params.PopulationSize) * params.ElitePercentile),
		matingPoolSize:   int(float64(params.PopulationSize) * params.MatingPoolPercentile),
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
				alg.LogStep(algos.GAStep{
					Elapsed: time.Since(alg.StartTimestamp), Solution: alg.Solution, Step: alg.generation,
				})
			}
		}
	}
}

func (alg *Algorithm) GetSteps() int {
	return alg.generation
}

func (alg *Algorithm) InitPopulation() {
	pop := make([]problems.Solution, alg.params.PopulationSize)
	for i := range pop {
		pop[i] = alg.Problem.RandomSolution()
	}
	alg.population = pop
}

func (alg *Algorithm) Evolve() {
	alg.evaluateGeneration()

	newPopulation := make([]problems.Solution, 0, alg.params.PopulationSize)

	// perform elitism
	newPopulation = append(newPopulation, alg.population[:alg.eliteSize]...)

	// generate rest of the population
	for len(newPopulation) < alg.params.PopulationSize {
		p1Ind := rand.IntN(alg.matingPoolSize)
		p2Ind := rand.IntN(alg.matingPoolSize)
		if p1Ind == p2Ind {
			continue
		}
		parent1 := alg.population[p1Ind]
		parent2 := alg.population[p2Ind]

		children := alg.params.CrossoverFunc(parent1, parent2)

		for _, child := range children {
			child = alg.params.MutationFunc(child)
			newPopulation = append(newPopulation, child)
			if len(newPopulation) >= alg.params.PopulationSize {
				break
			}
		}
	}

	alg.population = newPopulation
}

func (alg *Algorithm) evaluateGeneration() {
	sort.Slice(alg.population, func(i, j int) bool {
		return alg.population[i].Fitness() < alg.population[j].Fitness()
	})
	alg.Solution = alg.population[0]
}
