package shapes

// Line describe line
type Line struct {
	Point
	point2 Point
}

// NewLine return new line
func NewLine(p Point, p2 Point) Line {
	return Line{
		Point:  p,
		point2: p2,
	}
}

// Get return Spatial
func (l Line) Get() Spatial {
	return l
}

// Point1 return first point
func (l Line) Point1() Point {
	return l.Point
}

// Point2 return second point
func (l Line) Point2() Point {
	return l.point2
}

// Move return new object for given diff
func (l Line) Move(diff Point) Spatial {
	p1 := make([]float64, 3) // nolint:mnd
	p2 := make([]float64, 3) // nolint:mnd

	for k, v := range diff.Coordinates() {
		p1[k] = v + l.Point1().Coordinate(k)
		p2[k] = v + l.Point2().Coordinate(k)
	}

	return NewLine(NewPoint(p1...), NewPoint(p2...))
}

// Bounds return box of object
func (l Line) Bounds() Box {
	p := make([]float64, 3) // nolint:mnd
	s := make([]float64, 3) // nolint:mnd

	for k, v := range l.Point1().Coordinates() {
		if v < l.Point2().Coordinate(k) {
			p[k] = v
			s[k] = l.Point2().Coordinate(k) - p[k]
		} else {
			p[k] = l.Point2().Coordinate(k)
			s[k] = v - p[k]
		}
	}

	return NewBox(NewPoint(p...), s...)
}

// Center return center point
func (l Line) Center() Point {
	p := make([]float64, 3) // nolint:mnd
	for k, v := range l.Point1().Coordinates() {
		p[k] = (v + l.Point2().Coordinate(k)) / 2 // nolint:mnd
	}

	return NewPoint(p...)
}

// Support return support point for box
func (l Line) Support(d Point) Point {
	if l.Point1().Dot(d) > l.Point2().Dot(d) {
		return l.Point1()
	}

	return l.Point2()
}

// ProjectForPoint project point on line
func (l Line) ProjectForPoint(p Point) Point {
	v1, v2 := l.Point1(), l.Point2()
	e1 := v2.Subtract(v1)
	e2 := p.Subtract(v1)
	valDp := e1.Dot(e2)
	len2 := e1.NormalSquare()

	return v1.Add(e1.Scale(valDp).Scale(1 / len2))
}
