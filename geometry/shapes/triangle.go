package shapes

import "math"

// Triangle describe triangle by 3 points
type Triangle struct {
	points [3]Point
}

// NewTriangle return new triangle
func NewTriangle(ps ...Point) Triangle {
	a, b, c := NewPoint(), NewPoint(), NewPoint()
	l := len(ps)

	switch l {
	case 3: // nolint:mnd
		a, b, c = ps[0], ps[1], ps[2]
	case 2: // nolint:mnd
		a, b = ps[0], ps[1]
	case 1: // nolint:mnd
		a = ps[0]
	}

	return Triangle{
		points: [3]Point{a, b, c},
	}
}

// Coordinates return coordinate for dimension (0-x, 1-y, 2-z)
func (t Triangle) Coordinates() [3]float64 {
	return t.Bounds().Coordinates()
}

// Point1 return point
func (t Triangle) Point1() Point {
	return t.points[0]
}

// Bounds return rectangle of object
func (t Triangle) Bounds() Box {
	np := make([]float64, 3) //nolint:mnd
	s := make([]float64, 3)  //nolint:mnd

	for i := 0; i < 2; i++ {
		s1 := t.points[0].Coordinate(i)
		s2 := t.points[1].Coordinate(i)
		s3 := t.points[2].Coordinate(i)

		lmax := math.Max(s3, math.Max(s1, s2))
		lmin := math.Min(s3, math.Min(s1, s2))

		np[i] = lmin
		s[i] = lmax - lmin
	}

	return NewBox(NewPoint(np...), s...)
}

// Center return center point
func (t Triangle) Center() Point {
	np := make([]float64, 3) //nolint:mnd

	for i := 0; i < 2; i++ {
		var sum float64
		for _, p := range t.points {
			sum += p.Coordinate(i)
		}

		np[i] = sum / 3 //nolint:mnd
	}

	return NewPoint(np...)
}

// Get return Spatial
func (t Triangle) Get() Spatial {
	return t
}

// Move return new object for given diff
func (t Triangle) Move(diff Point) Spatial {
	np := make([]Point, 3) //nolint:mnd

	for i, p := range t.points {
		np[i] = p.Add(diff)
	}

	return NewTriangle(np...)
}

// CalculateSurfaceNormal return triangle normal
func (t Triangle) CalculateSurfaceNormal() Point {
	U := t.points[1].Subtract(t.points[0])
	V := t.points[2].Subtract(t.points[0])

	return U.Cross(V)
}

// Support return support point for triangle
func (t Triangle) Support(d Point) Point {
	dot0 := t.points[0].Dot(d)
	dot1 := t.points[1].Dot(d)
	dot2 := t.points[2].Dot(d)

	furthestPoint := t.points[0]
	if dot1 > dot0 {
		furthestPoint = t.points[1]
		if dot2 > dot1 {
			furthestPoint = t.points[2]
		}
	} else if dot2 > dot0 {
		furthestPoint = t.points[2]
	}

	return furthestPoint
}
