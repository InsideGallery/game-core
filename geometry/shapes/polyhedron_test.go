package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
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
