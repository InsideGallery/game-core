package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewLine(t *testing.T) {
	l := NewLine(NewPoint(1, 2, 3), NewPoint(4, 5, 6))
	testutils.Equal(t, l.Point1(), NewPoint(1, 2, 3))
	testutils.Equal(t, l.Point2(), NewPoint(4, 5, 6))
}

func TestLine_Get(t *testing.T) {
	l := NewLine(NewPoint(0, 0, 0), NewPoint(1, 1, 1))
	got := l.Get()
	_, ok := got.(Line)
	testutils.Equal(t, ok, true)
}

func TestLine_Point1Point2(t *testing.T) {
	l := NewLine(NewPoint(10, 20, 30), NewPoint(40, 50, 60))
	testutils.Equal(t, l.Point1(), NewPoint(10, 20, 30))
	testutils.Equal(t, l.Point2(), NewPoint(40, 50, 60))
}

func TestLine_Move(t *testing.T) {
	tests := map[string]struct {
		p1     Point
		p2     Point
		diff   Point
		wantP1 Point
		wantP2 Point
	}{
		"move positive": {
			p1:     NewPoint(0, 0, 0),
			p2:     NewPoint(10, 10, 10),
			diff:   NewPoint(5, 5, 5),
			wantP1: NewPoint(5, 5, 5),
			wantP2: NewPoint(15, 15, 15),
		},
		"move negative": {
			p1:     NewPoint(10, 20, 30),
			p2:     NewPoint(40, 50, 60),
			diff:   NewPoint(-10, -20, -30),
			wantP1: NewPoint(0, 0, 0),
			wantP2: NewPoint(30, 30, 30),
		},
		"move zero": {
			p1:     NewPoint(1, 2, 3),
			p2:     NewPoint(4, 5, 6),
			diff:   NewPoint(0, 0, 0),
			wantP1: NewPoint(1, 2, 3),
			wantP2: NewPoint(4, 5, 6),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := NewLine(tc.p1, tc.p2)
			moved := l.Move(tc.diff).(Line)
			testutils.Equal(t, moved.Point1(), tc.wantP1)
			testutils.Equal(t, moved.Point2(), tc.wantP2)
		})
	}
}

func TestLine_Bounds(t *testing.T) {
	tests := map[string]struct {
		p1     Point
		p2     Point
		wantP1 Point
		wantS  [3]float64
	}{
		"ascending points": {
			p1:     NewPoint(1, 2, 3),
			p2:     NewPoint(4, 5, 6),
			wantP1: NewPoint(1, 2, 3),
			wantS:  [3]float64{3, 3, 3},
		},
		"descending points": {
			p1:     NewPoint(4, 5, 6),
			p2:     NewPoint(1, 2, 3),
			wantP1: NewPoint(1, 2, 3),
			wantS:  [3]float64{3, 3, 3},
		},
		"mixed": {
			p1:     NewPoint(5, 1, 3),
			p2:     NewPoint(2, 4, 6),
			wantP1: NewPoint(2, 1, 3),
			wantS:  [3]float64{3, 3, 3},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := NewLine(tc.p1, tc.p2)
			b := l.Bounds()
			testutils.Equal(t, b.Point1(), tc.wantP1)
			testutils.Equal(t, b.Sizes(), tc.wantS)
		})
	}
}

func TestLine_Center(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		want Point
	}{
		"simple": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 10, 10),
			want: NewPoint(5, 5, 5),
		},
		"asymmetric": {
			p1:   NewPoint(2, 4, 6),
			p2:   NewPoint(8, 12, 14),
			want: NewPoint(5, 8, 10),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := NewLine(tc.p1, tc.p2)
			testutils.Equal(t, l.Center(), tc.want)
		})
	}
}

func TestLine_Support(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		dir  Point
		want Point
	}{
		"direction favors p1": {
			p1:   NewPoint(10, 0, 0),
			p2:   NewPoint(0, 10, 0),
			dir:  NewPoint(1, 0, 0),
			want: NewPoint(10, 0, 0),
		},
		"direction favors p2": {
			p1:   NewPoint(10, 0, 0),
			p2:   NewPoint(0, 10, 0),
			dir:  NewPoint(0, 1, 0),
			want: NewPoint(0, 10, 0),
		},
		"equal dot returns p2": {
			p1:   NewPoint(1, 0, 0),
			p2:   NewPoint(0, 1, 0),
			dir:  NewPoint(1, 1, 0),
			want: NewPoint(0, 1, 0),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := NewLine(tc.p1, tc.p2)
			testutils.Equal(t, l.Support(tc.dir), tc.want)
		})
	}
}

func TestLine_ProjectForPoint(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		p    Point
		want Point
	}{
		"project on x axis": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 0, 0),
			p:    NewPoint(5, 5, 0),
			want: NewPoint(5, 0, 0),
		},
		"project on y axis": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(0, 10, 0),
			p:    NewPoint(5, 5, 0),
			want: NewPoint(0, 5, 0),
		},
		"project at start": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 0, 0),
			p:    NewPoint(0, 5, 0),
			want: NewPoint(0, 0, 0),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			l := NewLine(tc.p1, tc.p2)
			testutils.Equal(t, l.ProjectForPoint(tc.p), tc.want)
		})
	}
}
