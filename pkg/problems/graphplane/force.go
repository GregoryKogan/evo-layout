package graphplane

import (
	"fmt"
	"math"
	"time"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type FDSParams struct {
	Steps int
	Temp  float64
	K     float64
}

// ForceDirectedSolver applies a Fruchterman–Reingold layout to the Graph.
type ForceDirectedSolver struct {
	*GraphPlaneSolution
	logger      algos.ProgressLoggerProvider
	adj         [][]int
	params      FDSParams
	k           float64
	temp        float64
	coolingStep float64
}

func NewForceDirectedSolver(initialSolution problems.Solution, params FDSParams, logger algos.ProgressLoggerProvider) ForceDirectedSolver {
	gpSol, ok := initialSolution.(*GraphPlaneSolution)
	if !ok {
		fmt.Printf("%#+v\n", initialSolution)
		panic("type error")
	}
	return ForceDirectedSolver{
		GraphPlaneSolution: gpSol,
		logger:             logger,
		params:             params,
	}
}

// Solve runs the spring-electrical simulation and returns a solution.
func (s *ForceDirectedSolver) Solve() problems.AlgorithmicSolution {
	start := time.Now()

	n := s.Graph.NumVertices
	s.k = math.Sqrt((s.Height*s.Width)/float64(n)) * s.params.K
	s.temp = math.Min(s.Height, s.Width) * s.params.Temp
	s.coolingStep = s.temp / float64(s.params.Steps)

	// precompute adjacency lists
	s.adj = make([][]int, n)
	for _, e := range s.Graph.Edges {
		s.adj[e.From] = append(s.adj[e.From], e.To)
		s.adj[e.To] = append(s.adj[e.To], e.From)
	}

	for step := range s.params.Steps {
		s.Iterate()

		s.CachedObjectives = nil
		s.Fitness()

		if s.logger != nil {
			s.logger.LogStep(algos.GAStep{Elapsed: time.Since(start), Solution: s.GraphPlaneSolution, Step: step + 1})
		}
	}

	return problems.AlgorithmicSolution{Solution: s.GraphPlaneSolution, TimeTook: time.Since(start)}
}

func (s *ForceDirectedSolver) Iterate() {
	n := s.Graph.NumVertices

	// displacement vectors
	disp := make([]VertexPos, n)

	// repulsive forces between all pairs
	for i := range n {
		for j := i + 1; j < n; j++ {
			dx := s.VertPositions[i].X - s.VertPositions[j].X
			dy := s.VertPositions[i].Y - s.VertPositions[j].Y
			d := math.Hypot(dx, dy) + 1e-9
			force := (s.k * s.k) / (d * d)
			disp[i].X += dx * force
			disp[i].Y += dy * force
			disp[j].X -= dx * force
			disp[j].Y -= dy * force
		}
	}

	// attractive forces along edges
	for u, neigh := range s.adj {
		for _, v := range neigh {
			dx := s.VertPositions[u].X - s.VertPositions[v].X
			dy := s.VertPositions[u].Y - s.VertPositions[v].Y
			d := math.Hypot(dx, dy) + 1e-9
			force := (d * d) / s.k
			dxNorm := dx / d
			dyNorm := dy / d
			disp[u].X -= dxNorm * force
			disp[u].Y -= dyNorm * force
			disp[v].X += dxNorm * force
			disp[v].Y += dyNorm * force
		}
	}

	for i := range n {
		dx := disp[i].X
		dy := disp[i].Y
		disp := math.Hypot(dx, dy)

		if disp > 0 {
			scale := math.Min(disp, s.temp) / disp
			dx *= scale
			dy *= scale
		}

		s.VertPositions[i].X = clamp(s.VertPositions[i].X+dx, 0, s.Width)
		s.VertPositions[i].Y = clamp(s.VertPositions[i].Y+dy, 0, s.Height)
	}

	s.temp -= s.coolingStep
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
