package isometric

import "github.com/InsideGallery/game-core/geometry/shapes"

// FromOrthographicWithStaticSize based on pixel calculation
func FromOrthographicWithStaticSize(p shapes.Point, tileSize float64) shapes.Point {
	return FromOrthographic(p, shapes.NewPoint(tileSize, tileSize/2)) //nolint:mnd
}

// ToOrthographicWithStaticSize based on pixel calculation
func ToOrthographicWithStaticSize(p shapes.Point, tileSize float64) shapes.Point {
	return ToOrthographic(p, shapes.NewPoint(tileSize, tileSize/2)) //nolint:mnd
}

// FromOrthographic make isometric point from orthographic
func FromOrthographic(location shapes.Point, tileSize shapes.Point) (screen shapes.Point) {
	return shapes.NewPoint(
		(location.Coordinate(0)-location.Coordinate(1))*(tileSize.Coordinate(0)/2), //nolint:mnd
		(location.Coordinate(0)+location.Coordinate(1))*(tileSize.Coordinate(1)/2), //nolint:mnd
	)
}

// ToOrthographic make orthographic point from isometric
func ToOrthographic(screen shapes.Point, tileSize shapes.Point) (location shapes.Point) {
	return shapes.NewPoint(
		(screen.Coordinate(0)/(tileSize.Coordinate(0)/2)+screen.Coordinate(1)/(tileSize.Coordinate(1)/2))/2, // nolint:mnd
		(screen.Coordinate(1)/(tileSize.Coordinate(1)/2)-screen.Coordinate(0)/(tileSize.Coordinate(0)/2))/2, // nolint:mnd
	)
}

// PixelToScreenCoords pixel to screen
func PixelToScreenCoords(pixel, tileSize shapes.Point, height float64) shapes.Point {
	tileWidth := tileSize.Coordinate(0)
	tileHeight := tileSize.Coordinate(1)
	originX := height * tileWidth / 2 // nolint:mnd
	tileY := pixel.Coordinate(1) / tileHeight
	tileX := pixel.Coordinate(0) / tileHeight

	return shapes.NewPoint(
		(tileX-tileY)*tileWidth/2+originX, //nolint:mnd
		(tileX+tileY)*tileHeight/2,        // nolint:mnd
	)
}

// ScreenToPixelCoords screen to pixel
func ScreenToPixelCoords(screen, tileSize shapes.Point, height float64) shapes.Point {
	x := screen.Coordinate(0)
	y := screen.Coordinate(1)
	tileWidth := tileSize.Coordinate(0)
	tileHeight := tileSize.Coordinate(1)
	x -= height * tileWidth / 2 // nolint:mnd
	tileY := y / tileHeight
	tileX := x / tileWidth

	return shapes.NewPoint(
		(tileY+tileX)*tileHeight,
		(tileY-tileX)*tileHeight,
	)
}

// TileToScreenCoords tile to screen
func TileToScreenCoords(pos, tileSize shapes.Point, height float64) shapes.Point {
	x := pos.Coordinate(0)
	y := pos.Coordinate(1)
	tileWidth := tileSize.Coordinate(0)
	tileHeight := tileSize.Coordinate(1)
	originX := height * tileWidth / 2 // nolint:mnd

	return shapes.NewPoint(
		(x-y)*tileWidth/2+originX,
		(x+y)*tileHeight/2, // nolint:mnd
	)
}

func ProjectISO(p shapes.Point) shapes.Point {
	x, y, z := p.Coordinate(0), p.Coordinate(1), p.Coordinate(2) //nolint:mnd

	return shapes.NewPoint(x-y, (x/2)+(y/2)-z) //nolint:mnd
}

func Nearness(p shapes.Point) float64 {
	return p.Coordinate(0) + p.Coordinate(1) + p.Coordinate(2) //nolint:mnd
}

func Closer(a, b shapes.Point) bool {
	return Nearness(a) > Nearness(b)
}
