package quickhull

import (
	"math"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

const (
	minUint = uint(0)
	maxUint = ^minUint
	// minInt  = -maxInt - 1
	maxInt = int(maxUint >> 1)
)

func triangleNormal(a shapes.Point, b shapes.Point, c shapes.Point) shapes.Point {
	k := a.Subtract(c)
	rhs := b.Subtract(c)

	return shapes.NewPoint(k.Coordinate(1)*rhs.Coordinate(2)-k.Coordinate(2)*rhs.Coordinate(1), k.Coordinate(2)*rhs.Coordinate(0)-k.Coordinate(0)*rhs.Coordinate(2), k.Coordinate(0)*rhs.Coordinate(1)-k.Coordinate(1)*rhs.Coordinate(0)) //nolint:lll,mnd
}

func signedDistanceToPlane(v shapes.Point, p plane) float64 {
	return p.n.Dot(v) + p.d
}

func squaredDistanceBetweenPointAndRay(p shapes.Point, r ray) float64 {
	s := p.Subtract(r.s)
	t := s.Dot(r.v)

	return s.NormalSquare() - t*t*r.vInvLengthSquared
}

// extremeValues find indices of extreme values (max x, min x, max y, min y, max z, min z) for the given point cloud
func extremeValues(vertexData []shapes.Point) (extremeValueIndices [6]int) { //nolint:mnd
	vd0 := vertexData[0]
	extremeVals := [6]float64{ //nolint:mnd
		vd0.Coordinate(0),
		vd0.Coordinate(0),
		vd0.Coordinate(1),
		vd0.Coordinate(1), //nolint:mnd
		vd0.Coordinate(2), //nolint:mnd
		vd0.Coordinate(2), //nolint:mnd
	}

	xv := func(i, i2 int, v float64, idx int) {
		if v > extremeVals[i] {
			extremeVals[i] = v
			extremeValueIndices[i] = idx
		} else if v < extremeVals[i2] {
			extremeVals[i2] = v
			extremeValueIndices[i2] = idx
		}
	}

	for i, pos := range vertexData[1:] {
		idx := i + 1
		xv(0, 1, pos.Coordinate(0), idx)
		xv(2, 3, pos.Coordinate(1), idx) //nolint:mnd
		xv(4, 5, pos.Coordinate(2), idx) //nolint:mnd
	}

	return
}

type vectorField int

const (
	vfX vectorField = iota
	vfY
	vfZ
)

// scale compute scale of the vertex data.
func scale(vertexData []shapes.Point, extremeValueIndices [6]int) (s float64) {
	scl := func(i int, valType vectorField) {
		v := vertexData[extremeValueIndices[i]]
		var a float64

		switch valType {
		case vfX:
			a = v.Coordinate(0)
		case vfY:
			a = v.Coordinate(1) //nolint:mnd
		case vfZ:
			a = v.Coordinate(2) //nolint:mnd
		}

		s = math.Max(s, math.Abs(a))
	}

	scl(0, vfX)
	scl(1, vfX) //nolint:mnd
	scl(2, vfY) //nolint:mnd
	scl(3, vfY) //nolint:mnd
	scl(4, vfZ) //nolint:mnd
	scl(5, vfZ) //nolint:mnd

	return
}
