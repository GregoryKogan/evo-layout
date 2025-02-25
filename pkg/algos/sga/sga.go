package sga

import (
	"fmt"
	"math/rand"
	"sort"

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
	Generation int               `json:"generation"`
	Solution   problems.Solution `json:"solution"`
}

func NewSimpleGeneticAlgorithm(problem problems.Problem, targetFitness float64, logFilepath string, params SGAParams) *SimpleGeneticAlgorithm {
	sga := &SimpleGeneticAlgorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, targetFitness, logFilepath),
		params:           params,
		generations:      0,
		eliteSize:        int(float64(params.PopulationSize) * params.ElitePercentile),
		matingPoolSize:   int(float64(params.PopulationSize) * params.MatingPoolPercentile),
	}

	sga.InitLogging()
	sga.LogProblem(problem)

	return sga
}

func (sga *SimpleGeneticAlgorithm) Run() {
	sga.InitPopulation()
	for sga.Solution.Fitness() < sga.TargetFitness {
		sga.Evolve()
		sga.generations++
		fmt.Println("Generation", sga.generations, "Best Fitness", sga.Solution.Fitness())
		sga.LogStep(SimpleGeneticAlgorithmStep{Generation: sga.generations, Solution: sga.Solution})
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
