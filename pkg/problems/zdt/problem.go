package zdt

import (
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// --------------------
// ZDT1 Problem
// --------------------
type ZDT1Problem struct {
	Dimensions int `json:"dimensions"`
}

// NewZDT1Problem creates a new ZDT1 problem instance with the specified number of dimensions.
func NewZDT1Problem(dimensions int) problems.Problem {
	if dimensions < 2 {
		panic("ZDT problems require at least 2 dimensions")
	}
	return &ZDT1Problem{Dimensions: dimensions}
}

func (p *ZDT1Problem) Name() string {
	return "ZDT1"
}

func (p *ZDT1Problem) RandomSolution() problems.Solution {
	return RandomZDT1Solution(p.Dimensions)
}

// --------------------
// ZDT2 Problem
// --------------------
type ZDT2Problem struct {
	Dimensions int `json:"dimensions"`
}

// NewZDT2Problem creates a new ZDT2 problem instance.
func NewZDT2Problem(dimensions int) problems.Problem {
	if dimensions < 2 {
		panic("ZDT problems require at least 2 dimensions")
	}
	return &ZDT2Problem{Dimensions: dimensions}
}

func (p *ZDT2Problem) Name() string {
	return "ZDT2"
}

func (p *ZDT2Problem) RandomSolution() problems.Solution {
	return RandomZDT2Solution(p.Dimensions)
}

// --------------------
// ZDT3 Problem
// --------------------
type ZDT3Problem struct {
	Dimensions int `json:"dimensions"`
}

// NewZDT3Problem creates a new ZDT3 problem instance.
func NewZDT3Problem(dimensions int) problems.Problem {
	if dimensions < 2 {
		panic("ZDT problems require at least 2 dimensions")
	}
	return &ZDT3Problem{Dimensions: dimensions}
}

func (p *ZDT3Problem) Name() string {
	return "ZDT3"
}

func (p *ZDT3Problem) RandomSolution() problems.Solution {
	return RandomZDT3Solution(p.Dimensions)
}

// --------------------
// ZDT4 Problem
// --------------------
type ZDT4Problem struct {
	Dimensions int `json:"dimensions"`
}

// NewZDT4Problem creates a new ZDT4 problem instance.
// Note: ZDT4 requires at least 2 dimensions. The first decision variable is in [0,1] and the rest in [-5,5].
func NewZDT4Problem(dimensions int) problems.Problem {
	if dimensions < 2 {
		panic("ZDT4 requires at least 2 dimensions")
	}
	return &ZDT4Problem{Dimensions: dimensions}
}

func (p *ZDT4Problem) Name() string {
	return "ZDT4"
}

func (p *ZDT4Problem) RandomSolution() problems.Solution {
	return RandomZDT4Solution(p.Dimensions)
}

// --------------------
// ZDT6 Problem
// --------------------
type ZDT6Problem struct {
	Dimensions int `json:"dimensions"`
}

// NewZDT6Problem creates a new ZDT6 problem instance.
func NewZDT6Problem(dimensions int) problems.Problem {
	if dimensions < 2 {
		panic("ZDT problems require at least 2 dimensions")
	}
	return &ZDT6Problem{Dimensions: dimensions}
}

func (p *ZDT6Problem) Name() string {
	return "ZDT6"
}

func (p *ZDT6Problem) RandomSolution() problems.Solution {
	return RandomZDT6Solution(p.Dimensions)
}
