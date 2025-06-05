package graphplane

import (
	"math"
	"math/rand"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"gonum.org/v1/gonum/stat"
)

// VertexPos is the (x,y) coordinate of a vertex.
type VertexPos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// GraphPlaneSolution represents a placement of Graph vertices in the plane.
type GraphPlaneSolution struct {
	Graph            *Graph      `json:"-"`
	Width, Height    float64     `json:"-"`
	Intersections    int         `json:"intersections"`
	VertPositions    []VertexPos `json:"vertices"`
	CachedObjectives []float64   `json:"objectives"`
	CachedFitness    float64     `json:"fitness"`
}

// RandomGraphPlaneSolution initializes vertices randomly in [0,width]×[0,height].
func RandomGraphPlaneSolution(g *Graph, width, height float64) problems.Solution {
	s := &GraphPlaneSolution{Graph: g, Width: width, Height: height}
	s.VertPositions = make([]VertexPos, g.NumVertices)
	for i := range s.VertPositions {
		s.VertPositions[i] = VertexPos{X: rand.Float64() * width, Y: rand.Float64() * height}
	}
	return s
}

// Objectives returns:
//
//	[0] intersections,
//	[1] dispersion penalty.
//	[2] min-angle penalty,
func (s *GraphPlaneSolution) Objectives() []float64 {
	if len(s.CachedObjectives) > 0 {
		return s.CachedObjectives
	}

	intersections := float64(s.CountIntersections())
	dispersion := s.dispersionPenalty()
	angle := s.anglePenalty()
	dispersion *= (intersections + 1.0)
	angle *= (intersections + 1.0)
	s.CachedObjectives = []float64{intersections, dispersion, angle}
	return s.CachedObjectives
}

// Fitness for single-objective algorithms
func (s *GraphPlaneSolution) Fitness() float64 {
	o := s.Objectives()
	s.CachedFitness = o[0] + o[1] + o[2]
	return s.CachedFitness
}

// anglePenalty penalizes uneven angles between edges
func (s *GraphPlaneSolution) anglePenalty() float64 {
	angles := make([]float64, 0, s.Graph.NumEdges*s.Graph.NumEdges)

	for v := range s.VertPositions {
		// collect incident neighbors
		var neigh []VertexPos
		for _, e := range s.Graph.Edges {
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
				angles = append(angles, ing)
			}
		}
	}

	return stat.StdDev(angles, nil)
}

// dispersionPenalty penalizes uneven distances between vertexes
func (s *GraphPlaneSolution) dispersionPenalty() float64 {
	n := s.Graph.NumVertices

	distances := make([]float64, 0, n*(n-1))
	for i := range n {
		for j := range n {
			if i == j {
				continue
			}
			d := math.Hypot(
				s.VertPositions[i].X-s.VertPositions[j].X,
				s.VertPositions[i].Y-s.VertPositions[j].Y,
			)
			distances = append(distances, d)
		}
	}

	mean, std := stat.MeanStdDev(distances, nil)
	desired := math.Min(s.Width, s.Height) / math.Sqrt(float64(n)) * 2
	return math.Abs(desired-mean)/desired + std
}

// CountIntersections counts all pairwise edge crossings.
func (s *GraphPlaneSolution) CountIntersections() int {
	cnt := 0
	for i := range s.Graph.Edges {
		e1 := s.Graph.Edges[i]
		p1, p2 := s.VertPositions[e1.From], s.VertPositions[e1.To]
		for j := i + 1; j < len(s.Graph.Edges); j++ {
			e2 := s.Graph.Edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			p3, p4 := s.VertPositions[e2.From], s.VertPositions[e2.To]
			if segmentsIntersect(p1, p2, p3, p4) {
				cnt++
			}
		}
	}
	s.Intersections = cnt
	return cnt
}

func (s *GraphPlaneSolution) TangledVertexes() []int {
	isTangled := make(map[int]bool)
	for i := range s.Graph.Edges {
		e1 := s.Graph.Edges[i]
		p1, p2 := s.VertPositions[e1.From], s.VertPositions[e1.To]
		for j := i + 1; j < len(s.Graph.Edges); j++ {
			e2 := s.Graph.Edges[j]
			if sharesVertex(e1, e2) {
				continue
			}
			p3, p4 := s.VertPositions[e2.From], s.VertPositions[e2.To]
			if segmentsIntersect(p1, p2, p3, p4) {
				isTangled[e1.From] = true
				isTangled[e1.To] = true
				isTangled[e2.From] = true
				isTangled[e2.To] = true
			}
		}
	}

	vertexes := make([]int, 0, len(isTangled))
	for v := range isTangled {
		vertexes = append(vertexes, v)
	}

	return vertexes
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
