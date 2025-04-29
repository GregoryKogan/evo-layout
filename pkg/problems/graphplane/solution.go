package graphplane

import (
	"math"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

// VertexPos is the (x,y) coordinate of a vertex.
type VertexPos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// GraphPlaneSolution represents a placement of graph vertices in the plane.
type GraphPlaneSolution struct {
	graph            *Graph
	width, height    float64
	VertPositions    []VertexPos `json:"vertices"`
	CachedObjectives []float64   `json:"objectives"`
}

// RandomGraphPlaneSolution initializes vertices randomly in [0,width]×[0,height].
func RandomGraphPlaneSolution(g *Graph, width, height float64) problems.Solution {
	s := &GraphPlaneSolution{graph: g, width: width, height: height}
	s.VertPositions = make([]VertexPos, g.NumVertices)
	for i := range s.VertPositions {
		s.VertPositions[i] = VertexPos{X: rand.Float64() * width, Y: rand.Float64() * height}
	}
	return s
}

// Crossover combines two parent solutions via uniform position mixing.
func (s *GraphPlaneSolution) Crossover(other problems.Solution) []problems.Solution {
	rhs, ok := other.(*GraphPlaneSolution)
	if !ok || len(s.VertPositions) != len(rhs.VertPositions) {
		return []problems.Solution{s}
	}
	// Create two children
	c1 := &GraphPlaneSolution{graph: s.graph, width: s.width, height: s.height}
	c2 := &GraphPlaneSolution{graph: s.graph, width: s.width, height: s.height}
	c1.VertPositions = make([]VertexPos, len(s.VertPositions))
	c2.VertPositions = make([]VertexPos, len(s.VertPositions))
	for i := range s.VertPositions {
		if rand.Float64() < 0.5 {
			c1.VertPositions[i] = s.VertPositions[i]
			c2.VertPositions[i] = rhs.VertPositions[i]
		} else {
			c1.VertPositions[i] = rhs.VertPositions[i]
			c2.VertPositions[i] = s.VertPositions[i]
		}
	}
	return []problems.Solution{c1, c2}
}

// Mutate perturbs vertices, focusing on those with many crossings.
func (s *GraphPlaneSolution) Mutate(rate float64) problems.Solution {
	m := &GraphPlaneSolution{graph: s.graph, width: s.width, height: s.height}
	m.VertPositions = make([]VertexPos, len(s.VertPositions))
	copy(m.VertPositions, s.VertPositions)
	// Intersection-based weights
	weights := s.intersectionWeights()
	for i := range m.VertPositions {
		// base probability plus weighted factor
		p := rate * (weights[i]/(weights[i]+1) + 0.1)
		if rand.Float64() < p {
			dx := rand.NormFloat64() * s.width * rate
			dy := rand.NormFloat64() * s.height * rate
			m.VertPositions[i].X = clamp(m.VertPositions[i].X+dx, 0, s.width)
			m.VertPositions[i].Y = clamp(m.VertPositions[i].Y+dy, 0, s.height)
		}
	}
	return m
}

// Objectives returns:
//
//	[0] intersections (min),
//	[1] smallest-angle penalty (min),
//	[2] normalized avg edge length (min),
//	[3] dispersion penalty (min).
func (s *GraphPlaneSolution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}
	inter := float64(s.countIntersections())
	anglePen := s.smallestAnglePenalty()
	avgLen := s.avgEdgeLength() / math.Hypot(s.width, s.height)
	disp := s.dispersionPenalty()
	s.CachedObjectives = []float64{inter, anglePen, avgLen, disp}
	return s.CachedObjectives
}

// Fitness for single-objective algorithms (fallback).
func (s *GraphPlaneSolution) Fitness() float64 {
	o := s.Objectives()
	sum := 0.0
	for _, v := range o {
		sum += v
	}
	return sum
}

// avgEdgeLength computes mean length of all edges.
func (s *GraphPlaneSolution) avgEdgeLength() float64 {
	tot := 0.0
	for _, e := range s.graph.Edges {
		p, q := s.VertPositions[e.From], s.VertPositions[e.To]
		tot += math.Hypot(p.X-q.X, p.Y-q.Y)
	}
	n := float64(len(s.graph.Edges))
	if n == 0 {
		return 0
	}
	return tot / n
}

