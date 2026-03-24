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

func TestRadianToDegree(t *testing.T) {
	tests := map[string]struct {
		radian float64
		want   float64
	}{
		"zero": {
			radian: 0,
			want:   0,
		},
		"pi": {
			radian: 3.141592653589793,
			want:   180,
		},
		"pi/2": {
			radian: 1.5707963267948966,
			want:   90,
		},
		"2pi": {
			radian: 6.283185307179586,
			want:   0,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, mathutils.RoundWithPrecision(RadianToDegree(tc.radian), 0.1), tc.want)
		})
	}
}

func TestGetAngle2D(t *testing.T) {
	tests := map[string]struct {
		v1 Point
		v2 Point
	}{
		"same point": {
			v1: NewPoint(0, 0),
			v2: NewPoint(0, 0),
		},
		"right": {
			v1: NewPoint(1, 0),
			v2: NewPoint(0, 0),
		},
		"up": {
			v1: NewPoint(0, 1),
			v2: NewPoint(0, 0),
		},
		"diagonal": {
			v1: NewPoint(1, 1),
			v2: NewPoint(0, 0),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			angle := GetAngle2D(tc.v1, tc.v2)
			_ = angle // should not panic
		})
	}

	// Test specific known angle
	angle := GetAngle2D(NewPoint(1, 0), NewPoint(0, 0))
	testutils.Equal(t, mathutils.RoundWithPrecision(angle, 0.0001), 0.0)
}

func TestGetDiffPoint2D(t *testing.T) {
	tests := map[string]struct {
		angle    float64
		velocity float64
	}{
		"zero angle": {
			angle:    0,
			velocity: 10,
		},
		"90 degrees": {
			angle:    1.5707963267948966,
			velocity: 10,
		},
		"zero velocity": {
			angle:    1,
			velocity: 0,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			p := GetDiffPoint2D(tc.angle, tc.velocity)
			testutils.Equal(t, p.Coordinate(2), 0.0) // z should always be 0
		})
	}

	// Test specific: angle=0, velocity=10 should give (10,0,0)
	p := GetDiffPoint2D(0, 10)
	testutils.Equal(t, p.Round(0.01), NewPoint(10, 0, 0))
}

func TestCoordinatesToPoint(t *testing.T) {
	c := [3]float64{1, 2, 3}
	p := CoordinatesToPoint(c)
	testutils.Equal(t, p, NewPoint(1, 2, 3))

	c2 := [3]float64{0, 0, 0}
	p2 := CoordinatesToPoint(c2)
	testutils.Equal(t, p2, NewPoint(0, 0, 0))
}
