package graphplane

import (
	"fmt"
	"math"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// ForceDirectedLayer applies a Fruchterman–Reingold layout to the Graph.
type ForceDirectedLayer struct {
	*GraphPlaneSolution
	logger      algos.ProgressLoggerProvider
	Iterations  int     // number of simulation steps
	SpringK     float64 // ideal edge length
	InitialTemp float64 // starting temperature
}

type Step struct {
	algos.GeneticAlgorithmStep
	Iteration int `json:"iteration"`
}

// NewForceDirectedProblem constructs a force-directed solver.
func NewForceDirectedLayer(sol problems.Solution, iterations int, initialTemp float64, logger algos.ProgressLoggerProvider) ForceDirectedLayer {
	gpSol, ok := sol.(*GraphPlaneSolution)
	if !ok {
		fmt.Printf("%#+v\n", sol)
		panic("type error")
	}
	return ForceDirectedLayer{
		GraphPlaneSolution: gpSol,
		logger:             logger,
		Iterations:         iterations,
		SpringK:            (math.Min(gpSol.Width, gpSol.Height) / math.Sqrt(float64(gpSol.Graph.NumVertices))) * 0.7,
		InitialTemp:        initialTemp,
	}
}

// Solve runs the spring-electrical simulation and returns a solution.
func (p *ForceDirectedLayer) Solve() problems.AlgorithmicSolution {
	bestSolution := &GraphPlaneSolution{Graph: p.Graph, Width: p.Width, Height: p.Height, VertPositions: make([]VertexPos, len(p.VertPositions))}
	copy(bestSolution.VertPositions, p.VertPositions)
	bestSolution.Objectives()

	start := time.Now()
	n := p.Graph.NumVertices

	// precompute adjacency lists
	adj := make([][]int, n)
	for _, e := range p.Graph.Edges {
		adj[e.From] = append(adj[e.From], e.To)
		adj[e.To] = append(adj[e.To], e.From)
	}

	// temperature schedule
	temp := p.InitialTemp
	deltaTemp := p.InitialTemp / float64(p.Iterations)

	// simulation loop
	for iter := range p.Iterations {
		// displacement vectors
		disp := make([]VertexPos, n)

		// repulsive forces between all pairs
		for i := range n {
			for j := i + 1; j < n; j++ {
				dx := p.VertPositions[i].X - p.VertPositions[j].X
				dy := p.VertPositions[i].Y - p.VertPositions[j].Y
				d := math.Hypot(dx, dy) + 1e-9
				force := (p.SpringK * p.SpringK) / d
				dxNorm := dx / d
				dyNorm := dy / d
				disp[i].X += dxNorm * force
				disp[i].Y += dyNorm * force
				disp[j].X -= dxNorm * force
				disp[j].Y -= dyNorm * force
			}
		}

		// attractive forces along edges
		for u, neigh := range adj {
			for _, v := range neigh {
				dx := p.VertPositions[u].X - p.VertPositions[v].X
				dy := p.VertPositions[u].Y - p.VertPositions[v].Y
				d := math.Hypot(dx, dy) + 1e-9
				force := (d * d) / p.SpringK
				dxNorm := dx / d
				dyNorm := dy / d
				disp[u].X -= dxNorm * force
				disp[u].Y -= dyNorm * force
				// symmetric for v
				disp[v].X += dxNorm * force
				disp[v].Y += dyNorm * force
			}
		}

		// limit displacement by temperature and update
		for i := range n {
			dMag := math.Hypot(disp[i].X, disp[i].Y)
			if dMag > 0 {
				scale := math.Min(dMag, temp) / dMag
				px := p.VertPositions[i].X + disp[i].X*scale
				py := p.VertPositions[i].Y + disp[i].Y*scale
				// keep within bounds
				p.VertPositions[i].X = math.Min(p.Width, math.Max(0, px))
				p.VertPositions[i].Y = math.Min(p.Height, math.Max(0, py))
			}
		}

		// cool down
		temp -= deltaTemp

		p.CachedObjectives = nil
		p.Objectives()
		if p.Intersections < bestSolution.Intersections {
			copy(bestSolution.VertPositions, p.VertPositions)
			bestSolution.CachedObjectives = nil
			bestSolution.Objectives()
		}
		p.logger.LogStep(Step{algos.GeneticAlgorithmStep{Elapsed: time.Since(start), Solution: p.GraphPlaneSolution}, iter})
	}

	p.logger.LogStep(Step{algos.GeneticAlgorithmStep{Elapsed: time.Since(start), Solution: bestSolution}, p.Iterations})
	return problems.AlgorithmicSolution{Solution: bestSolution, TimeTook: time.Since(start)}
}
