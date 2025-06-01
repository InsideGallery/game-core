package shapes

import (
	"github.com/InsideGallery/game-core/mathutils"
)

// RotatableBox describe rotatable box
type RotatableBox struct {
	Polyhedron
	byX, byY, byZ float64
	original      Box
}

// NewRotatableBox return new rotatable box by x,y,z
func NewRotatableBox(c Point, b Box, byX, byY, byZ float64) RotatableBox {
	byX, byY, byZ = DegreesToRadian(byX), DegreesToRadian(byY), DegreesToRadian(byZ)
	ph := b.ToPolyhedron()
	points := make([]Point, ph.Count())

	for i, p := range ph.Vectors() {
		if byX != 0 {
			points[i] = RotateBy(p.Subtract(c), byX, 0).Add(c)
		} else {
			points[i] = p
		}

		if byY != 0 {
			points[i] = RotateBy(p.Subtract(c), byY, 1).Add(c) //nolint:mnd
		} else {
			points[i] = p
		}

		if byZ != 0 {
			points[i] = RotateBy(p.Subtract(c), byZ, 2).Add(c) //nolint:mnd
		} else {
			points[i] = p
		}
	}

	return RotatableBox{
		Polyhedron: NewPolyhedron(points...),
		byX:        byX,
		byY:        byY,
		byZ:        byZ,
		original:   b,
	}
}

// Rotate rotate box
func (b RotatableBox) Rotate(c Point, byX, byY, byZ float64) RotatableBox {
	byX, byY, byZ = DegreesToRadian(byX), DegreesToRadian(byY), DegreesToRadian(byZ)
	ph := b.Polyhedron
	points := make([]Point, ph.Count())

	for i, p := range ph.Vectors() {
		if byX != 0 {
			points[i] = RotateBy(p.Subtract(c), byX, 0).Add(c) // nolint:mnd
		} else {
			points[i] = p
		}

		if byY != 0 {
			points[i] = RotateBy(p.Subtract(c), byY, 1).Add(c) // nolint:mnd
		} else {
			points[i] = p
		}

		if byZ != 0 {
			points[i] = RotateBy(p.Subtract(c), byZ, 2).Add(c) // nolint:mnd
		} else {
			points[i] = p
		}
	}

	return RotatableBox{
		Polyhedron: NewPolyhedron(points...),
		byX:        byX,
		byY:        byY,
		byZ:        byZ,
		original:   b.original,
	}
}

// Get return Spatial
func (b RotatableBox) Get() Spatial {
	return b
}

// Size return size by dimension
func (b RotatableBox) Size(i int) float64 {
	if i > 2 || i < 0 {
		return 0
	}

	return b.original.Point2().Coordinate(i) - b.original.Point1().Coordinate(i)
}

// Sizes return all sizes
func (b RotatableBox) Sizes() [3]float64 {
	var sizes [3]float64
	for i, p := range b.original.Point1().Coordinates() {
		sizes[i] = b.original.Point2().Coordinate(i) - p
	}

	return sizes
}

// Bounds return rectangle of object
func (b RotatableBox) Bounds() Box {
	ps := b.Polyhedron
	l := len(ps.Vectors())
	xs := make([]float64, l)
	ys := make([]float64, l)
	zs := make([]float64, l)

	for i, p := range ps.Vectors() {
		xs[i] = p.Coordinate(0)
		ys[i] = p.Coordinate(1)
		zs[i] = p.Coordinate(2) // nolint:mnd
	}

	sx := mathutils.Min(xs...)
	sy := mathutils.Min(ys...)
	sz := mathutils.Min(zs...)

	bx := mathutils.Max(xs...)
	by := mathutils.Max(ys...)
	bz := mathutils.Max(zs...)

	sizes := make([]float64, 3) // nolint:mnd
	sizes[0] = bx - sx
	sizes[1] = by - sy
	sizes[2] = bz - sz

	return NewBox(NewPoint(sx, sy, sz), sizes...)
}

// Move return new object for given diff
func (b RotatableBox) Move(diff Point) Spatial {
	return NewRotatableBox(b.Center(), b.original.Move(diff).(Box), b.byX, b.byY, b.byZ)
}

// Center return center point
func (b RotatableBox) Center() Point {
	p := make([]float64, 3) // nolint:mnd
	for k, v := range b.Coordinates() {
		p[k] = v + (b.Size(k) / 2) // nolint:mnd
	}

	return NewPoint(p...)
}

// Support return support point for box
func (b RotatableBox) Support(d Point) Point {
	return b.Polyhedron.Support(d)
}
