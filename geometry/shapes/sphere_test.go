package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewSphere(t *testing.T) {
	s := NewSphere(NewPoint(1, 2, 3), 5)
	testutils.Equal(t, s.Point1(), NewPoint(1, 2, 3))
	testutils.Equal(t, s.Radius(), 5.0)
}

func TestSphere_Point1(t *testing.T) {
	s := NewSphere(NewPoint(10, 20, 30), 15)
	testutils.Equal(t, s.Point1(), NewPoint(10, 20, 30))
}

func TestSphere_Radius(t *testing.T) {
	s := NewSphere(NewPoint(0, 0, 0), 42.5)
	testutils.Equal(t, s.Radius(), 42.5)
}

func TestSphere_Move(t *testing.T) {
	tests := map[string]struct {
		center  Point
		radius  float64
		diff    Point
		wantPos Point
		wantR   float64
	}{
		"move positive": {
			center:  NewPoint(1, 2, 3),
			radius:  5,
			diff:    NewPoint(10, 20, 30),
			wantPos: NewPoint(11, 22, 33),
			wantR:   5,
		},
		"move negative": {
			center:  NewPoint(10, 20, 30),
			radius:  10,
			diff:    NewPoint(-5, -10, -15),
			wantPos: NewPoint(5, 10, 15),
			wantR:   10,
		},
		"move zero": {
			center:  NewPoint(1, 2, 3),
			radius:  5,
			diff:    NewPoint(0, 0, 0),
			wantPos: NewPoint(1, 2, 3),
			wantR:   5,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := NewSphere(tc.center, tc.radius)
			moved := s.Move(tc.diff)
			ms := moved.(Sphere)
			testutils.Equal(t, ms.Point1(), tc.wantPos)
			testutils.Equal(t, ms.Radius(), tc.wantR)
		})
	}
}

func TestSphere_Get(t *testing.T) {
	s := NewSphere(NewPoint(1, 2, 3), 5)
	got := s.Get()
	_, ok := got.(Sphere)
	testutils.Equal(t, ok, true)
}

func TestSphere_Bounds(t *testing.T) {
	tests := map[string]struct {
		center Point
		radius float64
		wantP1 Point
		wantP2 Point
	}{
		"unit sphere at origin": {
			center: NewPoint(0, 0, 0),
			radius: 1,
			wantP1: NewPoint(-1, -1, -1),
			wantP2: NewPoint(1, 1, 1),
		},
		"sphere radius 5 at (10,10,10)": {
			center: NewPoint(10, 10, 10),
			radius: 5,
			wantP1: NewPoint(5, 5, 5),
			wantP2: NewPoint(15, 15, 15),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := NewSphere(tc.center, tc.radius)
			b := s.Bounds()
			testutils.Equal(t, b.Point1(), tc.wantP1)
			testutils.Equal(t, b.Point2(), tc.wantP2)
		})
	}
}

func TestSphere_CollisionSphere(t *testing.T) {
	tests := map[string]struct {
		s1   Sphere
		s2   Sphere
		want bool
	}{
		"overlapping": {
			s1:   NewSphere(NewPoint(0, 0, 0), 5),
			s2:   NewSphere(NewPoint(8, 0, 0), 5),
			want: true,
		},
		"touching": {
			s1:   NewSphere(NewPoint(0, 0, 0), 5),
			s2:   NewSphere(NewPoint(10, 0, 0), 5),
			want: true,
		},
		"separated": {
			s1:   NewSphere(NewPoint(0, 0, 0), 5),
			s2:   NewSphere(NewPoint(10.1, 0, 0), 5),
			want: false,
		},
		"concentric": {
			s1:   NewSphere(NewPoint(0, 0, 0), 5),
			s2:   NewSphere(NewPoint(0, 0, 0), 3),
			want: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.s1.CollisionSphere(tc.s2), tc.want)
		})
	}
}

func TestSphere_Center(t *testing.T) {
	s := NewSphere(NewPoint(5, 10, 15), 20)
	testutils.Equal(t, s.Center(), NewPoint(5, 10, 15))
}

func TestSphere_Support(t *testing.T) {
	tests := map[string]struct {
		center Point
		radius float64
		dir    Point
		want   Point
	}{
		"x direction": {
			center: NewPoint(0, 0, 0),
			radius: 10,
			dir:    NewPoint(1, 0, 0),
			want:   NewPoint(10, 0, 0),
		},
		"y direction": {
			center: NewPoint(0, 0, 0),
			radius: 10,
			dir:    NewPoint(0, 1, 0),
			want:   NewPoint(0, 10, 0),
		},
		"negative x direction": {
			center: NewPoint(0, 0, 0),
			radius: 10,
			dir:    NewPoint(-1, 0, 0),
			want:   NewPoint(-10, 0, 0),
		},
		"offset center": {
			center: NewPoint(100, 100, 100),
			radius: 20,
			dir:    NewPoint(1, 0, 0),
			want:   NewPoint(120, 100, 100),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := NewSphere(tc.center, tc.radius)
			testutils.Equal(t, s.Support(tc.dir), tc.want)
		})
	}
}
