package shapes

import (
	"math"

	"github.com/InsideGallery/core/mathutils"
)

// Point describe point
type Point struct {
	coordinates [3]float64
}

// NewPoint return new point
func NewPoint(s ...float64) Point {
	var x, y, z float64
	l := len(s)

	switch l {
	case 3: // nolint:mnd
		x, y, z = s[0], s[1], s[2]
	case 2: // nolint:mnd
		x, y = s[0], s[1]
	case 1: // nolint:mnd
		x = s[0]
	}

	return Point{
		coordinates: [3]float64{x, y, z},
	}
}

// DistanceSquare return distance square between positions (euclidean)
func (p Point) DistanceSquare(p2 Point) float64 {
	var sum float64

	for i := range p.Coordinates() {
		d := p.Coordinate(i) - p2.Coordinate(i)
		sum += d * d
	}

	return sum
}

// ManhattanDistance return manhattan distance
func (p Point) ManhattanDistance(p2 Point) float64 {
	var sum float64

	for i := range p.Coordinates() {
		d := math.Abs(p.Coordinate(i) - p2.Coordinate(i))
		sum += d
	}

	return sum
}

// Distance return distance between positions (euclidean)
func (p Point) Distance(p2 Point) float64 {
	return math.Sqrt(p.DistanceSquare(p2))
}

// Dot scalar multiply
func (p Point) Dot(p2 Point) float64 {
	var sum float64
	for i, v := range p.Coordinates() {
		sum += v * p2.Coordinate(i)
	}

	return sum
}

// Normal returns the vector's norm.
func (p Point) Normal() float64 {
	return math.Sqrt(p.Dot(p))
}

// NormalSquare returns the vector's norm square.
func (p Point) NormalSquare() float64 {
	return p.Dot(p)
}

// MinDistance computes the square of the distance from a point to a rectangle.
func (p Point) MinDistance(r Box) float64 {
	var sum float64

	for i, pi := range p.Coordinates() {
		if pi < r.Point1().Coordinate(i) { //nolint:gocritic
			d := pi - r.Point1().Coordinate(i)
			sum += d * d
		} else if pi > r.Point2().Coordinate(i) {
			d := pi - r.Point2().Coordinate(i)
			sum += d * d
		} else {
			sum += 0
		}
	}

	return sum
}

// MinMaxDistance computes the minimum of the maximum distances from p to points on r.
func (p Point) MinMaxDistance(r Box) float64 {
	// by definition, MinMaxDist(p, r) =
	// min{1<=k<=n}(|pk - rmk|^2 + sum{1<=i<=n, i != k}(|pi - rMi|^2))
	// where rmk and rMk are defined as follows:
	rm := func(k int) float64 {
		if p.Coordinate(k) <= (r.Point1().Coordinate(k)+r.Point2().Coordinate(k))/2 {
			return r.Point1().Coordinate(k)
		}

		return r.Point2().Coordinate(k)
	}

	rM := func(k int) float64 {
		if p.Coordinate(k) >= (r.Point1().Coordinate(k)+r.Point2().Coordinate(k))/2 {
			return r.Point1().Coordinate(k)
		}

		return r.Point2().Coordinate(k)
	}

	// This formula can be computed in linear time by precomputing
	// S = sum{1<=i<=n}(|pi - rMi|^2).

	var S float64

	for i := range p.Coordinates() {
		d := p.Coordinate(i) - rM(i)
		S += d * d
	}

	// Compute MinMaxDist using the precomputed S.
	smallest := math.MaxFloat64

	for k := range p.Coordinates() {
		d1 := p.Coordinate(k) - rM(k)
		d2 := p.Coordinate(k) - rm(k)
		d := S - d1*d1 + d2*d2

		if d < smallest {
			smallest = d
		}
	}

	return smallest
}

// Coordinate return coordinate for dimension (0-x, 1-y, 2-z)
func (p Point) Coordinate(i int) float64 {
	if i > 2 || i < 0 {
		return 0
	} // nolint:mnd

	return p.coordinates[i]
}

// Coordinates return all coordinates
func (p Point) Coordinates() [3]float64 {
	return p.coordinates
}

// Get return Spatial
func (p Point) Get() Spatial {
	return p
}

// Bounds return rectangle of object
func (p Point) Bounds() Box {
	return NewBox(p, 1, 1, 1) // nolint:mnd
}

// Scale multiply point on given point
func (p Point) Scale(v float64) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i * v
	}

	return NewPoint(np...)
}

// Increase add value to given point
func (p Point) Increase(v float64) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i + v
	}

	return NewPoint(np...)
}

// Decrease remove value from given point
func (p Point) Decrease(v float64) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i - v
	}

	return NewPoint(np...)
}

