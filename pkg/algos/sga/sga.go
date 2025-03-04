package sga

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type SimpleGeneticAlgorithm struct {
	algos.GeneticAlgorithm
	params         SGAParams
	generations    int
	population     []problems.Solution
	eliteSize      int
	matingPoolSize int
}

type SimpleGeneticAlgorithmStep struct {
	algos.GeneticAlgorithmStep
	Generation int `json:"generation"`
}

func NewSimpleGeneticAlgorithm(problem problems.Problem, targetFitness float64, params SGAParams, logger algos.ProgressLoggerProvider) *SimpleGeneticAlgorithm {
	return &SimpleGeneticAlgorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, targetFitness, logger),
		params:           params,
		generations:      0,
		eliteSize:        int(float64(params.PopulationSize) * params.ElitePercentile),
		matingPoolSize:   int(float64(params.PopulationSize) * params.MatingPoolPercentile),
	}
}

func (sga *SimpleGeneticAlgorithm) Run() {
	sga.InitPopulation()
	for fitness := 0.0; fitness < sga.TargetFitness; {
		sga.Evolve()
		sga.generations++
		bestFitness := sga.Solution.Fitness()
		fmt.Println("Generation", sga.generations, "Best Fitness", bestFitness)
		if fitness != bestFitness {
			fitness = bestFitness
			sga.LogStep(SimpleGeneticAlgorithmStep{
				GeneticAlgorithmStep: algos.GeneticAlgorithmStep{Elapsed: time.Since(sga.StartTimestamp), Solution: sga.Solution},
				Generation:           sga.generations,
			})
		}
	}
}

func (sga *SimpleGeneticAlgorithm) InitPopulation() {
	pop := make([]problems.Solution, sga.params.PopulationSize)
	for i := range pop {
		pop[i] = sga.Problem.RandomSolution()
	}
	sga.population = pop
}

func (sga *SimpleGeneticAlgorithm) Evolve() {
	sga.evaluateGeneration()

	newPopulation := make([]problems.Solution, 0, sga.params.PopulationSize)

	// perform elitism
	for range sga.eliteSize {
		newPopulation = append(newPopulation, sga.population[len(newPopulation)])
	}

	// generate rest of the population from mating pool
	for len(newPopulation) < sga.params.PopulationSize {
		p1Ind := rand.Intn(sga.matingPoolSize)
		p2Ind := rand.Intn(sga.matingPoolSize)
		if p1Ind == p2Ind {
			continue
		}
		parent1 := sga.population[p1Ind]
		parent2 := sga.population[p2Ind]
		children := parent1.Crossover(parent2)
		for i := range children {
			children[i] = children[i].Mutate(sga.params.MutationRate)
		}
		newPopulation = append(newPopulation, children...)
	}

	sga.population = newPopulation
}

func (sga *SimpleGeneticAlgorithm) evaluateGeneration() {
	// descending fitness
	sort.Slice(sga.population, func(i, j int) bool {
		return sga.population[j].Fitness() < sga.population[i].Fitness()
	})
	sga.Solution = sga.population[0]
}
