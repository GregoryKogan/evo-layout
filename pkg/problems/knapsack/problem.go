package knapsack

import "github.com/GregoryKogan/genetic-algorithms/pkg/problems"

type KnapsackProblem struct {
	Params KnapsackProblemParams `json:"parameters"`
	Items  []Item                `json:"items"`
}

type KnapsackProblemParams struct {
	Dimensions         int   `json:"dimensions"`
	ItemsNum           int   `json:"items_num"`
	InitialMaxValue    int   `json:"initial_max_value"`
	InitialMaxResource int   `json:"initial_max_resource"`
	InitialMaxAmount   int   `json:"initial_max_amount"`
	Constraints        []int `json:"constraints"`
}

func (pp *KnapsackProblemParams) validate() {
	if pp.Dimensions < 2 {
		panic("Knapsack problem is at least 2-dimensional!")
	}
	if pp.ItemsNum < 2 {
		panic("Knapsack problem requires at least 2 items to choose from")
	}
	if len(pp.Constraints) != pp.Dimensions-1 {
		panic("Constraints and dimensions do not match")
	}
}

func NewKnapsackProblem(params KnapsackProblemParams) problems.Problem {
	params.validate()

	items := make([]Item, 0, params.ItemsNum)
	for range params.ItemsNum {
		items = append(items, NewRandomItem(params))
	}
	return &KnapsackProblem{Params: params, Items: items}
}

func (p *KnapsackProblem) Name() string {
	return "Knapsack"
}

func (p *KnapsackProblem) RandomSolution() problems.Solution {
	return RandomKnapsackSolution(p.Params, p.Items)
}
