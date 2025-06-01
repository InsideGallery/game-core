package quickhull

import (
	"github.com/InsideGallery/game-core/geometry/shapes"
)

type ray struct {
	s                 shapes.Point
	v                 shapes.Point
	vInvLengthSquared float64
}

func newRay(s, v shapes.Point) ray {
	return ray{s: s, v: v, vInvLengthSquared: 1 / v.Dot(v)}
}
