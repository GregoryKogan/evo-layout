package sga

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// New types for progress logging
type BestSolutionRecord struct {
	Generation int               `json:"generation"`
	Solution   problems.Solution `json:"solution"`
}

type ProgressLog struct {
	Problem       problems.Problem     `json:"problem"`
	BestSolutions []BestSolutionRecord `json:"bestSolutions"`
}

type SimpleGeneticAlgorithm struct {
	algos.GeneticAlgorithm
	popSize     int
	generations int
	bestFitness float64
	population  []problems.Solution

	// new fields for progress logging
	progressLog  ProgressLog
	bestSolution problems.Solution
}

func NewSimpleGeneticAlgorithm(problem problems.Problem, targetFitness float64, popSize int) *SimpleGeneticAlgorithm {
	sga := &SimpleGeneticAlgorithm{
		GeneticAlgorithm: algos.GeneticAlgorithm{Problem: problem, TargetFitness: targetFitness},
		popSize:          popSize,
		generations:      0,
		bestFitness:      0,
	}
	// Initialize progress log with problem at the top
	sga.progressLog = ProgressLog{
		Problem:       problem,
		BestSolutions: []BestSolutionRecord{},
	}
	return sga
}

func (sga *SimpleGeneticAlgorithm) Run() {
	sga.InitPopulation()
	// Write initial log (with the problem) to the file
	sga.logProgress()
	for sga.bestFitness < sga.TargetFitness {
		fmt.Println("Generation", sga.generations, "Best Fitness", sga.bestFitness)
		sga.Evolve()
		sga.generations++
		// Append the best solution for this generation
		sga.progressLog.BestSolutions = append(sga.progressLog.BestSolutions, BestSolutionRecord{
			Generation: sga.generations,
			Solution:   sga.bestSolution,
		})
		sga.logProgress()
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
	sga.bestSolution = sga.population[bestIndex]

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

// New method: logProgress writes the current progress to a JSON file using a relative path.
func (sga *SimpleGeneticAlgorithm) logProgress() {
	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sga.progressLog); err != nil {
		fmt.Println("Error writing log:", err)
	}
}
