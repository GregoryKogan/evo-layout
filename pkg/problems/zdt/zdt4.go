package zdt

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// ZDT4Solution implements the ZDT4 problem.
// f1(x) = x₀
// g(x)  = 1 + 10*(n-1) + ∑₍ᵢ₌₁₎ⁿ₋₁ (xᵢ² - 10*cos(4πxᵢ))
// f2(x) = g * (1 - √(x₀/g))
type ZDT4Solution struct {
	Dimensions       int       `json:"dimensions"`
	X                []float64 `json:"x"`
	CachedObjectives []float64 `json:"objectives"`
}

func RandomZDT4Solution(dimensions int) problems.Solution {
	if dimensions < 2 {
		panic("ZDT4 requires at least 2 dimensions")
	}
	x := make([]float64, dimensions)
	x[0] = rand.Float64() // in [0,1]
	for i := 1; i < dimensions; i++ {
		// For ZDT4, remaining variables are in [-5,5]
		x[i] = -5 + rand.Float64()*10
	}
	return &ZDT4Solution{Dimensions: dimensions, X: x}
}

func (s *ZDT4Solution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}
	n := float64(s.Dimensions)
	f1 := s.X[0]
	sum := 0.0
	for i := 1; i < s.Dimensions; i++ {
		sum += s.X[i]*s.X[i] - 10*math.Cos(4*math.Pi*s.X[i])
	}
	g := 1 + 10*(n-1) + sum
	f2 := g * (1 - math.Sqrt(f1/g))
	s.CachedObjectives = []float64{f1, f2}
	return s.CachedObjectives
}

func (s *ZDT4Solution) Fitness() float64 {
	obj := s.Objectives()
	return obj[0] + obj[1]
}

func (s *ZDT4Solution) Crossover(other problems.Solution) []problems.Solution {
	otherSol, ok := other.(*ZDT4Solution)
	if !ok || s.Dimensions != otherSol.Dimensions {
		return []problems.Solution{s}
	}
	child1 := make([]float64, s.Dimensions)
	child2 := make([]float64, s.Dimensions)
	for i := 0; i < s.Dimensions; i++ {
		if rand.Float64() < 0.5 {
			child1[i] = s.X[i]
			child2[i] = otherSol.X[i]
		} else {
			child1[i] = otherSol.X[i]
			child2[i] = s.X[i]
		}
	}
	return []problems.Solution{
		&ZDT4Solution{Dimensions: s.Dimensions, X: child1},
		&ZDT4Solution{Dimensions: s.Dimensions, X: child2},
	}
}

func (s *ZDT4Solution) Mutate(rate float64) problems.Solution {
	mutant := make([]float64, s.Dimensions)
	copy(mutant, s.X)
	// Mutate first variable (in [0,1])
	if rand.Float64() < rate {
		delta := rand.NormFloat64() * 0.1
		mutant[0] += delta
		if mutant[0] < 0 {
			mutant[0] = 0
		}
		if mutant[0] > 1 {
			mutant[0] = 1
		}
	}
	// Mutate remaining variables (in [-5,5])
	for i := 1; i < s.Dimensions; i++ {
		if rand.Float64() < rate {
			delta := rand.NormFloat64() * 0.1 * 10
			mutant[i] += delta
			if mutant[i] < -5 {
				mutant[i] = -5
			}
			if mutant[i] > 5 {
				mutant[i] = 5
			}
		}
	}
	return &ZDT4Solution{Dimensions: s.Dimensions, X: mutant}
}
