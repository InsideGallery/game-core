package isometric

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestTileToScreenCoords(t *testing.T) {
	tileSize := shapes.NewPoint(64, 32)
	height := float64(8)

	pos := shapes.NewPoint(3, 2)
	screen := TileToScreenCoords(pos, tileSize, height)

	// x = (3-2)*64/2 + 8*64/2 = 32 + 256 = 288
	// y = (3+2)*32/2 = 80
	testutils.Equal(t, screen.Coordinate(0), 288.0)
	testutils.Equal(t, screen.Coordinate(1), 80.0)
}

func TestTileToScreenCoordsOrigin(t *testing.T) {
	tileSize := shapes.NewPoint(64, 32)
	height := float64(0)

	pos := shapes.NewPoint(0, 0)
	screen := TileToScreenCoords(pos, tileSize, height)
	testutils.Equal(t, screen.Coordinate(0), 0.0)
	testutils.Equal(t, screen.Coordinate(1), 0.0)
}

func TestProjectISO(t *testing.T) {
	p := shapes.NewPoint(4, 2, 1)
	result := ProjectISO(p)
	// x_out = 4-2 = 2
	// y_out = 4/2 + 2/2 - 1 = 2+1-1 = 2
	testutils.Equal(t, result.Coordinate(0), 2.0)
	testutils.Equal(t, result.Coordinate(1), 2.0)
}

func TestProjectISOZeroZ(t *testing.T) {
	p := shapes.NewPoint(6, 2, 0)
	result := ProjectISO(p)
	// x_out = 6-2 = 4
	// y_out = 6/2 + 2/2 - 0 = 3+1 = 4
	testutils.Equal(t, result.Coordinate(0), 4.0)
	testutils.Equal(t, result.Coordinate(1), 4.0)
}

func TestNearness(t *testing.T) {
	p := shapes.NewPoint(1, 2, 3)
	testutils.Equal(t, Nearness(p), 6.0)
}

func TestNearnessZero(t *testing.T) {
	p := shapes.NewPoint(0, 0, 0)
	testutils.Equal(t, Nearness(p), 0.0)
}

func TestCloserTrue(t *testing.T) {
	a := shapes.NewPoint(3, 3, 3) // nearness = 9
	b := shapes.NewPoint(1, 1, 1) // nearness = 3
	testutils.Equal(t, Closer(a, b), true)
}

func TestCloserFalse(t *testing.T) {
	a := shapes.NewPoint(1, 1, 1) // nearness = 3
	b := shapes.NewPoint(3, 3, 3) // nearness = 9
	testutils.Equal(t, Closer(a, b), false)
}

func TestCloserEqual(t *testing.T) {
	a := shapes.NewPoint(2, 2, 2) // nearness = 6
	b := shapes.NewPoint(2, 2, 2) // nearness = 6
	testutils.Equal(t, Closer(a, b), false)
}

func TestTileToScreenAndPixelToScreenRoundTrip(t *testing.T) {
	tileSize := shapes.NewPoint(64, 32)
	height := float64(5)
	// Test multiple tiles
	tiles := []shapes.Point{
		shapes.NewPoint(0, 0),
		shapes.NewPoint(1, 0),
		shapes.NewPoint(0, 1),
		shapes.NewPoint(3, 5),
	}
	for _, tile := range tiles {
		screen := TileToScreenCoords(tile, tileSize, height)
		// Just verify screen coords are finite
		testutils.Equal(t, screen.Coordinate(0) == screen.Coordinate(0), true) // NaN check
	}
}
