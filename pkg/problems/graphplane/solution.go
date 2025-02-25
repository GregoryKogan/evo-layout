package graphplane

import (
	"math"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type GraphPlaneSolution struct {
	graph         *Graph
	width, height float64
	cachedFitness float64
	VertPositions []VertexPos `json:"vertices"`
	Intersections int         `json:"intersections"`
	Dispersion    float64     `json:"dispersion"`
}

type VertexPos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func newSolutionBase(old *GraphPlaneSolution) *GraphPlaneSolution {
	return &GraphPlaneSolution{
		graph:         old.graph,
		width:         old.width,
		height:        old.height,
		VertPositions: make([]VertexPos, old.graph.NumVertices),
	}
}

func RandomGraphPlaneSolution(g *Graph, width, height float64) problems.Solution {
	s := newSolutionBase(&GraphPlaneSolution{graph: g, width: width, height: height})
	for i := range s.VertPositions {
		s.VertPositions[i] = VertexPos{
			X: rand.Float64() * width,
			Y: rand.Float64() * height,
		}
	}
	return s
}

func (s *GraphPlaneSolution) Crossover(other problems.Solution) problems.Solution {
	otherGPS, ok := other.(*GraphPlaneSolution)
	if !ok {
		return s
	}

	child := newSolutionBase(s)
	for i := range child.VertPositions {
		if rand.Float64() < 0.5 {
			child.VertPositions[i] = s.VertPositions[i]
		} else {
			child.VertPositions[i] = otherGPS.VertPositions[i]
		}
	}
	return child
}

func (s *GraphPlaneSolution) Mutate(rate float64) problems.Solution {
	mutant := newSolutionBase(s)
	copy(mutant.VertPositions, s.VertPositions)

	intersectionWeights := s.calculateIntersectionWeights()
	selectedVertex := weightedRandomSelect(intersectionWeights)
	// TODO: only one vertex is modified. Should be able to change multiple depending on mutation rate
	deltaX, deltaY := calculateMutationDeltas(rate, s.width, s.height)
	mutant.applyMutation(selectedVertex, deltaX, deltaY)

	return mutant
}

func (s *GraphPlaneSolution) Fitness() float64 {
	if s.cachedFitness > 0 {
		return s.cachedFitness
	}
	s.Intersections = s.countIntersections()
	s.Dispersion = s.calculateDispersionPenalty()
	s.cachedFitness = 1.0 / (float64(s.Intersections) + float64(s.Intersections+1)*s.Dispersion)
	return s.cachedFitness
}

func (s *GraphPlaneSolution) calculateIntersectionWeights() []float64 {
	weights := make([]float64, s.graph.NumVertices)
	s.forEachEdgePair(func(e1, e2 Edge, a, b, c, d VertexPos) {
		if segmentsIntersect(a, b, c, d) {
			weights[e1.From] += 1.0
			weights[e1.To] += 1.0
			weights[e2.From] += 1.0
			weights[e2.To] += 1.0
		}
	})
	return weights
}

func calculateMutationDeltas(rate, width, height float64) (float64, float64) {
	return rand.NormFloat64() * width * rate,
		rand.NormFloat64() * height * rate
}

func (s *GraphPlaneSolution) applyMutation(vertex int, deltaX, deltaY float64) {
	s.VertPositions[vertex].X = clamp(s.VertPositions[vertex].X+deltaX, 0, s.width)
	s.VertPositions[vertex].Y = clamp(s.VertPositions[vertex].Y+deltaY, 0, s.height)
}

// Common edge pair iteration logic
func (s *GraphPlaneSolution) forEachEdgePair(fn func(e1, e2 Edge, a, b, c, d VertexPos)) {
	edges := s.graph.Edges
	for i := 0; i < len(edges); i++ {
		e1 := edges[i]
		a := s.VertPositions[e1.From]
		b := s.VertPositions[e1.To]
		for j := i + 1; j < len(edges); j++ {
			e2 := edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			c := s.VertPositions[e2.From]
			d := s.VertPositions[e2.To]
			fn(e1, e2, a, b, c, d)
		}
	}
}

func (s *GraphPlaneSolution) countIntersections() int {
	count := 0
	s.forEachEdgePair(func(e1, e2 Edge, a, b, c, d VertexPos) {
		if segmentsIntersect(a, b, c, d) {
			count++
		}
	})
	return count
}

func weightedRandomSelect(weights []float64) int {
	// Calculate total weight
	total := 0.0
	for _, w := range weights {
		total += w
	}

	// If no intersections, select uniformly
	if total == 0 {
		return rand.Intn(len(weights))
	}

	// Generate random position
	r := rand.Float64() * total
	for i, w := range weights {
		r -= w
		if r < 0 {
			return i
		}
	}

	// Fallback to random selection
	return rand.Intn(len(weights))
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func sharesVertex(e1, e2 Edge) bool {
	return e1.From == e2.From || e1.From == e2.To || e1.To == e2.From || e1.To == e2.To
}

func segmentsIntersect(a, b, c, d VertexPos) bool {
	ccw := func(ax, ay, bx, by, cx, cy float64) bool {
		return (cy-ay)*(bx-ax) > (by-ay)*(cx-ax)
	}
	return ccw(a.X, a.Y, c.X, c.Y, d.X, d.Y) != ccw(b.X, b.Y, c.X, c.Y, d.X, d.Y) &&
		ccw(a.X, a.Y, b.X, b.Y, c.X, c.Y) != ccw(a.X, a.Y, b.X, b.Y, d.X, d.Y)
}

func (s *GraphPlaneSolution) calculateDispersionPenalty() float64 {
	desiredMin := math.Min(s.width, s.height) / math.Sqrt(float64(s.graph.NumVertices))
	minDist := s.computeMinVertexDistance()
	if minDist <= 0.001 {
		return math.MaxFloat64
	}
	if minDist < desiredMin {
		return (desiredMin - minDist) / desiredMin
	}
	return 0.0
}

func (s *GraphPlaneSolution) computeMinVertexDistance() float64 {
	minDist := math.MaxFloat64
	positions := s.VertPositions
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			dx := positions[i].X - positions[j].X
			dy := positions[i].Y - positions[j].Y
			dist := math.Hypot(dx, dy)
			if dist < minDist {
				minDist = dist
			}
		}
	}
	return minDist
}
