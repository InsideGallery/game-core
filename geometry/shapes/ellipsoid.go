package shapes

// Ellipsoid contains point and sizes
type Ellipsoid struct {
	Point
	sizes [3]float64
}

// NewEllipsoid return new circle
func NewEllipsoid(p Point, s ...float64) Ellipsoid {
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

	return Ellipsoid{
		Point: p,
		sizes: [3]float64{width, height, length},
	}
}

// Point1 return center point
func (e Ellipsoid) Point1() Point {
	return e.Point
}

// Size return size by dimension
func (e Ellipsoid) Size(i int) float64 {
	if i > 2 || i < 0 {
		return 0
	}

	return e.sizes[i]
}

// Sizes return all sizes
func (e Ellipsoid) Sizes() [3]float64 {
	return e.sizes
}

// Move return new object for given diff
func (e Ellipsoid) Move(diff Point) Spatial {
	d := make([]float64, 3) // nolint:mnd
	for k, v := range e.Coordinates() {
		d[k] = v + diff.Coordinate(k)
	}

	return NewEllipsoid(NewPoint(d...), e.Size(0), e.Size(1), e.Size(2)) // nolint:mnd
}

// Get return Spatial
func (e Ellipsoid) Get() Spatial {
	return e
}

// Bounds return rectangle of object
func (e Ellipsoid) Bounds() Box {
	points := make([]float64, 3) // nolint:mnd
	sizes := make([]float64, 3)  // nolint:mnd

	for k, v := range e.Point1().Coordinates() {
		points[k] = v - e.Size(k)
		sizes[k] = e.Size(k) * 2 // nolint:mnd
	}

	return NewBox(NewPoint(points...), sizes...)
}

// Center return center point
func (e Ellipsoid) Center() Point {
	return e.Point1()
}

// Support return support point for sphere
func (e Ellipsoid) Support(d Point) Point {
	return d.Multiply(NewPoint(e.Size(0), e.Size(1), e.Size(2))).Normalize(). // nolint:mnd
											Multiply(NewPoint(e.Size(0), e.Size(1), e.Size(2))). // nolint:mnd
											Add(e.Point1())
}