// smallestAnglePenalty computes penalty = (π - minAngle)/π, based on
// the smallest angle between any two edges sharing a vertex.
func (s *GraphPlaneSolution) smallestAnglePenalty() float64 {
	minAngle := math.Pi
	for v := range s.VertPositions {
		// collect incident neighbors
		var neigh []VertexPos
		for _, e := range s.graph.Edges {
			if e.From == v {
				neigh = append(neigh, s.VertPositions[e.To])
			} else if e.To == v {
				neigh = append(neigh, s.VertPositions[e.From])
			}
		}
		// compute angles between each pair
		for i := range neigh {
			for j := i + 1; j < len(neigh); j++ {
				u := VertexPos{X: neigh[i].X - s.VertPositions[v].X, Y: neigh[i].Y - s.VertPositions[v].Y}
				w := VertexPos{X: neigh[j].X - s.VertPositions[v].X, Y: neigh[j].Y - s.VertPositions[v].Y}
				dot := u.X*w.X + u.Y*w.Y
				du := math.Hypot(u.X, u.Y)
				dv := math.Hypot(w.X, w.Y)
				if du == 0 || dv == 0 {
					continue
				}
				ing := math.Acos(math.Min(1, math.Max(-1, dot/(du*dv))))
				if ing < minAngle {
					minAngle = ing
				}
			}
		}
	}
	if minAngle == math.Pi {
		return 0
	}
	return (math.Pi - minAngle) / math.Pi
}

// dispersionPenalty penalizes too-close vertices.
func (s *GraphPlaneSolution) dispersionPenalty() float64 {
	n := float64(len(s.VertPositions))
	desired := math.Min(s.width, s.height) / math.Sqrt(n)
	minD := math.MaxFloat64
	for i := range s.VertPositions {
		for j := i + 1; j < len(s.VertPositions); j++ {
			d := math.Hypot(
				s.VertPositions[i].X-s.VertPositions[j].X,
				s.VertPositions[i].Y-s.VertPositions[j].Y)
			if d < minD {
				minD = d
			}
		}
	}
	if minD >= desired {
		return 0
	}
	return (desired - minD) / desired
}

// countIntersections counts all pairwise edge crossings.
func (s *GraphPlaneSolution) countIntersections() int {
	cnt := 0
	for i := range s.graph.Edges {
		e1 := s.graph.Edges[i]
		p1, p2 := s.VertPositions[e1.From], s.VertPositions[e1.To]
		for j := i + 1; j < len(s.graph.Edges); j++ {
			e2 := s.graph.Edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			p3, p4 := s.VertPositions[e2.From], s.VertPositions[e2.To]
			if segmentsIntersect(p1, p2, p3, p4) {
				cnt++
			}
		}
	}
	return cnt
}

// intersectionWeights returns per-vertex count of incident crossings.
func (s *GraphPlaneSolution) intersectionWeights() []float64 {
	w := make([]float64, len(s.VertPositions))
	for i := range s.graph.Edges {
		e1 := s.graph.Edges[i]
		p1, p2 := s.VertPositions[e1.From], s.VertPositions[e1.To]
		for j := i + 1; j < len(s.graph.Edges); j++ {
			e2 := s.graph.Edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			p3, p4 := s.VertPositions[e2.From], s.VertPositions[e2.To]
			if segmentsIntersect(p1, p2, p3, p4) {
				w[e1.From]++
				w[e1.To]++
				w[e2.From]++
				w[e2.To]++
			}
		}
	}
	return w
}

// sharesVertex checks if two edges share an endpoint.
func sharesVertex(e1, e2 Edge) bool {
	return e1.From == e2.From || e1.From == e2.To || e1.To == e2.From || e1.To == e2.To
}

// segmentsIntersect tests segment intersection via CCW.
func segmentsIntersect(a, b, c, d VertexPos) bool {
	ccw := func(u, v, w VertexPos) bool {
		return (w.Y-u.Y)*(v.X-u.X) > (v.Y-u.Y)*(w.X-u.X)
	}
	return ccw(a, c, d) != ccw(b, c, d) && ccw(a, b, c) != ccw(a, b, d)
}

// clamp bounds v to [min,max].
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
