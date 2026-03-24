package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestPointOperation(t *testing.T) {
	p := NewPoint(0, 0, 0)
	p2 := NewPoint(1, 0, 0)
	p3 := NewPoint(0, 1, 0)

	mp := p.Copy().Add(p2).Add(p3)

	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	testutils.Equal(t, p.Add(p2).Coordinates(), [3]float64{1, 0, 0})
	testutils.Equal(t, p.Add(p3).Coordinates(), [3]float64{0, 1, 0})

	mp = mp.Scale(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{2, 2, 0})
	mp = mp.Divide(NewPoint(2, 2, 0))
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	mp = mp.Divide(NewPoint(0, 0, 0))
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	mp = mp.Increase(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{3, 3, 2})
	mp = mp.Invert()
	testutils.Equal(t, mp.Coordinates(), [3]float64{-3, -3, -2})
	mp = mp.Invert()
	mp = mp.Decrease(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	nmp := mp.Copy()
	nmp = nmp.Scale(20)
	testutils.Equal(t, nmp.Coordinates(), [3]float64{20, 20, 0})
	testutils.Equal(t, nmp.Cross(p2).Coordinates(), [3]float64{0, 0, -20})
	mp = mp.Scale(2)
	mp = mp.Add(NewPoint(0, 0, 1))
	testutils.Equal(t, mp.Coordinates(), [3]float64{2, 2, 1})
	p4 := mp.Copy()
	testutils.Equal(t, p4.Decrease(2).Coordinates(), [3]float64{0, 0, -1})
	testutils.Equal(t, p4.Subtract(NewPoint(1, 0, 0)).Coordinates(), [3]float64{1, 2, 1})
	testutils.Equal(t, p4.Multiply(NewPoint(1, 1, 2)).Coordinates(), [3]float64{2, 2, 2})
	testutils.Equal(t, p4.Divide(NewPoint(2, 1, 0)).Coordinates(), [3]float64{1, 2, 1})
	testutils.Equal(t, p4.Increase(2).Coordinates(), [3]float64{4, 4, 3})
	testutils.Equal(t, p4.Decrease(1).Coordinates(), [3]float64{1, 1, 0})
	testutils.Equal(t, p4.Invert().Coordinates(), [3]float64{-2, -2, -1})
	testutils.Equal(t, p4.Cross(p2).Coordinates(), [3]float64{0, 1, -2})
}

/*
BenchmarkNonModifiablePoint-4                           20000000                60.0 ns/op            32 B/op          1 allocs/op
BenchmarkNonModifiablePointManyUpdates-4                10000000               183 ns/op              32 B/op          1 allocs/op
*/
var (
	globalVector Point
)

func TestPoint_ManhattanDistance(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		want float64
	}{
		"same point": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(0, 0, 0),
			want: 0,
		},
		"one axis": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(5, 0, 0),
			want: 5,
		},
		"two axes": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(3, 4, 0),
			want: 7,
		},
		"three axes": {
			p1:   NewPoint(1, 2, 3),
			p2:   NewPoint(4, 6, 9),
			want: 13,
		},
		"negative": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(-3, -4, 0),
			want: 7,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p1.ManhattanDistance(tc.p2), tc.want)
		})
	}
}

func TestPoint_Distance(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		want float64
	}{
		"same point": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(0, 0, 0),
			want: 0,
		},
		"unit x": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(1, 0, 0),
			want: 1,
		},
		"3-4-5 triangle": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(3, 4, 0),
			want: 5,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p1.Distance(tc.p2), tc.want)
		})
	}
}

func TestPoint_MinDistance(t *testing.T) {
	tests := map[string]struct {
		p    Point
		r    Box
		want float64
	}{
		"point inside box": {
			p:    NewPoint(5, 5, 5),
			r:    NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			want: 0,
		},
		"point outside on x": {
			p:    NewPoint(15, 5, 5),
			r:    NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			want: 25,
		},
		"point below on y": {
			p:    NewPoint(5, -3, 5),
			r:    NewBox(NewPoint(0, 0, 0), 10, 10, 10),
			want: 9,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p.MinDistance(tc.r), tc.want)
		})
	}
}

