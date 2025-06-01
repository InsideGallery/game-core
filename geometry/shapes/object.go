package shapes

import (
	"math"
)

// Spatial describe boundable object
type Spatial interface {
	Point1() Point
	Coordinates() [3]float64
	Center() Point
	Get() Spatial
	Move(diff Point) Spatial
	Bounds() Box
}

// Collide describe collide interface
type Collide interface {
	Spatial
	Support(d Point) Point
}

// RadianToDegree return degree for given rad
func RadianToDegree(radian float64) float64 {
	return NormalizeDegrees((180 / math.Pi) * radian) // nolint:mnd
}

// DegreesToRadian return radian for given degrees
func DegreesToRadian(degrees float64) float64 {
	return (math.Pi / 180) * NormalizeDegrees(degrees) // nolint:mnd
}

// NormalizeDegrees return degrees n range 0-360
func NormalizeDegrees(value float64) float64 {
	return value + math.Ceil(-value/360)*360 //nolint:mnd
}

// GetAngle2D get angle between points in radian
func GetAngle2D(v1 Point, v2 Point) float64 {
	return math.Atan2(v1.Coordinate(1)-v2.Coordinate(1), v1.Coordinate(0)-v2.Coordinate(0)) // nolint:mnd
}

// GetDiffPoint2D return  diff point 2D
func GetDiffPoint2D(angle, velocity float64) Point {
	ps := make([]float64, 3) // nolint:mnd
	ps[0] = velocity * math.Cos(angle)
	ps[1] = velocity * math.Sin(angle)
	ps[2] = 0

	return NewPoint(ps...)
}

// CoordinatesToPoint return point by coordinates
func CoordinatesToPoint(c [3]float64) Point {
	p := make([]float64, 3) // nolint:mnd
	copy(p, c[:])

	return NewPoint(p...)
}

/*
yaw - around z
pitch - around y
roll - around x
*/

// RotatePoint rotate point by roll, yaw and pitch
func RotatePoint(p Point, yaw, pitch, roll float64) Point {
	cosa := math.Cos(yaw)
	sina := math.Sin(yaw)

	cosb := math.Cos(pitch)
	sinb := math.Sin(pitch)

	cosc := math.Cos(roll)
	sinc := math.Sin(roll)

	Axx := cosa * cosb
	Axy := cosa*sinb*sinc - sina*cosc
	Axz := cosa*sinb*cosc + sina*sinc

	Ayx := sina * cosb
	Ayy := sina*sinb*sinc + cosa*cosc
	Ayz := sina*sinb*cosc - cosa*sinc

	Azx := -sinb
	Azy := cosb * sinc
	Azz := cosb * cosc

	px := p.Coordinate(0) // nolint:mnd
	py := p.Coordinate(1) // nolint:mnd
	pz := p.Coordinate(2) // nolint:mnd

	nx := Axx*px + Axy*py + Axz*pz
	ny := Ayx*px + Ayy*py + Ayz*pz
	nz := Azx*px + Azy*py + Azz*pz

	return NewPoint(nx, ny, nz)
}

// RotateBy rotate by single coordinate
func RotateBy(p Point, angle float64, which int) Point {
	x, y, z := p.Coordinate(0), p.Coordinate(1), p.Coordinate(2) // nolint:mnd
	var dx, dy, dz float64

	switch which {
	case 0: // nolint:mnd
		dy = y*math.Cos(angle) - z*math.Sin(angle)
		dz = y*math.Sin(angle) + z*math.Cos(angle)
		dx = x
	case 1: // nolint:mnd
		dz = z*math.Cos(angle) - x*math.Sin(angle)
		dx = z*math.Sin(angle) + x*math.Cos(angle)
		dy = y
	case 2: // nolint:mnd
		dx = x*math.Cos(angle) - y*math.Sin(angle)
		dy = x*math.Sin(angle) + y*math.Cos(angle)
		dz = z
	}

	return NewPoint(dx, dy, dz)
}
