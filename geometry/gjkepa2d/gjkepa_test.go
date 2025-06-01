package gjkepa2d

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestCollision(t *testing.T) {
	cc := NewGJKEPA()
	s3 := shapes.NewSphere(shapes.NewPoint(95, 110), 10)
	bb := shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp := cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(-5, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))

	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))

	s3.Point = s3.Add(pp).Round(0.1)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))

	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(92, 110), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(-2, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(90, 110, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(608, 110), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(2, 0, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(610, 110, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(110, 92), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, -2, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(110, 90, 0))

	s3 = shapes.NewSphere(shapes.NewPoint(110, 608), 10)
	bb = shapes.NewBox(shapes.NewPoint(100, 100), 500, 500)
	_, pp = cc.GJK(s3, bb, true)
	testutils.Equal(t, pp.Round(0.1), shapes.NewPoint(0, 2, 0))
	s3.Point = s3.Add(pp).Round(0.1)
	testutils.Equal(t, s3.Point, shapes.NewPoint(110, 610, 0))
}
