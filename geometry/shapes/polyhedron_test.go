package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestPolyhedron_Bounds(t *testing.T) {
	pl := NewPolyhedron(
		NewPoint(411.333, 508.667),
		NewPoint(506.9997, 471.667),
		NewPoint(481.66630000000004, 401),
		NewPoint(379.9997, 434.667),
	)
	testutils.Equal(t, pl.Bounds(), NewBox(NewPoint(379.9997, 401), 127, 107.66699999999997))
}

func TestPolyhedron_ToLines(t *testing.T) {
	pl := NewPolyhedron(NewPoint(0, 0), NewPoint(10, 0), NewPoint(10, 10), NewPoint(0, 10))
	lines := pl.ToLines()
	testutils.Equal(t, len(lines), 4)
	// First line
	testutils.Equal(t, lines[0].Point1(), NewPoint(0, 0))
	testutils.Equal(t, lines[0].Point2(), NewPoint(10, 0))
	// Last line wraps back to first point
	testutils.Equal(t, lines[3].Point1(), NewPoint(0, 10))
	testutils.Equal(t, lines[3].Point2(), NewPoint(0, 0))
}

func TestPolyhedron_Coordinates(t *testing.T) {
	pl := NewPolyhedron(NewPoint(5, 10), NewPoint(15, 20))
	coords := pl.Coordinates()
	testutils.Equal(t, coords[0], 5.0)
	testutils.Equal(t, coords[1], 10.0)
}

func TestPolyhedron_Move(t *testing.T) {
	pl := NewPolyhedron(NewPoint(0, 0, 0), NewPoint(10, 0, 0), NewPoint(10, 10, 0))
	moved := pl.Move(NewPoint(5, 5, 5)).(Polyhedron)
	testutils.Equal(t, moved.Vector(0), NewPoint(5, 5, 5))
	testutils.Equal(t, moved.Vector(1), NewPoint(15, 5, 5))
	testutils.Equal(t, moved.Vector(2), NewPoint(15, 15, 5))
}

func TestPolyhedron_Point1(t *testing.T) {
	pl := NewPolyhedron(NewPoint(42, 43, 44), NewPoint(10, 20, 30))
	testutils.Equal(t, pl.Point1(), NewPoint(42, 43, 44))
}

func TestPolyhedron_Get(t *testing.T) {
	pl := NewPolyhedron(NewPoint(0, 0), NewPoint(1, 1))
	got := pl.Get()
	_, ok := got.(Polyhedron)
	testutils.Equal(t, ok, true)
}
