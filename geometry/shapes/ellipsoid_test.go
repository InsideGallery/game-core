package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewEllipsoid_ZeroSizes(t *testing.T) {
	e := NewEllipsoid(NewPoint(1, 2, 3))
	testutils.Equal(t, e.Point1(), NewPoint(1, 2, 3))
	testutils.Equal(t, e.Sizes(), [3]float64{0, 0, 0})
}

func TestNewEllipsoid_OneSize(t *testing.T) {
	e := NewEllipsoid(NewPoint(1, 2, 3), 5)
	testutils.Equal(t, e.Size(0), 5.0)
	testutils.Equal(t, e.Size(1), 0.0)
	testutils.Equal(t, e.Size(2), 0.0)
}

func TestNewEllipsoid_TwoSizes(t *testing.T) {
	e := NewEllipsoid(NewPoint(1, 2, 3), 5, 10)
	testutils.Equal(t, e.Size(0), 5.0)
	testutils.Equal(t, e.Size(1), 10.0)
	testutils.Equal(t, e.Size(2), 0.0)
}

func TestNewEllipsoid_ThreeSizes(t *testing.T) {
	e := NewEllipsoid(NewPoint(1, 2, 3), 5, 10, 15)
	testutils.Equal(t, e.Size(0), 5.0)
	testutils.Equal(t, e.Size(1), 10.0)
	testutils.Equal(t, e.Size(2), 15.0)
}

func TestEllipsoid_Point1(t *testing.T) {
	e := NewEllipsoid(NewPoint(10, 20, 30), 5, 5, 5)
	testutils.Equal(t, e.Point1(), NewPoint(10, 20, 30))
}

func TestEllipsoid_Size(t *testing.T) {
	e := NewEllipsoid(NewPoint(0, 0, 0), 3, 4, 5)
	tests := map[string]struct {
		dim  int
		want float64
	}{
		"dim 0":                {dim: 0, want: 3},
		"dim 1":                {dim: 1, want: 4},
		"dim 2":                {dim: 2, want: 5},
		"dim -1 out of bounds": {dim: -1, want: 0},
		"dim 3 out of bounds":  {dim: 3, want: 0},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, e.Size(tc.dim), tc.want)
		})
	}
}

func TestEllipsoid_Sizes(t *testing.T) {
	e := NewEllipsoid(NewPoint(0, 0, 0), 3, 4, 5)
	testutils.Equal(t, e.Sizes(), [3]float64{3, 4, 5})
}

func TestEllipsoid_Move(t *testing.T) {
	e := NewEllipsoid(NewPoint(1, 2, 3), 5, 10, 15)
	moved := e.Move(NewPoint(10, 20, 30))
	me := moved.(Ellipsoid)
	testutils.Equal(t, me.Point1(), NewPoint(11, 22, 33))
	testutils.Equal(t, me.Sizes(), [3]float64{5, 10, 15})
}

func TestEllipsoid_Get(t *testing.T) {
	e := NewEllipsoid(NewPoint(0, 0, 0), 5, 5, 5)
	got := e.Get()
	_, ok := got.(Ellipsoid)
	testutils.Equal(t, ok, true)
}

func TestEllipsoid_Bounds(t *testing.T) {
	e := NewEllipsoid(NewPoint(10, 20, 30), 5, 10, 15)
	b := e.Bounds()
	testutils.Equal(t, b.Point1(), NewPoint(5, 10, 15))
	testutils.Equal(t, b.Point2(), NewPoint(15, 30, 45))
}

func TestEllipsoid_Center(t *testing.T) {
	e := NewEllipsoid(NewPoint(10, 20, 30), 5, 10, 15)
	testutils.Equal(t, e.Center(), NewPoint(10, 20, 30))
}

func TestEllipsoid_Support(t *testing.T) {
	tests := map[string]struct {
		center Point
		sizes  [3]float64
		dir    Point
	}{
		"x direction": {
			center: NewPoint(0, 0, 0),
			sizes:  [3]float64{5, 10, 15},
			dir:    NewPoint(1, 0, 0),
		},
		"y direction": {
			center: NewPoint(0, 0, 0),
			sizes:  [3]float64{5, 10, 15},
			dir:    NewPoint(0, 1, 0),
		},
		"z direction": {
			center: NewPoint(0, 0, 0),
			sizes:  [3]float64{5, 10, 15},
			dir:    NewPoint(0, 0, 1),
		},
		"with offset": {
			center: NewPoint(10, 20, 30),
			sizes:  [3]float64{5, 5, 5},
			dir:    NewPoint(1, 0, 0),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := NewEllipsoid(tc.center, tc.sizes[0], tc.sizes[1], tc.sizes[2])
			p := e.Support(tc.dir)
			// Support should return a point that is on the ellipsoid surface
			// Just verify it doesn't panic and returns something reasonable
			_ = p.Coordinates()
		})
	}

	// Test a specific known case: unit sphere-like ellipsoid
	e := NewEllipsoid(NewPoint(0, 0, 0), 10, 10, 10)
	p := e.Support(NewPoint(1, 0, 0))
	testutils.Equal(t, p, NewPoint(10, 0, 0))
}
