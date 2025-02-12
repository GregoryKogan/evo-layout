package graphplane

import (
	"math"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneSolution struct {
	graph         *Graph
	width, height float64
	VertPositions []VertexPos `json:"vertices"`
}

type VertexPos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func RandomGraphPlaneSolution(g *Graph, width, height float64) problems.Solution {
	numVertices := g.NumVertices
	vertPositions := make([]VertexPos, numVertices)
	for i := range vertPositions {
		vertPositions[i].X = rand.Float64() * width
		vertPositions[i].Y = rand.Float64() * height
	}
	return &GraphPlaneSolution{g, width, height, vertPositions}
}

func (s *GraphPlaneSolution) Crossover(other problems.Solution) problems.Solution {
	otherGPS, ok := other.(*GraphPlaneSolution)
	if !ok {
		return s
	}
	child := &GraphPlaneSolution{
		VertPositions: make([]VertexPos, s.graph.NumVertices),
		width:         s.width,
		height:        s.height,
		graph:         s.graph, // assume same graph
	}
	// For each vertex choose allele from one of the parents.
	for i := 0; i < s.graph.NumVertices; i++ {
		if rand.Float64() < 0.5 {
			child.VertPositions[i] = s.VertPositions[i]
		} else {
			child.VertPositions[i] = otherGPS.VertPositions[i]
		}
	}
	return child
}

func (s *GraphPlaneSolution) Mutate() problems.Solution {
	mutant := &GraphPlaneSolution{
		VertPositions: make([]VertexPos, s.graph.NumVertices),
		width:         s.width,
		height:        s.height,
		graph:         s.graph,
	}
	copy(mutant.VertPositions, s.VertPositions)
	// pick a random vertex and change its position by a small delta
	i := rand.Intn(s.graph.NumVertices)
	deltaX := (rand.Float64() - 0.5) * s.width * 1
	deltaY := (rand.Float64() - 0.5) * s.height * 1
	mutant.VertPositions[i].X = math.Max(0, math.Min(s.width, mutant.VertPositions[i].X+deltaX))
	mutant.VertPositions[i].Y = math.Max(0, math.Min(s.height, mutant.VertPositions[i].Y+deltaY))
	return mutant
}

func (s *GraphPlaneSolution) Fitness() float64 {
	// count intersections between all pairs of non-adjacent edges
	count := 0
	edges := s.graph.Edges
	// helper function for ccw test
	ccw := func(ax, ay, bx, by, cx, cy float64) bool {
		return (cy-ay)*(bx-ax) > (by-ay)*(cx-ax)
	}
	segmentsIntersect := func(ax, ay, bx, by, cx, cy, dx, dy float64) bool {
		return ccw(ax, ay, cx, cy, dx, dy) != ccw(bx, by, cx, cy, dx, dy) &&
			ccw(ax, ay, bx, by, cx, cy) != ccw(ax, ay, bx, by, dx, dy)
	}
	n := len(edges)
	for i := 0; i < n; i++ {
		a := edges[i].From
		b := edges[i].To
		ax, ay := s.VertPositions[a].X, s.VertPositions[a].Y
		bx, by := s.VertPositions[b].X, s.VertPositions[b].Y
		for j := i + 1; j < n; j++ {
			// skip if edges share a vertex
			if edges[i].From == edges[j].From ||
				edges[i].From == edges[j].To ||
				edges[i].To == edges[j].From ||
				edges[i].To == edges[j].To {
				continue
			}
			c := edges[j].From
			d := edges[j].To
			cx, cy := s.VertPositions[c].X, s.VertPositions[c].Y
			dx, dy := s.VertPositions[d].X, s.VertPositions[d].Y
			if segmentsIntersect(ax, ay, bx, by, cx, cy, dx, dy) {
				count++
			}
		}
	}

	// Calculate dispersion penalty based on minimal distance between vertices.
	// Desired minimum distance is proportional to the area coverage.
	desiredMin := math.Min(s.width, s.height) / math.Sqrt(float64(s.graph.NumVertices))
	minDist := math.MaxFloat64
	for i := 0; i < len(s.VertPositions); i++ {
		for j := i + 1; j < len(s.VertPositions); j++ {
			dx := s.VertPositions[i].X - s.VertPositions[j].X
			dy := s.VertPositions[i].Y - s.VertPositions[j].Y
			d := math.Sqrt(dx*dx + dy*dy)
			if d < minDist {
				minDist = d
			}
		}
	}
	dispersionPenalty := 0.0
	if minDist < desiredMin {
		dispersionPenalty = (desiredMin - minDist) / desiredMin
	}

	// Calculate collinearity penalty.
	// Compute mean of vertices.
	meanX, meanY := 0.0, 0.0
	nv := float64(len(s.VertPositions))
	for _, v := range s.VertPositions {
		meanX += v.X
		meanY += v.Y
	}
	meanX /= nv
	meanY /= nv
	// Compute covariance components.
	sxx, syy, sxy := 0.0, 0.0, 0.0
	for _, v := range s.VertPositions {
		dx := v.X - meanX
		dy := v.Y - meanY
		sxx += dx * dx
		syy += dy * dy
		sxy += dx * dy
	}
	sxx /= nv
	syy /= nv
	sxy /= nv
	// Calculate eigenvalues of the covariance matrix.
	trace := sxx + syy
	det := sxx*syy - sxy*sxy
	eigenVal2 := (trace - math.Sqrt(trace*trace-4*det)) / 2
	eigenVal1 := trace - eigenVal2
	ratio := 1.0
	if eigenVal1 > 0 {
		ratio = eigenVal2 / eigenVal1
	}
	collinearityPenalty := 0.0
	threshold := 0.1
	if ratio < threshold {
		collinearityPenalty = (threshold - ratio) / threshold
	}

	// Higher fitness for fewer intersections, good dispersion, and non-collinearity.
	return 1.0 / (1.0 + float64(count) + (dispersionPenalty+collinearityPenalty)*10)
}
