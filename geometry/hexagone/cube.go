package hexagone

import "math"

// Cube for hexagon map
type Cube struct {
	X, Y, Z int
}

// NewCube return new cube
func NewCube(x, y, z int) Cube {
	return Cube{
		X: x,
		Y: y,
		Z: z,
	}
}

// ToAxis convert cube to axis
func (c Cube) ToAxis() Axial {
	return Axial{Col: c.X, Row: c.Z}
}

// Distance distance between cubes
func (c Cube) Distance(cube Cube) int {
	return int((math.Abs(float64(c.X)-float64(cube.X)) +
		math.Abs(float64(c.Y)-float64(cube.Y)) +
		math.Abs(float64(c.Z)-float64(cube.Z))) / 2) //nolint:mnd
}
