package gjkepa3d

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestGJKEPA(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(100, 120, 100), 10)
	s3 := shapes.NewSphere(shapes.NewPoint(115, 100, 100), 10)
	m := shapes.NewMultiObject(s1, s2)

	g := NewGJKEPA()

	r, p := g.GJK(s1, s3, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Round(0.1), shapes.NewPoint(-5.9, -2.1, -3.4))

	r, p = g.GJK(s3, s1, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Round(0.1), shapes.NewPoint(5.9, 2.1, -3.4))
	testutils.Equal(t, s3.Add(p).Round(0.1), shapes.NewPoint(120.9, 102.1, 96.6))

	r, p = g.GJK(s2, s3, true)
	testutils.Equal(t, r, false)
	var v shapes.Point
	testutils.Equal(t, p, v)

	r, p = g.GJK(m, s3, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Round(0.1), shapes.NewPoint(-5, 0, 0))

	r, p = g.GJK(s3, m, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Round(0.1), shapes.NewPoint(5, 0, 0))

	testutils.Equal(t, s3.Add(p).Round(0.1), shapes.NewPoint(120, 100, 100))

	s3 = shapes.NewSphere(shapes.NewPoint(95, 110), 10)
	bb := shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp := g.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(-5, 0, 0))
	testutils.Equal(t, s3.Add(pp).Round(0.1), shapes.NewPoint(90, 110, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(606, 110), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = g.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(4, 0, 0))
	testutils.Equal(t, s3.Add(pp).Round(0.1), shapes.NewPoint(610, 110, -0))

	s3 = shapes.NewSphere(shapes.NewPoint(110, 581), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = g.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 0, 10))
	testutils.Equal(t, s3.Add(pp).Round(0.1), shapes.NewPoint(110, 581, 10))
	testutils.Equal(t, s3.Point1().Round(0.1), shapes.NewPoint(110, 581, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(110, 95), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = g.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, -5, 0))
	testutils.Equal(t, s3.Add(pp).Round(0.1), shapes.NewPoint(110, 90, 0))
}
