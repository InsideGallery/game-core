package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
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
