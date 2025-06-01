package gjkepa2d //nolint:mnd

import (
	"math"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

// GJKEPA check related constants
const (
	NoIntersection = iota
	FoundIntersection
	StillEvolving

	Clockwise
	CounterClockwise

	GJKMaxNumIterations = 64
	EPAMaxNumIterations = 64
)

// Edge describe simplex edge
type Edge struct {
	distance float64
	normal   shapes.Point
	index    int
}

// NewEdge return new edge
func NewEdge(distance float64, normal shapes.Point, index int) *Edge {
	return &Edge{
		distance: distance,
		normal:   normal,
		index:    index,
	}
}

// GJKEPA gjk epa collision checker
type GJKEPA struct {
	vertices  []shapes.Point
	direction shapes.Point
	shapeA    shapes.Collide
	shapeB    shapes.Collide
}

// NewGJKEPA return new 2d collision checker
func NewGJKEPA() *GJKEPA {
	return &GJKEPA{}
}

// calculateSupport calculate support point
func (g *GJKEPA) calculateSupport(direction shapes.Point) shapes.Point {
	return g.shapeA.Support(direction).Subtract(g.shapeB.Support(direction.Invert()))
}

// addSupport add support point
func (g *GJKEPA) addSupport(direction shapes.Point) bool {
	newVertex := g.calculateSupport(direction)
	g.vertices = append(g.vertices, newVertex)

	return direction.Dot(newVertex) >= 0
}

// tripleProduct calculate triple
func (g *GJKEPA) tripleProduct(pa, pb, pc shapes.Point) shapes.Point {
	pa = shapes.NewPoint(pa.Coordinate(0), pa.Coordinate(1))
	pb = shapes.NewPoint(pb.Coordinate(0), pb.Coordinate(1))
	pc = shapes.NewPoint(pc.Coordinate(0), pc.Coordinate(1))

	first := pa.Cross(pb)
	second := first.Cross(pc)

	return shapes.NewPoint(second.Coordinate(0), second.Coordinate(1))
}

// remove remove vector from vertices
func (g *GJKEPA) remove(v shapes.Point) {
	for i, e := range g.vertices {
		if e.Equals(v) {
			g.vertices = append(g.vertices[:i], g.vertices[i+1:]...)
		}
	}
}

// evolveSimplex evolve simplex
func (g *GJKEPA) evolveSimplex() int {
	switch len(g.vertices) {
	case 0: //nolint:mnd
		g.direction = g.shapeB.Point1().Subtract(g.shapeA.Point1())
	case 1:
		g.direction = g.direction.Scale(-1)
	case 2: //nolint:mnd
		ab := g.vertices[1].Subtract(g.vertices[0])
		a0 := g.vertices[0].Scale(-1)
		g.direction = g.tripleProduct(ab, a0, ab)
	case 3: //nolint:mnd
		c0 := g.vertices[2].Scale(-1)
		bc := g.vertices[1].Subtract(g.vertices[2])
		ca := g.vertices[0].Subtract(g.vertices[2])

		bcNorm := g.tripleProduct(ca, bc, bc)
		caNorm := g.tripleProduct(bc, ca, ca)

		if bcNorm.Dot(c0) > 0 { //nolint:gocritic
			g.remove(g.vertices[0])
			g.direction = bcNorm
		} else if caNorm.Dot(c0) > 0 {
			g.remove(g.vertices[1])
			g.direction = caNorm
		} else {
			return FoundIntersection
		}
	}

	if g.addSupport(g.direction) {
		return StillEvolving
	}

	return NoIntersection
}

// GJK check GJK collision
func (g *GJKEPA) GJK(shapeA, shapeB shapes.Collide, calculateMTV bool) (bool, shapes.Point) {
	g.vertices = []shapes.Point{}
	g.shapeA = shapeA
	g.shapeB = shapeB

	result := StillEvolving
	for iteration := 0; iteration < GJKMaxNumIterations && result == StillEvolving; iteration++ {
		result = g.evolveSimplex()
	}

	i := result == FoundIntersection

	var p shapes.Point
	if i && calculateMTV && len(g.vertices) == 3 {
		p = g.EPA(shapeA, shapeB)
	}

	return i, p
}

// findClosestEdge calculate closest edge
func (g *GJKEPA) findClosestEdge(winding int) *Edge {
	var closestIndex int
	var closestNormal shapes.Point
	closestDistance := math.MaxFloat64

	for i := range g.vertices {
		j := i + 1
		if j >= len(g.vertices) {
			j = 0
		}

		line := g.vertices[j].Subtract(g.vertices[i])
		var norm shapes.Point

		switch winding {
		case Clockwise:
			norm = shapes.NewPoint(line.Coordinate(1), line.Coordinate(0)*-1)
		case CounterClockwise:
			norm = shapes.NewPoint(line.Coordinate(1)*-1, line.Coordinate(0))
		default:
			return nil
		}

		norm = norm.Normalize()
		dist := norm.Dot(g.vertices[i])

		if dist < closestDistance {
			closestDistance = dist
			closestNormal = norm
			closestIndex = j
		}
	}

	return NewEdge(closestDistance, closestNormal, closestIndex)
}

// EPA calculate EPA intersection point
func (g *GJKEPA) EPA(_, _ shapes.Collide) shapes.Point {
	// calculate the winding of the existing simplex
	e0 := (g.vertices[1].Coordinate(0) - g.vertices[0].Coordinate(0)) *
		(g.vertices[1].Coordinate(1) + g.vertices[0].Coordinate(1))
	e1 := (g.vertices[2].Coordinate(0) - g.vertices[1].Coordinate(0)) *
		(g.vertices[2].Coordinate(1) + g.vertices[1].Coordinate(1))
	e2 := (g.vertices[0].Coordinate(0) - g.vertices[2].Coordinate(0)) *
		(g.vertices[0].Coordinate(1) + g.vertices[2].Coordinate(1))

	var winding int
	if e0+e1+e2 >= 0 {
		winding = Clockwise
	} else {
		winding = CounterClockwise
	}

	var minIntersection shapes.Point
	minDistance := math.MaxFloat64
	var intersection shapes.Point

	for i := 0; i < EPAMaxNumIterations; i++ {
		edge := g.findClosestEdge(winding)
		support := g.calculateSupport(edge.normal)
		distance := support.Dot(edge.normal)

		intersection = edge.normal.Copy()
		intersection = intersection.Scale(distance).Invert()

		if distance < minDistance {
			minDistance = distance
			minIntersection = intersection
		}

		if math.Abs(distance-edge.distance) <= 0.000001 { //nolint:mnd
			return intersection
		}

		g.vertices[edge.index] = support
	}

	return minIntersection
}
