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
	edges     []edgeData
	direction shapes.Point
	shapeA    shapes.Collide
	shapeB    shapes.Collide
}

// NewGJKEPA return new 2d collision checker
func NewGJKEPA() *GJKEPA {
	return &GJKEPA{
		// GJK needs at most 3 vertices; EPA inserts at most one per iteration.
		// Preallocating lets GJK/EPA run without growing the buffers.
		vertices: make([]shapes.Point, 0, 3+EPAMaxNumIterations),
		edges:    make([]edgeData, 0, 3+EPAMaxNumIterations),
	}
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
	g.vertices = g.vertices[:0]
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

// edgeData caches a polytope edge's outward normal and its distance from
// the origin, so EPA recomputes only the two edges an insertion creates
// instead of every edge on every iteration.
type edgeData struct {
	dist   float64
	nx, ny float64
}

// computeEdge calculates the polytope edge starting at vertex i. The scalar
// math mirrors the shapes.Point operations it replaces (Subtract, Normalize
// with its zero-length guard, Dot with a zero z-term) so results are
// bit-identical; it just avoids allocating Points in EPA's hottest loop.
func (g *GJKEPA) computeEdge(winding, i int) edgeData {
	j := i + 1
	if j >= len(g.vertices) {
		j = 0
	}

	lx := g.vertices[j].Coordinate(0) - g.vertices[i].Coordinate(0)
	ly := g.vertices[j].Coordinate(1) - g.vertices[i].Coordinate(1)

	var nx, ny float64

	switch winding {
	case Clockwise:
		nx, ny = ly, lx*-1
	case CounterClockwise:
		nx, ny = ly*-1, lx
	default:
		return edgeData{} // unreachable: EPA always passes a valid winding
	}

	// Point.Normalize: length via sqrt(dot), zero guard, then divide
	n := math.Sqrt(nx*nx + ny*ny)
	if n < math.SmallestNonzeroFloat64 {
		nx, ny = 0, 0
	} else {
		nx /= n
		ny /= n
	}

	dist := nx*g.vertices[i].Coordinate(0) + ny*g.vertices[i].Coordinate(1)

	return edgeData{dist: dist, nx: nx, ny: ny}
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

	g.edges = g.edges[:0]
	for i := range g.vertices {
		g.edges = append(g.edges, g.computeEdge(winding, i))
	}

	minDistance := math.MaxFloat64

	var minIX, minIY float64

	for i := 0; i < EPAMaxNumIterations; i++ {
		// closest cached edge; seeding with MaxFloat64 keeps the earliest
		// strict minimum and skips NaN/+Inf distances exactly like the full
		// rescan this replaces
		var edge edgeData

		edge.dist = math.MaxFloat64
		closest := -1

		for e := range g.edges {
			if g.edges[e].dist < edge.dist {
				edge = g.edges[e]
				closest = e
			}
		}

		normal := shapes.NewPoint(edge.nx, edge.ny)
		support := g.calculateSupport(normal)
		distance := support.Dot(normal)

		ix := edge.nx * distance * -1
		iy := edge.ny * distance * -1

		if distance < minDistance {
			minDistance = distance
			minIX, minIY = ix, iy
		}

		if math.Abs(distance-edge.dist) <= 0.000001 { //nolint:mnd
			return shapes.NewPoint(ix, iy)
		}

		// expand the polytope: insert the support point between the closest
		// edge's endpoints (index k, its second endpoint), keeping hull
		// order; shifting in place avoids reallocating every iteration.
		// closest == -1 (every distance NaN) mirrors the original's
		// zero-value index 0
		k := 0
		if closest >= 0 && closest+1 < len(g.vertices) {
			k = closest + 1
		}

		g.vertices = append(g.vertices, shapes.Point{})
		copy(g.vertices[k+1:], g.vertices[k:])
		g.vertices[k] = support

		// the insertion replaces one edge with two; every other edge keeps
		// its vertex pair, so only those two need recomputing
		g.edges = append(g.edges, edgeData{})
		copy(g.edges[k+1:], g.edges[k:])

		prev := k - 1
		if prev < 0 {
			prev = len(g.edges) - 1
		}

		g.edges[prev] = g.computeEdge(winding, prev)
		g.edges[k] = g.computeEdge(winding, k)
	}

	return shapes.NewPoint(minIX, minIY)
}
