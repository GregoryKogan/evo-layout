package zdt

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// ZDT6Solution implements the ZDT6 problem.
// f1(x) = 1 - exp(-4 * x₀) * (sin(6πx₀))⁶
// g(x)  = 1 + 9 * ( (∑₍ᵢ₌₁₎ⁿ₋₁ xᵢ)/(n-1) )^0.25
// f2(x) = g * (1 - (f1/g)²)
type ZDT6Solution struct {
	Dimensions       int       `json:"dimensions"`
	X                []float64 `json:"x"`
	CachedObjectives []float64 `json:"objectives"`
}

func RandomZDT6Solution(dimensions int) problems.Solution {
	x := make([]float64, dimensions)
	for i := range x {
		x[i] = rand.Float64()
	}
	return &ZDT6Solution{Dimensions: dimensions, X: x}
}

func (s *ZDT6Solution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}
	n := float64(s.Dimensions)
	f1 := 1 - math.Exp(-4*s.X[0])*math.Pow(math.Sin(6*math.Pi*s.X[0]), 6)
	sum := 0.0
	for i := 1; i < s.Dimensions; i++ {
		sum += s.X[i]
	}
	avg := sum / (n - 1)
	g := 1 + 9*math.Pow(avg, 0.25)
	f2 := g * (1 - math.Pow(f1/g, 2))
	s.CachedObjectives = []float64{f1, f2}
	return s.CachedObjectives
}

func (s *ZDT6Solution) Fitness() float64 {
	obj := s.Objectives()
	return obj[0] + obj[1]
}

func (s *ZDT6Solution) Crossover(other problems.Solution) []problems.Solution {
	otherSol, ok := other.(*ZDT6Solution)
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
		&ZDT6Solution{Dimensions: s.Dimensions, X: child1},
		&ZDT6Solution{Dimensions: s.Dimensions, X: child2},
	}
}

func (s *ZDT6Solution) Mutate(rate float64) problems.Solution {
	mutant := make([]float64, s.Dimensions)
	copy(mutant, s.X)
	for i := 0; i < s.Dimensions; i++ {
		if rand.Float64() < rate {
			delta := rand.NormFloat64() * 0.1
			mutant[i] += delta
			if mutant[i] < 0 {
				mutant[i] = 0
			}
			if mutant[i] > 1 {
				mutant[i] = 1
			}
		}
	}
	return &ZDT6Solution{Dimensions: s.Dimensions, X: mutant}
}
