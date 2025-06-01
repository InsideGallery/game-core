package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestRotatableBox(t *testing.T) {
	bb := NewBox(NewPoint(0, 0, 0), 5, 5, 5)
	rb := NewRotatableBox(bb.Center(), bb, 0, 0, 0)
	testutils.Equal(t, rb.Polyhedron, bb.ToPolyhedron())
	rb = rb.Rotate(bb.Center(), 0, 0, 90)
	testutils.Equal(t, rb.Polyhedron, NewPolyhedron(NewPoint(5, 0, 0), NewPoint(0, 0, 0), NewPoint(0, 5, 0), NewPoint(5, 5, 0), NewPoint(5, 0, 5), NewPoint(0, 0, 5), NewPoint(0, 5, 5), NewPoint(5, 5, 5)))
}
