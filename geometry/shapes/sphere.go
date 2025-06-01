package shapes

// Sphere contains point and radius
type Sphere struct {
	Point
	radius float64
}

// NewSphere return new circle
func NewSphere(p Point, radius float64) Sphere {
	return Sphere{
		Point:  p,
		radius: radius,
	}
}

// Point1 return center point
func (s Sphere) Point1() Point {
	return s.Point
}

// Radius return radius
func (s Sphere) Radius() float64 {
	return s.radius
}

// Move return new object for given diff
func (s Sphere) Move(diff Point) Spatial {
	d := make([]float64, 3) // nolint:mnd
	for k, v := range s.Coordinates() {
		d[k] = v + diff.Coordinate(k)
	}

	return NewSphere(NewPoint(d...), s.Radius())
}

// Get return Spatial
func (s Sphere) Get() Spatial {
	return s
}

// Bounds return rectangle of object
func (s Sphere) Bounds() Box {
	points := make([]float64, 3) // nolint:mnd
	sizes := make([]float64, 3)  // nolint:mnd

	for k, v := range s.Point1().Coordinates() {
		points[k] = v - s.Radius()
		sizes[k] = 2 * s.Radius() // nolint:mnd
	}

	return NewBox(NewPoint(points...), sizes...)
}

// CollisionSphere return true if two spheres collided
func (s Sphere) CollisionSphere(s2 Sphere) bool {
	d := s.Point1().DistanceSquare(s2.Point1())
	r := s.Radius() + s2.Radius()

	return d <= r*r
}

// Center return center point
func (s Sphere) Center() Point {
	return s.Point1()
}

// Support return support point for sphere
func (s Sphere) Support(d Point) Point {
	return d.Normalize().Scale(s.Radius()).Add(s.Point1())
}
