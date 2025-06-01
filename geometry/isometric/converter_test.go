package isometric

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestConverter(t *testing.T) {
	tile := shapes.NewPoint(64, 32, 0)
	tilemap := []shapes.Point{
		shapes.NewPoint(0, 0),
		shapes.NewPoint(1, 0),
		shapes.NewPoint(2, 0),
		shapes.NewPoint(0, 1),
		shapes.NewPoint(1, 1),
		shapes.NewPoint(2, 1),
		shapes.NewPoint(0, 2),
		shapes.NewPoint(1, 2),
		shapes.NewPoint(2, 2),
	}
	for _, p := range tilemap {
		iso := FromOrthographic(p, tile)
		ort := ToOrthographic(iso, tile)
		testutils.Equal(t, p, ort)
	}
}

func TestConverterWithStaticSize(t *testing.T) {
	var size float64 = 64
	tilemap := []shapes.Point{
		shapes.NewPoint(0, 0),
		shapes.NewPoint(1, 0),
		shapes.NewPoint(2, 0),
		shapes.NewPoint(0, 1),
		shapes.NewPoint(1, 1),
		shapes.NewPoint(2, 1),
		shapes.NewPoint(0, 2),
		shapes.NewPoint(1, 2),
		shapes.NewPoint(2, 2),
	}
	for _, p := range tilemap {
		iso := FromOrthographicWithStaticSize(p, size)
		ort := ToOrthographicWithStaticSize(iso, size)
		testutils.Equal(t, p, ort)
	}
}

func TestPixelToScreenCoords(t *testing.T) {
	size := shapes.NewPoint(64, 32)
	tilemap := []shapes.Point{
		shapes.NewPoint(253, 39),
	}
	for _, p := range tilemap {
		screen := PixelToScreenCoords(p, size, 8)
		pixel := ScreenToPixelCoords(screen, size, 8)
		testutils.Equal(t, p, pixel)
	}
}
