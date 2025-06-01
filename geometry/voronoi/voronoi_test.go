package voronoi

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestVoronoi(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 600, 480)

	sites := []shapes.Point{
		shapes.NewPoint(110, 20),
		shapes.NewPoint(140, 40),
		shapes.NewPoint(155, 80),
		shapes.NewPoint(350, 120),
		shapes.NewPoint(200, 240),
	}

	v := NewFromPoints(sites, rect)
	v.Generate()
	data := v.ToPolyhedrons()
	testutils.Equal(t, data[0].Point1(), shapes.NewPoint(-128, 246))
}
