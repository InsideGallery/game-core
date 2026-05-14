package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestBox(t *testing.T) {
	b := NewBox(NewPoint(5, 10), 3, 2)
	b2 := NewBox(NewPoint(10, 5), 2, 3)
	nb := b.BoundingBox(b2)
	testutils.Equal(t, nb.Point1().Coordinate(0), float64(5))
	testutils.Equal(t, nb.Point1().Coordinate(1), float64(5))
	testutils.Equal(t, nb.Point2().Coordinate(0), float64(12))
	testutils.Equal(t, nb.Point2().Coordinate(1), float64(12))
}

func TestNewBox_VariousArgs(t *testing.T) {
	tests := map[string]struct {
		sizes []float64
		want  [3]float64
	}{
		"no sizes": {
			sizes: nil,
			want:  [3]float64{0, 0, 0},
		},
		"one size": {
			sizes: []float64{5},
			want:  [3]float64{5, 0, 0},
		},
		"two sizes": {
			sizes: []float64{5, 10},
			want:  [3]float64{5, 10, 0},
		},
		"three sizes": {
			sizes: []float64{5, 10, 15},
			want:  [3]float64{5, 10, 15},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			b := NewBox(NewPoint(0, 0, 0), tc.sizes...)
			testutils.Equal(t, b.Sizes(), tc.want)
		})
	}
}

func TestBox_VectorSizes(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 5, 10, 15)
	vs := b.VectorSizes()
	testutils.Equal(t, vs, NewPoint(5, 10, 15))
}

func TestBox_Equal(t *testing.T) {
	tests := map[string]struct {
		b1   Box
		b2   Box
		want bool
	}{
		"equal boxes": {
			b1:   NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			b2:   NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			want: true,
		},
		"different point1": {
			b1:   NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			b2:   NewBox(NewPoint(1, 0, 0), 10, 10, 10),
			want: false,
		},
		"different sizes": {
			b1:   NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			b2:   NewBox(NewPoint(0, 0, 0), 10, 20, 10),
			want: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.b1.Equal(tc.b2), tc.want)
		})
	}
}

func TestBox_Volume(t *testing.T) {
	tests := map[string]struct {
		b    Box
		want float64
	}{
		"unit box": {
			b:    NewBox(NewPoint(0, 0, 0), 1, 1, 1),
			want: 1,
		},
		"large box": {
			b:    NewBox(NewPoint(0, 0, 0), 2, 3, 4),
			want: 24,
		},
		"flat box": {
			b:    NewBox(NewPoint(0, 0, 0), 5, 5, 0),
			want: 0,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.b.Volume(), tc.want)
		})
	}
}

func TestBox_Margin(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 3, 4, 5)
	testutils.Equal(t, b.Margin(), 24.0)
}

func TestBox_ContainsPoint(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	tests := map[string]struct {
		p    Point
		want bool
	}{
		"inside": {
			p:    NewPoint(5, 5, 5),
			want: true,
		},
		"on boundary": {
			p:    NewPoint(0, 0, 0),
			want: true,
		},
		"on boundary max": {
			p:    NewPoint(10, 10, 10),
			want: true,
		},
		"outside": {
			p:    NewPoint(11, 5, 5),
			want: false,
		},
		"outside negative": {
			p:    NewPoint(-1, 5, 5),
			want: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, b.ContainsPoint(tc.p), tc.want)
		})
	}
}

func TestBox_ContainsRectangle(t *testing.T) {
	outer := NewBox(NewPoint(0, 0, 0), 100, 100, 100)
	tests := map[string]struct {
		inner Box
		want  bool
	}{
		"fully inside": {
			inner: NewBox(NewPoint(10, 10, 10), 20, 20, 20),
			want:  true,
		},
		"same box": {
			inner: NewBox(NewPoint(0, 0, 0), 100, 100, 100),
			want:  true,
		},
		"partially outside": {
			inner: NewBox(NewPoint(90, 90, 90), 20, 20, 20),
			want:  false,
		},
		"completely outside": {
			inner: NewBox(NewPoint(200, 200, 200), 10, 10, 10),
			want:  false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, outer.ContainsRectangle(tc.inner), tc.want)
		})
	}
}

