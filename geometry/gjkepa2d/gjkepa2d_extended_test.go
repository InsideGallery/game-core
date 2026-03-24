package gjkepa2d

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestRoundedGJKEPA2D(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100), 10)
	col, p := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, true)
	testutils.Equal(t, p.Normal() > 0, true)
}

func TestRoundedGJKEPA2DNoCollision(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(200, 200), 10)
	col, _ := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, false)
}

func TestRoundedGJKEPA2DSameCenter(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	col, p := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, true)
	testutils.Equal(t, p.Normal() > 0, true)
}

func TestSolveKinematicBody2D(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100), 10)
	force1 := shapes.NewPoint(1, 0)
	force2 := shapes.NewPoint(-1, 0)
	a, b := SolveKinematicBody(s1, s2, force1, force2, 0.8, 0.8)
	testutils.Equal(t, a.Center().Coordinate(0) < b.Center().Coordinate(0), true)
}

func TestSolveKinematicBody2DNoCollision(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(200, 200), 10)
	force1 := shapes.NewPoint(1, 0)
	force2 := shapes.NewPoint(-1, 0)
	a, b := SolveKinematicBody(s1, s2, force1, force2, 0.8, 0.8)
	testutils.Equal(t, a.Center(), s1.Center())
	testutils.Equal(t, b.Center(), s2.Center())
}

func TestGJK2DNoMTV(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100), 10)
	g := NewGJKEPA()
	r, p := g.GJK(s1, s2, false)
	testutils.Equal(t, r, true)
	var zero shapes.Point
	testutils.Equal(t, p, zero)
}

func TestGJK2DNoIntersection(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(0, 0), 5)
	s2 := shapes.NewSphere(shapes.NewPoint(100, 100), 5)
	g := NewGJKEPA()
	r, _ := g.GJK(s1, s2, true)
	testutils.Equal(t, r, false)
}

func TestGJK2DSphereVsBox(t *testing.T) {
	s := shapes.NewSphere(shapes.NewPoint(95, 150), 10)
	bb := shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	g := NewGJKEPA()
	r, p := g.GJK(s, bb, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Round(0.1), shapes.NewPoint(-5, 0, 0))
}

func TestGJK2DMultipleDirections(t *testing.T) {
	g := NewGJKEPA()
	tests := []struct {
		name string
		s    shapes.Sphere
		b    shapes.Box
	}{
		{"left", shapes.NewSphere(shapes.NewPoint(95, 300), 10), shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)},
		{"right", shapes.NewSphere(shapes.NewPoint(605, 300), 10), shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)},
		{"top", shapes.NewSphere(shapes.NewPoint(300, 95), 10), shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)},
		{"bottom", shapes.NewSphere(shapes.NewPoint(300, 605), 10), shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, p := g.GJK(tc.s, tc.b, true)
			testutils.Equal(t, r, true)
			testutils.Equal(t, p.Round(0.1).Normal() > 0, true)
		})
	}
}

func TestNewEdge(t *testing.T) {
	e := NewEdge(1.5, shapes.NewPoint(1, 0), 3)
	testutils.Equal(t, e.distance, 1.5)
	testutils.Equal(t, e.normal, shapes.NewPoint(1, 0))
	testutils.Equal(t, e.index, 3)
}
