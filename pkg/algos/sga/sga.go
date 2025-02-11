package sga

import (
	"fmt"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type SimpleGeneticAlgorithm struct {
	algos.GeneticAlgorithm
	popSize     int
	generations int
	bestFitness float64
	population  []problems.Solution
}

func NewSimpleGeneticAlgorithm(problem problems.Problem, targetFitness float64, popSize int) *SimpleGeneticAlgorithm {
	return &SimpleGeneticAlgorithm{
		GeneticAlgorithm: algos.GeneticAlgorithm{Problem: problem, TargetFitness: targetFitness},
		popSize:          popSize,
		generations:      0,
		bestFitness:      0,
	}
}

func (sga *SimpleGeneticAlgorithm) Run() {
	sga.InitPopulation()
	for sga.bestFitness < sga.TargetFitness {
		fmt.Println("Generation", sga.generations, "Best Fitness", sga.bestFitness)
		sga.Evolve()
		sga.generations++
	}
}

func (sga *SimpleGeneticAlgorithm) InitPopulation() {
	pop := make([]problems.Solution, sga.popSize)
	for i := range pop {
		pop[i] = sga.Problem.RandomSolution()
	}
	sga.population = pop
}

func (sga *SimpleGeneticAlgorithm) Evolve() {
	// evaluate current fitness and determine best solution
	bestIndex := 0
	bestFit := sga.population[0].Fitness()
	for i, sol := range sga.population {
		fit := sol.Fitness()
		if fit > bestFit {
			bestFit = fit
			bestIndex = i
		}
	}
	sga.bestFitness = bestFit

	// create new population with elitism
	newPop := make([]problems.Solution, sga.popSize)
	newPop[0] = sga.population[bestIndex] // carry over best solution

	// generate rest of the population
	for i := 1; i < sga.popSize; i++ {
		// tournament selection: pick two random solutions and choose the better one
		a := sga.population[rand.Intn(len(sga.population))]
		b := sga.population[rand.Intn(len(sga.population))]
		parent1 := a
		if a.Fitness() < b.Fitness() {
			parent1 = b
		}
		// select second parent
		c := sga.population[rand.Intn(len(sga.population))]
		d := sga.population[rand.Intn(len(sga.population))]
		parent2 := c
		if c.Fitness() < d.Fitness() {
			parent2 = d
		}
		// create child through crossover and then mutation
		child := parent1.Crossover(parent2)
		child = child.Mutate()
		newPop[i] = child
	}
	sga.population = newPop
}
