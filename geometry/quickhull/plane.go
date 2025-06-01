package quickhull

import (
	"github.com/InsideGallery/game-core/geometry/shapes"
)

type plane struct {
	n          shapes.Point
	d          float64 // Signed distance (if normal is of length 1) to the plane from origin
	sqrNLength float64 // Normal length squared
}

func (p plane) isPointOnPositiveSide(q shapes.Point) bool {
	return p.n.Dot(q)+p.d >= 0
}

func newPlane(n shapes.Point, p shapes.Point) plane {
	return plane{n: n, d: -n.Dot(p), sqrNLength: n.Dot(n)}
}