func TestPoint_MinMaxDistance(t *testing.T) {
	p := NewPoint(0, 0, 0)
	r := NewBox(NewPoint(1, 1, 1), 2, 2, 2)
	result := p.MinMaxDistance(r)
	// Should return a positive number
	testutils.Equal(t, result > 0, true)

	// Point at center of box
	p2 := NewPoint(2, 2, 2)
	result2 := p2.MinMaxDistance(r)
	testutils.Equal(t, result2 >= 0, true)
}

func TestPoint_Coordinate_OutOfBounds(t *testing.T) {
	p := NewPoint(1, 2, 3)
	tests := map[string]struct {
		dim  int
		want float64
	}{
		"dim -1": {dim: -1, want: 0},
		"dim 3":  {dim: 3, want: 0},
		"dim 0":  {dim: 0, want: 1},
		"dim 1":  {dim: 1, want: 2},
		"dim 2":  {dim: 2, want: 3},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, p.Coordinate(tc.dim), tc.want)
		})
	}
}

func TestPoint_GetMinAxis(t *testing.T) {
	tests := map[string]struct {
		p    Point
		want int
	}{
		"x is min": {
			p:    NewPoint(1, 2, 3),
			want: 0,
		},
		"y is min": {
			p:    NewPoint(2, 1, 3),
			want: 1,
		},
		"z is min via x<y path": {
			p:    NewPoint(2, 3, 1),
			want: 2,
		},
		"z is min via y<x path": {
			p:    NewPoint(3, 2, 1),
			want: 2,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p.GetMinAxis(), tc.want)
		})
	}
}

func TestPoint_Abs(t *testing.T) {
	tests := map[string]struct {
		p    Point
		want Point
	}{
		"all positive": {
			p:    NewPoint(1, 2, 3),
			want: NewPoint(1, 2, 3),
		},
		"all negative": {
			p:    NewPoint(-1, -2, -3),
			want: NewPoint(1, 2, 3),
		},
		"mixed": {
			p:    NewPoint(-1, 2, -3),
			want: NewPoint(1, 2, 3),
		},
		"zeros": {
			p:    NewPoint(0, 0, 0),
			want: NewPoint(0, 0, 0),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p.Abs(), tc.want)
		})
	}
}

func TestPoint_Equals(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		want bool
	}{
		"equal": {
			p1:   NewPoint(1, 2, 3),
			p2:   NewPoint(1, 2, 3),
			want: true,
		},
		"not equal": {
			p1:   NewPoint(1, 2, 3),
			p2:   NewPoint(4, 5, 6),
			want: false,
		},
		"zero": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(0, 0, 0),
			want: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p1.Equals(tc.p2), tc.want)
		})
	}
}

func TestPoint_Refract(t *testing.T) {
	tests := map[string]struct {
		p      Point
		normal Point
		eta    float64
	}{
		"normal refraction": {
			p:      NewPoint(1, -1, 0).Normalize(),
			normal: NewPoint(0, 1, 0),
			eta:    0.5,
		},
		"total internal reflection": {
			p:      NewPoint(1, -0.1, 0).Normalize(),
			normal: NewPoint(0, 1, 0),
			eta:    10,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result := tc.p.Refract(tc.normal, tc.eta)
			_ = result.Coordinates() // should not panic
		})
	}
}

func TestPoint_Min(t *testing.T) {
	p1 := NewPoint(1, 5, 3)
	p2 := NewPoint(4, 2, 6)
	result := p1.Min(p2)
	testutils.Equal(t, result, NewPoint(1, 2, 3))
}

func TestPoint_Max(t *testing.T) {
	p1 := NewPoint(1, 5, 3)
	p2 := NewPoint(4, 2, 6)
	result := p1.Max(p2)
	testutils.Equal(t, result, NewPoint(4, 5, 6))
}

