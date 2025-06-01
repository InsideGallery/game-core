package shapes

import (
	"math"
)

// Box describe rectangle
type Box struct {
	Point
	sizes [3]float64
}

// NewBox return new box
func NewBox(p Point, s ...float64) Box {
	var width, height, length float64
	l := len(s)

	switch l {
	case 3: // nolint:mnd
		width, height, length = s[0], s[1], s[2]
	case 2: // nolint:mnd
		width, height = s[0], s[1]
	case 1: // nolint:mnd
		width = s[0]
	}

	return Box{
		Point: p,
		sizes: [3]float64{width, height, length},
	}
}

// Point1 return first point
func (b Box) Point1() Point {
	return b.Point
}

// Point2 return second point of rectangle
func (b Box) Point2() Point {
	p2 := make([]float64, 3) //nolint:mnd
	for k, v := range b.Point1().Coordinates() {
		p2[k] = v + b.Size(k)
	}

	return NewPoint(p2...)
}

// Get return Spatial
func (b Box) Get() Spatial {
	return b
}

// Size return size by dimension
func (b Box) Size(i int) float64 {
	if i > 2 || i < 0 {
		return 0
	}

	return b.sizes[i]
}

// Sizes return all sizes
func (b Box) Sizes() [3]float64 {
	return b.sizes
}

// VectorSizes return all sizes as vector
func (b Box) VectorSizes() Point {
	return NewPoint(b.sizes[0], b.sizes[1], b.sizes[2])
}

// Bounds return rectangle of object
func (b Box) Bounds() Box {
	return b
}

// Move return new object for given diff
func (b Box) Move(diff Point) Spatial {
	d := make([]float64, 3) //nolint:mnd
	s := make([]float64, 3) //nolint:mnd

	for k, v := range b.Coordinates() {
		d[k] = v + diff.Coordinate(k)
		s[k] = b.Size(k)
	}

	return NewBox(NewPoint(d...), s...)
}

// Equal returns true if the two rectangles are equal
func (b Box) Equal(r2 Box) bool {
	for i, e := range b.Point1().Coordinates() {
		if e != r2.Point1().Coordinate(i) {
			return false
		}
	}

	for i, e := range b.Point2().Coordinates() {
		if e != r2.Point2().Coordinate(i) {
			return false
		}
	}

	return true
}

// Volume return size of rectangle
func (b Box) Volume() float64 {
	size := 1.0
	for _, a := range b.Sizes() {
		size *= a
	}

	return size
}

// Margin computes the sum of the edge lengths of a rectangle.
func (b Box) Margin() float64 {
	// The number of edges in an n-dimensional rectangle is n * 2^(n-1),
	// for our case it`s always 2 (http://en.wikipedia.org/wiki/Hypercube_graph). 2^(3-1)=2
	sum := 0.0
	for _, a := range b.Sizes() {
		sum += a
	}

	return 2 * sum // nolint:mnd
}

// ContainsPoint tests whether p is located inside or on the boundary of r
func (b Box) ContainsPoint(p Point) bool {
	for i, a := range p.Coordinates() {
		// p is contained in (or on) b if and only if p <= a <= q for
		// every dimension.
		if a < b.Point1().Coordinate(i) || a > b.Point2().Coordinate(i) {
			return false
		}
	}

	return true
}

// ContainsRectangle tests whether r2 is is located inside r1.
func (b Box) ContainsRectangle(r2 Box) bool {
	for i, a1 := range b.Point1().Coordinates() {
		b1, a2, b2 := b.Point2().Coordinate(i), r2.Point1().Coordinate(i), r2.Point2().Coordinate(i)
		// enforced by constructor: a1 <= b1 and a2 <= b2.
		// so containment holds if and only if a1 <= a2 <= b2 <= b1
		// for every dimension.
		if a1 > a2 || b2 > b1 {
			return false
		}
	}

	return true
}

// Intersect computes the intersection of two rectangles
func (b Box) Intersect(r2 Box) (Box, bool) {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// Thus there are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2
	//
	// Enforced by constructor: a1 <= b1 and a2 <= b2.  So we can just
	// check the endpoints.
	p := make([]float64, 3) //nolint:mnd
	s := make([]float64, 3) //nolint:mnd

	for i := range p {
		a1, b1, a2, b2 := b.Point1().Coordinate(i), b.Point2().Coordinate(i), r2.Point1().Coordinate(i), r2.Point2().Coordinate(i) //nolint:lll
		if a1 > b2 || a2 > b1 {
			return NewBox(NewPoint()), false
		}

		p[i] = math.Max(a1, a2)
		s[i] = math.Min(b1, b2) - p[i]
	}

	return NewBox(NewPoint(p...), s...), true
}

// BoundingBox constructs the smallest rectangle containing both r1 and r2.
func (b Box) BoundingBox(r2 Box) Box {
	p := make([]float64, 3) //nolint:mnd
	s := make([]float64, 3) //nolint:mnd

	for i := 0; i < 3; i++ {
		if b.Point1().Coordinate(i) <= r2.Point1().Coordinate(i) {
			p[i] = b.Point1().Coordinate(i)
		} else {
			p[i] = r2.Point1().Coordinate(i)
		}

		if b.Point2().Coordinate(i) <= r2.Point2().Coordinate(i) {
			s[i] = r2.Point2().Coordinate(i) - p[i]
		} else {
			s[i] = b.Point2().Coordinate(i) - p[i]
		}
	}

	return NewBox(NewPoint(p...), s...)
}

