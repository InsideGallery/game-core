package shapes

import (
	"math"
)

// Border special type of Box. Collision will return true, if object out of border
type Border struct {
	Box
}

// NewBorder return new border
func NewBorder(r Box) Border {
	return Border{
		Box: r,
	}
}

// Get return Spatial
func (b Border) Get() Spatial {
	return b
}

// Collision check collisions
func (b Border) Collision(o Spatial, dimensions int) (bool, []float64) {
	depth := make([]float64, dimensions)

	switch v := o.Get().(type) {
	case Point:
		for k := 0; k < dimensions; k++ {
			if v.Point1().Coordinate(k) < b.Point1().Coordinate(k) {
				depth[k] = b.Point1().Coordinate(k) + v.Point1().Coordinate(k)
			} else if v.Point1().Coordinate(k) > b.Point2().Coordinate(k) {
				depth[k] = v.Point1().Coordinate(k) - b.Point2().Coordinate(k)
			}
		}
	case Sphere:
		for k := 0; k < dimensions; k++ {
			if v.Point1().Coordinate(k)-v.Radius() < b.Point1().Coordinate(k) {
				depth[k] = b.Point1().Coordinate(k) + v.Point1().Coordinate(k) - v.Radius()
			} else if v.Point1().Coordinate(k)+v.Radius() > b.Point2().Coordinate(k) {
				depth[k] = v.Point1().Coordinate(k) + v.Radius() - b.Point2().Coordinate(k)
			}
		}
	case Box:
		for k := 0; k < dimensions; k++ {
			if v.Point1().Coordinate(k) < b.Point1().Coordinate(k) {
				depth[k] = b.Point1().Coordinate(k) + v.Point1().Coordinate(k)
			} else if v.Point2().Coordinate(k) > b.Point2().Coordinate(k) {
				depth[k] = v.Point1().Coordinate(k) - b.Point2().Coordinate(k)
			}
		}
	case Line:
		for k := 0; k < dimensions; k++ {
			switch {
			case v.Point1().Coordinate(k) < b.Point1().Coordinate(k):
				depth[k] = b.Point1().Coordinate(k) + v.Point1().Coordinate(k)
			case v.Point1().Coordinate(k) > b.Point2().Coordinate(k):
				depth[k] = v.Point2().Coordinate(k) - b.Point2().Coordinate(k)
			case v.Point2().Coordinate(k) < b.Point1().Coordinate(k):
				depth[k] = b.Point1().Coordinate(k) + v.Point2().Coordinate(k)
			case v.Point2().Coordinate(k) > b.Point2().Coordinate(k):
				depth[k] = v.Point2().Coordinate(k) - b.Point2().Coordinate(k)
			}
		}

	case Polyhedron:
		next := 0
		for current := 0; current < v.Count(); current++ {
			next = current + 1
			if next == v.Count() {
				next = 0
			}

			vc := v.Vector(current)
			vn := v.Vector(next)

			collision, depth := b.Collision(NewLine(vc, vn), dimensions)
			if collision {
				return true, depth
			}
		}
	case MultiObject:
		for _, o := range v.Objects() {
			if collision, depth := b.Collision(o, dimensions); collision {
				return true, depth
			}
		}
	}

	for _, v := range depth {
		if math.Abs(v) > 0 {
			return true, depth
		}
	}

	return false, depth
}
