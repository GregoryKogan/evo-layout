package knapsack

import "math/rand/v2"

type Item struct {
	Value     int   `json:"value"`
	Resources []int `json:"resources"`
}

func NewRandomItem(params KnapsackProblemParams) Item {
	resources := make([]int, params.Dimensions-1)
	for i := range params.Dimensions - 1 {
		resources[i] = 1 + rand.IntN(params.InitialMaxResource-1)
	}
	return Item{Value: rand.IntN(params.InitialMaxValue), Resources: resources}
}