func TestPoint_Lerp(t *testing.T) {
	tests := map[string]struct {
		p1   Point
		p2   Point
		f    float64
		want Point
	}{
		"f=0": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 10, 10),
			f:    0,
			want: NewPoint(0, 0, 0),
		},
		"f=1": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 10, 10),
			f:    1,
			want: NewPoint(10, 10, 10),
		},
		"f=0.5": {
			p1:   NewPoint(0, 0, 0),
			p2:   NewPoint(10, 10, 10),
			f:    0.5,
			want: NewPoint(5, 5, 5),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p1.Lerp(tc.p2, tc.f), tc.want)
		})
	}
}

func TestPoint_IsEmpty(t *testing.T) {
	tests := map[string]struct {
		p    Point
		want bool
	}{
		"zero point": {
			p:    NewPoint(0, 0, 0),
			want: true,
		},
		"non-zero x": {
			p:    NewPoint(1, 0, 0),
			want: false,
		},
		"non-zero y": {
			p:    NewPoint(0, 1, 0),
			want: false,
		},
		"non-zero z": {
			p:    NewPoint(0, 0, 1),
			want: false,
		},
		"default point": {
			p:    NewPoint(),
			want: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, tc.p.IsEmpty(), tc.want)
		})
	}
}

func TestPoint_Normalize_ZeroVector(t *testing.T) {
	p := NewPoint(0, 0, 0)
	n := p.Normalize()
	testutils.Equal(t, n, NewPoint(0, 0, 0))
}

func TestPoint_Move(t *testing.T) {
	p := NewPoint(1, 2, 3)
	moved := p.Move(NewPoint(10, 20, 30))
	mp, ok := moved.(Point)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, mp, NewPoint(11, 22, 33))
}

func TestPoint_Bounds_Default(t *testing.T) {
	p := NewPoint(5, 10, 15)
	b := p.Bounds()
	testutils.Equal(t, b.Point1(), NewPoint(5, 10, 15))
	testutils.Equal(t, b.Size(0), 1.0)
	testutils.Equal(t, b.Size(1), 1.0)
	testutils.Equal(t, b.Size(2), 1.0)
}

func TestPoint_Get(t *testing.T) {
	p := NewPoint(1, 2, 3)
	got := p.Get()
	_, ok := got.(Point)
	testutils.Equal(t, ok, true)
}

func TestPoint_Center_And_Point1(t *testing.T) {
	p := NewPoint(5, 10, 15)
	testutils.Equal(t, p.Center(), p)
	testutils.Equal(t, p.Point1(), p)
}

func TestPoint_Support_ReturnsItself(t *testing.T) {
	p := NewPoint(3, 4, 5)
	testutils.Equal(t, p.Support(NewPoint(1, 0, 0)), p)
}

func TestNewPoint_VariousArgs(t *testing.T) {
	tests := map[string]struct {
		args []float64
		want [3]float64
	}{
		"no args": {
			args: nil,
			want: [3]float64{0, 0, 0},
		},
		"one arg": {
			args: []float64{5},
			want: [3]float64{5, 0, 0},
		},
		"two args": {
			args: []float64{5, 10},
			want: [3]float64{5, 10, 0},
		},
		"three args": {
			args: []float64{5, 10, 15},
			want: [3]float64{5, 10, 15},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			p := NewPoint(tc.args...)
			testutils.Equal(t, p.Coordinates(), tc.want)
		})
	}
}

func BenchmarkNonModifiablePoint(b *testing.B) {
	p := NewPoint(1, 1, 1)
	p1 := NewPoint(0, 1, 0)
	for i := 0; i < b.N; i++ {
		p = p.Add(p1).Scale(1.1)
	}
	globalVector = p
}

func BenchmarkNonModifiablePointManyUpdates(b *testing.B) {
	p := NewPoint(1, 1, 1)
	p1 := NewPoint(0, 1, 0)
	for i := 0; i < b.N; i++ {
		p = p.Add(p1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1)
	}
	globalVector = p
}
