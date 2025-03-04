package knapsack

import (
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type KnapsackProblem struct {
	Params KnapsackProblemParams `json:"parameters"`
	Items  []Item                `json:"items"`
}

type KnapsackProblemParams struct {
	Dimensions         int   `json:"dimensions"`
	ItemsNum           int   `json:"items_num"`
	InitialMaxValue    int   `json:"initial_max_value"`
	InitialMaxResource int   `json:"initial_max_resource"`
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

func NewKnapsackProblem(params KnapsackProblemParams) problems.AlgorithmicProblem {
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

func (p *KnapsackProblem) AlgorithmicSolution() problems.AlgorithmicSolution {
	start := time.Now()

	if p.Params.Dimensions != 2 {
		panic("Algorithmic solutions available only for 2D knapsack problem")
	}

	n := p.Params.ItemsNum
	capacity := p.Params.Constraints[0]

	weights := make([]int, n)
	values := make([]int, n)
	for i := range n {
		values[i] = p.Items[i].Value
		weights[i] = p.Items[i].Resources[0]
	}

	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	// solving
	for i := 1; i <= n; i++ {
		for w := 1; w <= capacity; w++ {
			if weights[i-1] <= w {
				include_item := values[i-1] + dp[i-1][w-weights[i-1]]
				exclude_item := dp[i-1][w]
				dp[i][w] = max(include_item, exclude_item)
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	// tracing solution
	selectedBits := make([]bool, n)
	for i, w := n, capacity; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			selectedBits[i-1] = true
			w -= weights[i-1]
		}
	}

	return problems.AlgorithmicSolution{
		Solution: &KnapsackSolution{problemParams: p.Params, items: p.Items, Bits: selectedBits},
		TimeTook: time.Since(start),
	}
}
