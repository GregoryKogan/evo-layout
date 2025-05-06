package zdt

import (
	"math"
	"math/rand/v2"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// ZDT1Solution represents a candidate solution for the ZDT1 problem.
type ZDT1Solution struct {
	Dimensions       int       `json:"dimensions"`
	X                []float64 `json:"x"` // decision variables in [0,1]
	CachedObjectives []float64 `json:"objectives"`
}

// RandomZDT1Solution returns a new random solution for a given dimensionality.
func RandomZDT1Solution(dimensions int) problems.Solution {
	x := make([]float64, dimensions)
	for i := range x {
		x[i] = rand.Float64() // uniformly in [0,1]
	}
	return &ZDT1Solution{
		Dimensions: dimensions,
		X:          x,
	}
}

// Objectives computes and returns the two objectives for ZDT1:
// f1 = x[0]
// g  = 1 + (9/(n-1)) * sum(x[1:n])
// f2 = g * (1 - sqrt(x[0]/g))
func (s *ZDT1Solution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}
	n := float64(s.Dimensions)
	f1 := s.X[0]
	sum := 0.0
	for i := 1; i < s.Dimensions; i++ {
		sum += s.X[i]
	}
	g := 1.0 + (9.0/(n-1))*sum
	f2 := g * (1 - math.Sqrt(f1/g))
	s.CachedObjectives = []float64{f1, f2}
	return s.CachedObjectives
}

// Fitness computes a scalar fitness value as the sum of the two objectives (minimization).
func (s *ZDT1Solution) Fitness() float64 {
	objs := s.Objectives()
	return objs[0] + objs[1]
}

func ZDT1CrossoverFunc() problems.CrossoverFunc {
	return func(parentA, parentB problems.Solution) []problems.Solution {
		a, aOk := parentA.(*ZDT1Solution)
		b, bOk := parentB.(*ZDT1Solution)
		if !aOk || !bOk {
			panic("invalid parents")
		}
		return a.crossover(b)
	}
}

func ZDT1MutationFunc() problems.MutationFunc {
	return func(individual problems.Solution) problems.Solution {
		s, ok := individual.(*ZDT1Solution)
		if !ok {
			panic("invalid individual")
		}
		return s.mutate()
	}
}

// Crossover applies a uniform crossover between two ZDT solutions and returns two offspring.
func (s *ZDT1Solution) crossover(other problems.Solution) []problems.Solution {
	otherZDT, ok := other.(*ZDT1Solution)
	if !ok || s.Dimensions != otherZDT.Dimensions {
		return []problems.Solution{s}
	}

	child1X := make([]float64, s.Dimensions)
	child2X := make([]float64, s.Dimensions)
	for i := range s.Dimensions {
		if rand.Float64() < 0.5 {
			child1X[i] = s.X[i]
			child2X[i] = otherZDT.X[i]
		} else {
			child1X[i] = otherZDT.X[i]
			child2X[i] = s.X[i]
		}
	}
	child1 := &ZDT1Solution{Dimensions: s.Dimensions, X: child1X}
	child2 := &ZDT1Solution{Dimensions: s.Dimensions, X: child2X}
	return []problems.Solution{child1, child2}
}

// mutate applies mutation by perturbing each decision variable with a small probability.
// The mutated value is clamped to remain within [0, 1].
func (s *ZDT1Solution) mutate() problems.Solution {
	mutantX := make([]float64, s.Dimensions)
	copy(mutantX, s.X)
	for i := range s.Dimensions {
		// Add a normally distributed perturbation (scale factor is arbitrary).
		delta := rand.NormFloat64() * 0.1
		mutantX[i] += delta
		// Clamp to [0, 1].
		if mutantX[i] < 0 {
			mutantX[i] = 0
		}
		if mutantX[i] > 1 {
			mutantX[i] = 1
		}
	}
	return &ZDT1Solution{
		Dimensions: s.Dimensions,
		X:          mutantX,
	}
}
