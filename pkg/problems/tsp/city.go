package tsp

import (
	"math"
	"math/rand/v2"
)

type City struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func NewRandomCity() City {
	return City{
		Latitude:  rand.Float64() * 100,
		Longitude: rand.Float64() * 100,
	}
}

func (c *City) Distance(other City) float64 {
	return math.Hypot(c.Latitude-other.Latitude, c.Longitude-other.Longitude)
}
