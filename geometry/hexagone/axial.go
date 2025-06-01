package hexagone

import (
	"math"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

// Axial axial coordinates
type Axial struct {
	Col, Row int
}

// NewAxial return new axial
func NewAxial(col, row int) Axial {
	return Axial{
		Col: col,
		Row: row,
	}
}

// ToCube convert axial to cube
func (a Axial) ToCube() Cube {
	x := a.Col
	z := a.Row
	y := -x - z

	return Cube{X: x, Y: y, Z: z}
}

// ToPosition return position
func (a Axial) ToPosition() shapes.Point {
	x := 1 * (math.Sqrt(3)*float64(a.Col) + math.Sqrt(3)/2*float64(a.Row)) //nolint:mnd
	y := 1 * (3. / 2 * float64(a.Row))

	return shapes.NewPoint(x, y)
}

// Distance return distance between axials
func (a Axial) Distance(axial Axial) int {
	return a.ToCube().Distance(axial.ToCube())
}