// Center return center point
func (b Box) Center() Point {
	p := make([]float64, 3) // nolint:mnd
	for k, v := range b.Coordinates() {
		p[k] = v + (b.Size(k) / 2) // nolint:mnd
	}

	return NewPoint(p...)
}

// ToPolyhedron return polyhedron from box
func (b Box) ToPolyhedron() Polyhedron {
	p1, p2 := b.Point1(), b.Point2()
	v := make([]Point, 8)                                                 //nolint:mnd
	v[0] = NewPoint(p1.Coordinate(0), p1.Coordinate(1), p1.Coordinate(2)) //nolint:mnd
	v[1] = NewPoint(p1.Coordinate(0), p2.Coordinate(1), p1.Coordinate(2)) //nolint:mnd
	v[2] = NewPoint(p2.Coordinate(0), p2.Coordinate(1), p1.Coordinate(2)) //nolint:mnd
	v[3] = NewPoint(p2.Coordinate(0), p1.Coordinate(1), p1.Coordinate(2)) //nolint:mnd

	v[4] = NewPoint(p1.Coordinate(0), p1.Coordinate(1), p2.Coordinate(2)) //nolint:mnd
	v[5] = NewPoint(p1.Coordinate(0), p2.Coordinate(1), p2.Coordinate(2)) //nolint:mnd
	v[6] = NewPoint(p2.Coordinate(0), p2.Coordinate(1), p2.Coordinate(2)) //nolint:mnd
	v[7] = NewPoint(p2.Coordinate(0), p1.Coordinate(1), p2.Coordinate(2)) //nolint:mnd

	return NewPolyhedron(v...)
}

// ToPolygon return polygon from box
func (b Box) ToPolygon() Polyhedron {
	p1, p2 := b.Point1(), b.Point2()
	v := make([]Point, 4) //nolint:mnd
	v[0] = NewPoint(p1.Coordinate(0), p1.Coordinate(1))
	v[1] = NewPoint(p1.Coordinate(0), p2.Coordinate(1))
	v[2] = NewPoint(p2.Coordinate(0), p2.Coordinate(1))
	v[3] = NewPoint(p2.Coordinate(0), p1.Coordinate(1))

	return NewPolyhedron(v...)
}

// Support return support point for box
func (b Box) Support(d Point) Point {
	var x, y, z float64

	if d.Coordinate(0) > 0 {
		x = b.Point2().Coordinate(0)
	} else {
		x = b.Point1().Coordinate(0)
	}

	if d.Coordinate(1) > 0 {
		y = b.Point2().Coordinate(1)
	} else {
		y = b.Point1().Coordinate(1)
	}

	if d.Coordinate(2) > 0 { // nolint:mnd
		z = b.Point2().Coordinate(2) // nolint:mnd
	} else {
		z = b.Point1().Coordinate(2) // nolint:mnd
	}

	return NewPoint(x, y, z)
}

// Split split box into 8 boxes
func (b Box) Split() [8]Box {
	center := b.Center()
	mn := b.Point1()
	mx := b.Point2()

	return [8]Box{
		NewBox(NewPoint(mn.Coordinate(0), mn.Coordinate(1), mn.Coordinate(2)), center.Coordinate(0), // nolint:mnd
			center.Coordinate(1), center.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(mn.Coordinate(0), mn.Coordinate(1), center.Coordinate(2)), center.Coordinate(0), // nolint:mnd
			center.Coordinate(1), mx.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(mn.Coordinate(0), center.Coordinate(1), mn.Coordinate(2)), center.Coordinate(0), // nolint:mnd
			mx.Coordinate(1), center.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(mn.Coordinate(0), center.Coordinate(1), center.Coordinate(2)), center.Coordinate(0), // nolint:mnd
			mx.Coordinate(1), mx.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(center.Coordinate(0), mn.Coordinate(1), mn.Coordinate(2)), mx.Coordinate(0), // nolint:mnd
			center.Coordinate(1), center.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(center.Coordinate(0), mn.Coordinate(1), center.Coordinate(2)), mx.Coordinate(0), // nolint:mnd
			center.Coordinate(1), mx.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(center.Coordinate(0), center.Coordinate(1), mn.Coordinate(2)), mx.Coordinate(0), // nolint:mnd
			mx.Coordinate(1), center.Coordinate(2)), // nolint:mnd
		NewBox(NewPoint(center.Coordinate(0), center.Coordinate(1), center.Coordinate(2)), mx.Coordinate(0), // nolint:mnd
			mx.Coordinate(1), mx.Coordinate(2)), // nolint:mnd
	}
}

// Contains returns whether the specified point is contained in this box.
func (b Box) Contains(v Point) bool {
	mn := b.Point1()
	mx := b.Point2()

	return (mn.Coordinate(0) <= v.Coordinate(0) && v.Coordinate(0) <= mx.Coordinate(0)) && // nolint:mnd
		(mn.Coordinate(1) <= v.Coordinate(1) && v.Coordinate(1) <= mx.Coordinate(1)) && // nolint:mnd
		(mn.Coordinate(2) <= v.Coordinate(2) && v.Coordinate(2) <= mx.Coordinate(2)) // nolint:mnd
}

// Fit returns whether the specified area is fully contained in the other area.
func (b Box) Fit(o Box) bool {
	return o.Contains(b.Point2()) && o.Contains(b.Point1())
}
