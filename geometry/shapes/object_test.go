package shapes

import (
	"testing"

	"github.com/InsideGallery/core/mathutils"
	"github.com/InsideGallery/core/testutils"
)

func TestSphere(t *testing.T) {
	testcases := map[string]struct {
		s1 Sphere
		s2 Sphere
		r  bool
	}{
		"sphere;z;true": {
			s1: NewSphere(NewPoint(160, 160, 60), 30),
			s2: NewSphere(NewPoint(160, 160, 120), 30),
			r:  true,
		},
		"sphere;z;false": {
			s1: NewSphere(NewPoint(160, 160, 60), 30),
			s2: NewSphere(NewPoint(160, 160, 120.1), 30),
			r:  false,
		},
		"sphere;y;true": {
			s1: NewSphere(NewPoint(160, 60, 160), 30),
			s2: NewSphere(NewPoint(160, 120, 160), 30),
			r:  true,
		},
		"sphere;y;false": {
			s1: NewSphere(NewPoint(160, 60, 160), 30),
			s2: NewSphere(NewPoint(160, 120.1, 160), 30),
			r:  false,
		},
		"sphere;x;true": {
			s1: NewSphere(NewPoint(60, 160, 160), 30),
			s2: NewSphere(NewPoint(120, 160, 160), 30),
			r:  true,
		},
		"sphere;x;false": {
			s1: NewSphere(NewPoint(60, 160, 160), 30),
			s2: NewSphere(NewPoint(120.1, 160, 160), 30),
			r:  false,
		},
		"circle;x;true": {
			s1: NewSphere(NewPoint(60, 160), 30),
			s2: NewSphere(NewPoint(120, 160), 30),
			r:  true,
		},
		"circle;x;false": {
			s1: NewSphere(NewPoint(60, 160), 30),
			s2: NewSphere(NewPoint(120.1, 160), 30),
			r:  false,
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.s1.CollisionSphere(test.s2), test.r)
		})
	}
}

func TestIntersect(t *testing.T) {
	testcases := map[string]struct {
		bb1 Box
		bb2 Box
		c   [3]float64
		s   [3]float64
	}{
		"2d;intersect": {
			bb1: NewBox(NewPoint(100, 100), 200, 200),
			bb2: NewBox(NewPoint(190, 190), 100, 100),
			c:   [3]float64{190, 190},
			s:   [3]float64{100, 100},
		},
		"3d;intersect": {
			bb1: NewBox(NewPoint(100, 100, 100), 200, 200, 200),
			bb2: NewBox(NewPoint(190, 190, 190), 100, 100, 100),
			c:   [3]float64{190, 190, 190},
			s:   [3]float64{100, 100, 100},
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			bb, _ := test.bb1.Intersect(test.bb2)
			c, s := bb.Coordinates(), bb.Sizes()
			testutils.Equal(t, c, test.c)
			testutils.Equal(t, s, test.s)
		})
	}
}

func TestSupport(t *testing.T) {
	s := NewSphere(NewPoint(100, 100, 100), 20)
	p := s.Support(NewPoint(1, 0, 0))
	testutils.Equal(t, p, NewPoint(120, 100, 100))

	pl := NewPolyhedron(NewPoint(10, 10, 0), NewPoint(10, 100, 0), NewPoint(100, 100, 0), NewPoint(100, 10, 0), NewPoint(10, 10, 0))
	p = pl.Support(NewPoint(0.5, -0.5, 0))
	testutils.Equal(t, p, NewPoint(100, 10, 0))

	bb := NewBox(NewPoint(100, 100, 100), 20, 20, 20)
	p = bb.Support(NewPoint(1, 0, 0))
	testutils.Equal(t, p, NewPoint(120, 100, 100))

	tr := NewTriangle(NewPoint(100, 100, 100), NewPoint(100, 100, 200), NewPoint(100, 200, 100))
	p = tr.Support(NewPoint(200, 100, 100))
	testutils.Equal(t, p, NewPoint(100, 100, 200))
}

func TestReflect(t *testing.T) {
	p := NewPoint(10, 10, 0)
	s := NewPoint(1, 0, 0)

	testutils.Equal(t, p.Reflect(s), NewPoint(-10, 10, 0))
}

func TestNormalizeDegrees(t *testing.T) {
	testcases := map[string]struct {
		value  float64
		result float64
	}{
		"value;0": {
			value:  0,
			result: 0,
		},
		"value;360": {
			value:  360,
			result: 0,
		},
		"value;360.9": {
			value:  360.9,
			result: 0.9,
		},
		"value;1243": {
			value:  1243,
			result: 163,
		},
		"value;-0.1": {
			value:  -0.1,
			result: 359.9,
		},
		"value;-1243": {
			value:  -1243,
			result: 197,
		},
		"value;-1079": {
			value:  -1079,
			result: 1,
		},
		"value;-1080": {
			value:  -1080,
			result: 0,
		},
		"value;-1081": {
			value:  -1081,
			result: 359,
		},
	}
	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, mathutils.RoundWithPrecision(NormalizeDegrees(test.value), 0.1), mathutils.RoundWithPrecision(test.result, 0.1))
		})
	}
}

func TestRotatePoint(t *testing.T) {
	p := RotatePoint(NewPoint(1, 1, 1), DegreesToRadian(0), DegreesToRadian(0), DegreesToRadian(0))
	testutils.Equal(t, p, NewPoint(1, 1, 1))
	p = RotatePoint(NewPoint(1, 1, 1), DegreesToRadian(90), DegreesToRadian(0), DegreesToRadian(0))
	testutils.Equal(t, p.Round(0.01), NewPoint(-1, 1, 1))
	p = RotatePoint(NewPoint(1, 1, 1), DegreesToRadian(0), DegreesToRadian(90), DegreesToRadian(0))
	testutils.Equal(t, p.Round(0.01), NewPoint(1, 1, -1))
	p = RotatePoint(NewPoint(1, 1, 1), DegreesToRadian(0), DegreesToRadian(0), DegreesToRadian(90))
	testutils.Equal(t, p.Round(0.01), NewPoint(1, -1, 1))
}

func TestRotateBy(t *testing.T) {
	p := RotateBy(NewPoint(1, 2, 3), DegreesToRadian(90), 0)
	testutils.Equal(t, p, NewPoint(1, -3, 2))
	p = RotateBy(NewPoint(1, 2, 3), DegreesToRadian(90), 1)
	testutils.Equal(t, p.Round(0.1), NewPoint(3, 2, -1))
	p = RotateBy(NewPoint(1, 2, 3), DegreesToRadian(90), 2)
	testutils.Equal(t, p.Round(0.1), NewPoint(-2, 1, 3))
}
