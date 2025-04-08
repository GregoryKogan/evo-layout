package ssga

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type Algorithm struct {
	algos.GeneticAlgorithm
	params     Params
	iterations int
	population []problems.Solution
}

type Step struct {
	algos.GeneticAlgorithmStep
	Iteration int `json:"iteration"`
}

func NewAlgorithm(
	problem problems.Problem,
	targetFitness float64,
	params Params,
	logger algos.ProgressLoggerProvider,
) *Algorithm {
	return &Algorithm{
		GeneticAlgorithm: *algos.NewGeneticAlgorithm(problem, targetFitness, logger),
		params:           params,
		iterations:       0,
	}
}

func (alg *Algorithm) Run() {
	alg.InitPopulation()
	for fitness := 0.0; fitness < alg.TargetFitness; {
		alg.Evolve()
		alg.iterations++
		bestFitness := alg.Solution.Fitness()
		fmt.Println("Iteration", alg.iterations, "Best Fitness", bestFitness)
		if fitness != bestFitness {
			fitness = bestFitness
			alg.LogStep(Step{
				GeneticAlgorithmStep: algos.GeneticAlgorithmStep{Elapsed: time.Since(alg.StartTimestamp), Solution: alg.Solution},
				Iteration:            alg.iterations,
			})
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
		children := parent1.Crossover(parent2)
		for i := range children {
			children[i] = children[i].Mutate(alg.params.MutationRate)
			alg.population[alg.params.PopulationSize-i-1] = children[i]
		}
		replaced = true
	}
}

func (alg *Algorithm) evaluateGeneration() {
	// descending fitness
	sort.Slice(alg.population, func(i, j int) bool {
		return alg.population[j].Fitness() < alg.population[i].Fitness()
	})
	alg.Solution = alg.population[0]
}

func (alg *Algorithm) tournamentSelect() int {
	ind1 := rand.IntN(alg.params.PopulationSize)
	ind2 := rand.IntN(alg.params.PopulationSize)
	if ind1 == ind2 {
		return ind1
	}
	if alg.population[ind1].Fitness() > alg.population[ind2].Fitness() {
		return ind1
	}
	return ind2
}