// GetMinAxis return the axis with the minimal value
func (p Point) GetMinAxis() int {
	if p.Coordinate(0) < p.Coordinate(1) { // nolint:mnd
		if p.Coordinate(0) < p.Coordinate(2) { // nolint:mnd
			return 0
		} // nolint:mnd

		return 2 //nolint:mnd
	}

	if p.Coordinate(1) < p.Coordinate(2) { // nolint:mnd
		return 1 //nolint:mnd
	} // nolint:mnd

	return 2 //nolint:mnd
}

// Normalize normalize point
func (p Point) Normalize() Point {
	n2 := math.Abs(p.Normal())
	if mathutils.ApproximatelyEqual(n2, 0) {
		return NewPoint(0, 0, 0)
	} // nolint:mnd

	r := make([]float64, 3) // nolint:mnd
	for i, c := range p.Coordinates() {
		r[i] = c / n2
	}

	return NewPoint(r...)
}

// Move return new object for given diff
func (p Point) Move(diff Point) Spatial {
	return p.Add(diff)
}

// Add add diff to p
func (p Point) Add(diff Point) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i + diff.Coordinate(k)
	}

	return NewPoint(np...)
}

// Round round all coordinate with precision
func (p Point) Round(precision float64) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = mathutils.RoundWithPrecision(i, precision)
	}

	return NewPoint(np...)
}

// Subtract get diff between points
func (p Point) Subtract(diff Point) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i - diff.Coordinate(k)
	}

	return NewPoint(np...)
}

// Multiply get multiple point
func (p Point) Multiply(diff Point) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = i * diff.Coordinate(k)
	}

	return NewPoint(np...)
}

// Divide get divide point
func (p Point) Divide(diff Point) Point {
	np := make([]float64, 3) // nolint:mnd

	for k, i := range p.Coordinates() {
		c := diff.Coordinate(k)
		if c != 0 {
			np[k] = i / c
		} else {
			np[k] = i
		}
	}

	return NewPoint(np...)
}

// Cross calculate cross
func (p Point) Cross(p2 Point) Point {
	x := p.Coordinate(1)*p2.Coordinate(2) - p.Coordinate(2)*p2.Coordinate(1) // nolint:mnd
	y := p.Coordinate(2)*p2.Coordinate(0) - p.Coordinate(0)*p2.Coordinate(2) // nolint:mnd
	z := p.Coordinate(0)*p2.Coordinate(1) - p.Coordinate(1)*p2.Coordinate(0) // nolint:mnd

	return NewPoint(x, y, z)
}

// Center return center point
func (p Point) Center() Point {
	return p
}

// Point1 return point1 point
func (p Point) Point1() Point {
	return p
}

// Copy return copy point
func (p Point) Copy() Point {
	np := make([]float64, 3) // nolint:mnd
	for k, v := range p.Coordinates() {
		np[k] = v
	}

	return NewPoint(np...)
}

// Invert invert the point
func (p Point) Invert() Point {
	np := make([]float64, 3) // nolint:mnd
	for k, v := range p.Coordinates() {
		np[k] = v * -1 // nolint:mnd
	}

	return NewPoint(np...)
}

// Abs abs the point
func (p Point) Abs() Point {
	np := make([]float64, 3) // nolint:mnd
	for k, v := range p.Coordinates() {
		np[k] = math.Abs(v)
	}

	return NewPoint(np...)
}

// Equals return true if points approximately equal
func (p Point) Equals(p2 Point) bool {
	for k, v := range p.Coordinates() {
		if !mathutils.ApproximatelyEqual(v, p2.Coordinate(k)) {
			return false
		}
	}

	return true
}

// Reflect reflection by normal
func (p Point) Reflect(normal Point) Point {
	return p.Subtract(normal.Scale(2 * p.Dot(normal))) // nolint:mnd
}

// Refract refraction by normal and eta
func (p Point) Refract(normal Point, eta float64) Point {
	n := normal.Dot(p)

	k := 1 - eta*eta*(1-n*n) // nolint:mnd
	if k < 0 {
		return NewPoint()
	}

	return p.Scale(eta).Subtract(normal.Scale(eta*n + math.Sqrt(k)))
}

// Support return support point
func (p Point) Support(_ Point) Point {
	return p
}

// Min Returns the a vector where each component is the lesser of the
// corresponding component in this and the specified vector
func (p Point) Min(other Point) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = math.Min(i, other.Coordinate(k))
	}

	return NewPoint(np...)
}

// Max Returns the a vector where each component is the greater of the
// corresponding component in this and the specified vector
func (p Point) Max(other Point) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = math.Max(i, other.Coordinate(k))
	}

	return NewPoint(np...)
}

// Lerp Returns the linear interpolation between two Point(s)
func (p Point) Lerp(other Point, f float64) Point {
	np := make([]float64, 3) // nolint:mnd
	for k, i := range p.Coordinates() {
		np[k] = (other.Coordinate(k)-i)*f + i
	}

	return NewPoint(np...)
}

// IsEmpty return true if this is zero point
func (p Point) IsEmpty() bool {
	for _, i := range p.Coordinates() {
		if i != 0 {
			return false
		}
	}

	return true
}
