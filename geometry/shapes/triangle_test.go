package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewTriangle_ZeroArgs(t *testing.T) {
	tr := NewTriangle()
	testutils.Equal(t, tr.Point1(), NewPoint())
}

func TestNewTriangle_OneArg(t *testing.T) {
	tr := NewTriangle(NewPoint(1, 2, 3))
	testutils.Equal(t, tr.points[0], NewPoint(1, 2, 3))
	testutils.Equal(t, tr.points[1], NewPoint())
	testutils.Equal(t, tr.points[2], NewPoint())
}

func TestNewTriangle_TwoArgs(t *testing.T) {
	tr := NewTriangle(NewPoint(1, 2, 3), NewPoint(4, 5, 6))
	testutils.Equal(t, tr.points[0], NewPoint(1, 2, 3))
	testutils.Equal(t, tr.points[1], NewPoint(4, 5, 6))
	testutils.Equal(t, tr.points[2], NewPoint())
}

func TestNewTriangle_ThreeArgs(t *testing.T) {
	tr := NewTriangle(NewPoint(1, 2, 3), NewPoint(4, 5, 6), NewPoint(7, 8, 9))
	testutils.Equal(t, tr.points[0], NewPoint(1, 2, 3))
	testutils.Equal(t, tr.points[1], NewPoint(4, 5, 6))
	testutils.Equal(t, tr.points[2], NewPoint(7, 8, 9))
}

func TestTriangle_Coordinates(t *testing.T) {
	tr := NewTriangle(NewPoint(1, 5, 0), NewPoint(3, 1, 0), NewPoint(7, 3, 0))
	coords := tr.Coordinates()
	// Bounds min should be (1,1,0)
	testutils.Equal(t, coords[0], 1.0)
	testutils.Equal(t, coords[1], 1.0)
	testutils.Equal(t, coords[2], 0.0)
}

func TestTriangle_Point1(t *testing.T) {
	tr := NewTriangle(NewPoint(10, 20, 30), NewPoint(40, 50, 60), NewPoint(70, 80, 90))
	testutils.Equal(t, tr.Point1(), NewPoint(10, 20, 30))
}

func TestTriangle_Bounds(t *testing.T) {
	tests := map[string]struct {
		a, b, c Point
		wantP   Point
		wantS   [3]float64
	}{
		"simple triangle": {
			a:     NewPoint(0, 0, 0),
			b:     NewPoint(10, 0, 0),
			c:     NewPoint(5, 10, 0),
			wantP: NewPoint(0, 0, 0),
			wantS: [3]float64{10, 10, 0},
		},
		"negative coords": {
			a:     NewPoint(-5, -5, 0),
			b:     NewPoint(5, 0, 0),
			c:     NewPoint(0, 5, 0),
			wantP: NewPoint(-5, -5, 0),
			wantS: [3]float64{10, 10, 0},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewTriangle(tc.a, tc.b, tc.c)
			b := tr.Bounds()
			testutils.Equal(t, b.Point1(), tc.wantP)
			testutils.Equal(t, b.Sizes(), tc.wantS)
		})
	}
}

func TestTriangle_Center(t *testing.T) {
	tr := NewTriangle(NewPoint(0, 0, 0), NewPoint(3, 0, 0), NewPoint(0, 3, 0))
	c := tr.Center()
	testutils.Equal(t, c.Coordinate(0), 1.0)
	testutils.Equal(t, c.Coordinate(1), 1.0)
}

func TestTriangle_Get(t *testing.T) {
	tr := NewTriangle(NewPoint(0, 0, 0), NewPoint(1, 0, 0), NewPoint(0, 1, 0))
	got := tr.Get()
	_, ok := got.(Triangle)
	testutils.Equal(t, ok, true)
}

func TestTriangle_Move(t *testing.T) {
	tr := NewTriangle(NewPoint(0, 0, 0), NewPoint(10, 0, 0), NewPoint(5, 10, 0))
	moved := tr.Move(NewPoint(5, 5, 5)).(Triangle)
	testutils.Equal(t, moved.points[0], NewPoint(5, 5, 5))
	testutils.Equal(t, moved.points[1], NewPoint(15, 5, 5))
	testutils.Equal(t, moved.points[2], NewPoint(10, 15, 5))
}

func TestTriangle_CalculateSurfaceNormal(t *testing.T) {
	tests := map[string]struct {
		a, b, c Point
		wantX   float64
		wantY   float64
		wantZ   float64
	}{
		"xy plane triangle": {
			a:     NewPoint(0, 0, 0),
			b:     NewPoint(1, 0, 0),
			c:     NewPoint(0, 1, 0),
			wantX: 0,
			wantY: 0,
			wantZ: 1,
		},
		"xz plane triangle": {
			a:     NewPoint(0, 0, 0),
			b:     NewPoint(1, 0, 0),
			c:     NewPoint(0, 0, 1),
			wantX: 0,
			wantY: -1,
			wantZ: 0,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewTriangle(tc.a, tc.b, tc.c)
			n := tr.CalculateSurfaceNormal()
			testutils.Equal(t, n.Coordinate(0), tc.wantX)
			testutils.Equal(t, n.Coordinate(1), tc.wantY)
			testutils.Equal(t, n.Coordinate(2), tc.wantZ)
		})
	}
}

func TestTriangle_Support(t *testing.T) {
	tests := map[string]struct {
		a, b, c Point
		dir     Point
		want    Point
	}{
		"point0 is furthest": {
			a:    NewPoint(10, 0, 0),
			b:    NewPoint(0, 5, 0),
			c:    NewPoint(0, 0, 5),
			dir:  NewPoint(1, 0, 0),
			want: NewPoint(10, 0, 0),
		},
		"point1 is furthest": {
			a:    NewPoint(0, 0, 0),
			b:    NewPoint(0, 10, 0),
			c:    NewPoint(0, 0, 5),
			dir:  NewPoint(0, 1, 0),
			want: NewPoint(0, 10, 0),
		},
		"point2 is furthest": {
			a:    NewPoint(0, 0, 0),
			b:    NewPoint(5, 0, 0),
			c:    NewPoint(0, 0, 10),
			dir:  NewPoint(0, 0, 1),
			want: NewPoint(0, 0, 10),
		},
		"point2 via else-if branch": {
			a:    NewPoint(5, 0, 0),
			b:    NewPoint(0, 0, 0),
			c:    NewPoint(0, 0, 10),
			dir:  NewPoint(0, 0, 1),
			want: NewPoint(0, 0, 10),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewTriangle(tc.a, tc.b, tc.c)
			testutils.Equal(t, tr.Support(tc.dir), tc.want)
		})
	}
}