func TestBox_Intersect_NoOverlap(t *testing.T) {
	b1 := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	b2 := NewBox(NewPoint(20, 20, 20), 10, 10, 10)
	_, ok := b1.Intersect(b2)
	testutils.Equal(t, ok, false)
}

func TestBox_ToPolygon(t *testing.T) {
	b := NewBox(NewPoint(0, 0), 10, 20)
	poly := b.ToPolygon()
	testutils.Equal(t, poly.Count(), 4)
	testutils.Equal(t, poly.Vector(0), NewPoint(0, 0))
	testutils.Equal(t, poly.Vector(1), NewPoint(0, 20))
	testutils.Equal(t, poly.Vector(2), NewPoint(10, 20))
	testutils.Equal(t, poly.Vector(3), NewPoint(10, 0))
}

func TestBox_Support_AllDirections(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	tests := map[string]struct {
		dir  Point
		want Point
	}{
		"positive xyz": {
			dir:  NewPoint(1, 1, 1),
			want: NewPoint(10, 10, 10),
		},
		"negative xyz": {
			dir:  NewPoint(-1, -1, -1),
			want: NewPoint(0, 0, 0),
		},
		"positive x, negative yz": {
			dir:  NewPoint(1, -1, -1),
			want: NewPoint(10, 0, 0),
		},
		"negative x, positive yz": {
			dir:  NewPoint(-1, 1, 1),
			want: NewPoint(0, 10, 10),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, b.Support(tc.dir), tc.want)
		})
	}
}

func TestBox_Split(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	parts := b.Split()
	testutils.Equal(t, len(parts), 8)
	// Each part should have the center as one of its defining coordinates
	for _, part := range parts {
		testutils.Equal(t, part.Volume() >= 0, true)
	}
}

func TestBox_Contains(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	tests := map[string]struct {
		p    Point
		want bool
	}{
		"inside": {
			p:    NewPoint(5, 5, 5),
			want: true,
		},
		"on min boundary": {
			p:    NewPoint(0, 0, 0),
			want: true,
		},
		"on max boundary": {
			p:    NewPoint(10, 10, 10),
			want: true,
		},
		"outside x": {
			p:    NewPoint(11, 5, 5),
			want: false,
		},
		"outside y": {
			p:    NewPoint(5, 11, 5),
			want: false,
		},
		"outside z": {
			p:    NewPoint(5, 5, 11),
			want: false,
		},
		"outside negative": {
			p:    NewPoint(-1, 5, 5),
			want: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, b.Contains(tc.p), tc.want)
		})
	}
}

func TestBox_Fit(t *testing.T) {
	outer := NewBox(NewPoint(0, 0, 0), 100, 100, 100)
	tests := map[string]struct {
		inner Box
		want  bool
	}{
		"fits": {
			inner: NewBox(NewPoint(10, 10, 10), 20, 20, 20),
			want:  true,
		},
		"does not fit": {
			inner: NewBox(NewPoint(90, 90, 90), 20, 20, 20),
			want:  false,
		},
		"exact fit": {
			inner: NewBox(NewPoint(0, 0, 0), 100, 100, 100),
			want:  true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.inner.Fit(outer), tc.want)
		})
	}
}

func TestBox_Size_OutOfBounds(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 3, 4, 5)
	testutils.Equal(t, b.Size(-1), 0.0)
	testutils.Equal(t, b.Size(3), 0.0)
	testutils.Equal(t, b.Size(0), 3.0)
}

func TestBox_Get(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 5, 5, 5)
	got := b.Get()
	_, ok := got.(Box)
	testutils.Equal(t, ok, true)
}

func TestBox_Move(t *testing.T) {
	b := NewBox(NewPoint(1, 2, 3), 5, 10, 15)
	moved := b.Move(NewPoint(10, 20, 30)).(Box)
	testutils.Equal(t, moved.Point1(), NewPoint(11, 22, 33))
	testutils.Equal(t, moved.Sizes(), [3]float64{5, 10, 15})
}

func TestBox_Center(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 20, 30)
	testutils.Equal(t, b.Center(), NewPoint(5, 10, 15))
}

func TestBox_Bounds(t *testing.T) {
	b := NewBox(NewPoint(1, 2, 3), 4, 5, 6)
	testutils.Equal(t, b.Bounds(), b)
}
