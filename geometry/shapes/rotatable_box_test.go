package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestRotatableBox(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 5, 5, 5)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	testutils.Equal(t, rb.Polyhedron, bb.ToPolyhedron())
	rb = rb.Rotate(bb.Center(), 0, 0, 90)
	testutils.Equal(t, rb.Polyhedron, NewPolyhedron(NewPoint(5, 0, 0), NewPoint(0, 0, 0), NewPoint(0, 5, 0), NewPoint(5, 5, 0), NewPoint(5, 0, 5), NewPoint(0, 0, 5), NewPoint(0, 5, 5), NewPoint(5, 5, 5)))
}

func TestRotatableBox_Get(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 5, 5, 5)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	got := rb.Get()
	_, ok := got.(RotatableBox)
	testutils.Equal(t, ok, true)
}

func TestRotatableBox_Size(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 5, 10, 15)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)

	tests := map[string]struct {
		dim  int
		want float64
	}{
		"dim 0":                {dim: 0, want: 5},
		"dim 1":                {dim: 1, want: 10},
		"dim 2":                {dim: 2, want: 15},
		"dim -1 out of bounds": {dim: -1, want: 0},
		"dim 3 out of bounds":  {dim: 3, want: 0},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, rb.Size(tc.dim), tc.want)
		})
	}
}

func TestRotatableBox_Sizes(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 5, 10, 15)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	testutils.Equal(t, rb.Sizes(), [3]float64{5, 10, 15})
}

func TestRotatableBox_Bounds(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	bounds := rb.Bounds()
	// For unrotated box, bounds should match original
	testutils.Equal(t, bounds.Point1(), NewPoint(0, 0, 0))
	testutils.Equal(t, bounds.Sizes(), [3]float64{10, 10, 10})
}

func TestRotatableBox_Move(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	moved := rb.Move(NewPoint(5, 5, 5))
	// Should not panic
	_ = moved.Bounds()
}

func TestRotatableBox_Center(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	c := rb.Center()
	testutils.Equal(t, c, NewPoint(5, 5, 5))
}

func TestRotatableBox_Support(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	// Use a direction that uniquely selects one vertex
	p := rb.Support(NewPoint(1, 1, 1))
	testutils.Equal(t, p, NewPoint(10, 10, 10))
}

func TestRotatableBox_RotateByX(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 90, 0, 0)
	// Should not panic, and should produce valid bounds
	bounds := rb.Bounds()
	_ = bounds.Sizes()
}

func TestRotatableBox_RotateByY(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 90, 0)
	bounds := rb.Bounds()
	_ = bounds.Sizes()
}

func TestRotatableBox_RotateMultipleAxes(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 45, 45, 45)
	bounds := rb.Bounds()
	_ = bounds.Sizes()
}

func TestRotatableBox_Rotate_ByXAxis(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	rotated := rb.Rotate(bb.Center(), 90, 0, 0)
	bounds := rotated.Bounds()
	_ = bounds.Sizes()
}

func TestRotatableBox_Rotate_ByYAxis(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	rotated := rb.Rotate(bb.Center(), 0, 90, 0)
	bounds := rotated.Bounds()
	_ = bounds.Sizes()
}
