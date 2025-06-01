package gjkepa3d

import "github.com/InsideGallery/game-core/geometry/shapes"

const (
	diff              = 0.1
	divide2Multiplier = 0.5
)

// RoundedGJKEPA return rounded collision
func RoundedGJKEPA(shapeA, shapeB shapes.Collide) (bool, shapes.Point) {
	if shapeA.Center().Equals(shapeB.Center()) {
		shapeB = shapeB.Move(shapes.NewPoint(diff, diff, diff)).(shapes.Collide)
	}

	c, p := NewGJKEPA().GJK(shapeA, shapeB, true)
	p = p.Round(diff)

	if c && p.Normal() == 0 {
		return false, p
	}

	return c, p
}

// SolveKinematicBody return rounded collision
func SolveKinematicBody(shapeA, shapeB shapes.Collide,
	force1 shapes.Point, force2 shapes.Point,
	elasticity1, elasticity2 float64,
) (shapes.Collide, shapes.Collide) {
	col, p := RoundedGJKEPA(shapeA, shapeB)
	if col {
		rb := p.Scale(divide2Multiplier)
		p1 := force1.Reflect(p.Normalize()).Scale(elasticity1)
		p2 := force2.Reflect(p.Normalize()).Scale(elasticity2)

		shapeA = shapeA.Move(rb.Invert().Add(p1)).(shapes.Sphere)
		shapeB = shapeB.Move(rb.Add(p2)).(shapes.Sphere)
	}

	return shapeA, shapeB
}
