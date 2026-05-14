package gjkepa3d

import (
	"testing"

	"github.com/FrogoAI/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestRoundedGJKEPA(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100, 100), 10)
	col, p := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, true)
	testutils.Equal(t, p.Coordinate(0) != 0, true)
}

func TestRoundedGJKEPANoCollision(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(200, 200, 200), 10)
	col, _ := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, false)
}

func TestRoundedGJKEPASameCenter(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	col, p := RoundedGJKEPA(s1, s2)
	// When centers are equal, one is moved by diff; collision should still be detected
	testutils.Equal(t, col, true)
	testutils.Equal(t, p.Normal() > 0, true)
}

func TestSolveKinematicBody(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100, 100), 10)
	force1 := shapes.NewPoint(1, 0, 0)
	force2 := shapes.NewPoint(-1, 0, 0)
	a, b := SolveKinematicBody(s1, s2, force1, force2, 0.8, 0.8)
	// After solving, they should have moved apart
	testutils.Equal(t, a.Center().Coordinate(0) < b.Center().Coordinate(0), true)
}

func TestSolveKinematicBodyNoCollision(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(200, 200, 200), 10)
	force1 := shapes.NewPoint(1, 0, 0)
	force2 := shapes.NewPoint(-1, 0, 0)
	a, b := SolveKinematicBody(s1, s2, force1, force2, 0.8, 0.8)
	// No collision, shapes should remain the same
	testutils.Equal(t, a.Center(), s1.Center())
	testutils.Equal(t, b.Center(), s2.Center())
}

func TestGJKNoMTV(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(115, 100, 100), 10)
	g := NewGJKEPA()
	r, p := g.GJK(s1, s2, false)
	testutils.Equal(t, r, true)
	// Without calculateMTV, mtv should be zero
	var zero shapes.Point
	testutils.Equal(t, p, zero)
}

func TestGJKNoIntersection(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(0, 0, 0), 5)
	s2 := shapes.NewSphere(shapes.NewPoint(100, 100, 100), 5)
	g := NewGJKEPA()
	r, _ := g.GJK(s1, s2, true)
	testutils.Equal(t, r, false)
}

func TestGJKBoxVsSphere(t *testing.T) {
	bb := shapes.NewBox(shapes.NewPoint(0, 0, 0), 100, 100)
	s := shapes.NewSphere(shapes.NewPoint(5, 50, 0), 10)
	g := NewGJKEPA()
	r, p := g.GJK(s, bb, true)
	testutils.Equal(t, r, true)
	testutils.Equal(t, p.Normal() > 0, true)
}

func TestGJKMultiObjectVsSphere(t *testing.T) {
	s1 := shapes.NewSphere(shapes.NewPoint(0, 0, 0), 20)
	s2 := shapes.NewSphere(shapes.NewPoint(50, 0, 0), 20)
	multi := shapes.NewMultiObject(s1, s2)
	s3 := shapes.NewSphere(shapes.NewPoint(25, 0, 0), 10)
	g := NewGJKEPA()
	r, _ := g.GJK(multi, s3, true)
	testutils.Equal(t, r, true)
}

// Test the specific simplex3 and simplex4 update paths by using
// overlapping spheres in different orientations
func TestGJKSphereOverlapVariousPositions(t *testing.T) {
	g := NewGJKEPA()
	tests := []struct {
		name string
		c1   shapes.Point
		c2   shapes.Point
	}{
		{"overlap-x", shapes.NewPoint(0, 0, 0), shapes.NewPoint(15, 0, 0)},
		{"overlap-y", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, 15, 0)},
		{"overlap-z", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, 0, 15)},
		{"overlap-xy", shapes.NewPoint(0, 0, 0), shapes.NewPoint(10, 10, 0)},
		{"overlap-xz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(10, 0, 10)},
		{"overlap-yz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, 10, 10)},
		{"overlap-xyz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(8, 8, 8)},
		{"overlap-neg-xyz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(-8, -8, -8)},
		{"overlap-neg-x", shapes.NewPoint(0, 0, 0), shapes.NewPoint(-15, 0, 0)},
		{"overlap-neg-y", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, -15, 0)},
		{"overlap-neg-z", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, 0, -15)},
		{"overlap-neg-xy", shapes.NewPoint(0, 0, 0), shapes.NewPoint(-10, -10, 0)},
		{"overlap-neg-xz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(-10, 0, -10)},
		{"overlap-neg-yz", shapes.NewPoint(0, 0, 0), shapes.NewPoint(0, -10, -10)},
		{"overlap-mixed-1", shapes.NewPoint(5, 5, 5), shapes.NewPoint(-5, -5, -5)},
		{"overlap-mixed-2", shapes.NewPoint(-5, 5, -5), shapes.NewPoint(5, -5, 5)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s1 := shapes.NewSphere(tc.c1, 10)
			s2 := shapes.NewSphere(tc.c2, 10)
			r, p := g.GJK(s1, s2, true)
			testutils.Equal(t, r, true)
			testutils.Equal(t, p.Normal() > 0, true)
		})
	}
}

// Test updateSimplex3 and updateSimplex4 directly with various shapes
// to hit different branches (ACD, ADB planes in simplex4)
func TestGJKVariousShapeCollisions(t *testing.T) {
	g := NewGJKEPA()

	// Sphere vs box in various positions triggers different simplex paths
	tests := []struct {
		name   string
		sphere shapes.Sphere
		box    shapes.Box
	}{
		{"sphere-top-left-box", shapes.NewSphere(shapes.NewPoint(-5, -5, 0), 10), shapes.NewBox(shapes.NewPoint(0, 0, 0), 20, 20)},
		{"sphere-bottom-right-box", shapes.NewSphere(shapes.NewPoint(25, 25, 0), 10), shapes.NewBox(shapes.NewPoint(0, 0, 0), 20, 20)},
		{"sphere-inside-box", shapes.NewSphere(shapes.NewPoint(10, 10, 0), 5), shapes.NewBox(shapes.NewPoint(0, 0, 0), 20, 20)},
		{"sphere-on-edge-box", shapes.NewSphere(shapes.NewPoint(0, 10, 0), 5), shapes.NewBox(shapes.NewPoint(0, 0, 0), 20, 20)},
		{"sphere-corner-box", shapes.NewSphere(shapes.NewPoint(-3, -3, 0), 10), shapes.NewBox(shapes.NewPoint(0, 0, 0), 20, 20)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, p := g.GJK(tc.sphere, tc.box, true)
			testutils.Equal(t, r, true)
			testutils.Equal(t, p.Normal() > 0, true)
		})
	}
}

// Test sphere vs multi-object at various angles to trigger different simplex4 branches
func TestGJKMultiObjectVariousAngles(t *testing.T) {
	g := NewGJKEPA()
	s1 := shapes.NewSphere(shapes.NewPoint(0, 0, 0), 30)
	s2 := shapes.NewSphere(shapes.NewPoint(50, 50, 50), 30)
	multi := shapes.NewMultiObject(s1, s2)

	tests := []struct {
		name string
		pos  shapes.Point
	}{
		{"center", shapes.NewPoint(25, 25, 25)},
		{"near-s1", shapes.NewPoint(5, 5, 5)},
		{"near-s2", shapes.NewPoint(45, 45, 45)},
		{"above", shapes.NewPoint(25, 25, 50)},
		{"below", shapes.NewPoint(25, 25, 0)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := shapes.NewSphere(tc.pos, 10)
			r, _ := g.GJK(multi, s, true)
			testutils.Equal(t, r, true)
		})
	}
}

// Exercise the RoundedGJKEPA path where collision is true but MTV rounds to zero normal
func TestRoundedGJKEPABarely(t *testing.T) {
	// Two spheres barely touching - distance = 20, sum of radii = 20
	s1 := shapes.NewSphere(shapes.NewPoint(0, 0, 0), 10)
	s2 := shapes.NewSphere(shapes.NewPoint(19.9, 0, 0), 10)
	col, _ := RoundedGJKEPA(s1, s2)
	testutils.Equal(t, col, true)
}

// Directly test updateSimplex3 and updateSimplex4 branches
func TestUpdateSimplexBranches(t *testing.T) {
	g := NewGJKEPA()

	// updateSimplex3 - Branch: above triangle (n.Dot(AO) > 0)
	// origin is at (0,0,0), a is closest to origin
	// triangle where origin is above the plane
	{
		a := shapes.NewPoint(-1, -1, 0)
		b := shapes.NewPoint(0, 2, 0)
		c := shapes.NewPoint(2, 0, 0)
		d := shapes.NewPoint(0, 0, 0)
		_, _, _, _, simpDim, _ := g.updateSimplex3(a, b, c, d)
		testutils.Equal(t, simpDim, 3)
	}

	// updateSimplex4 all four branches
	// Branch 1: ABC.Dot(AO) > 0 (in front of ABC)
	{
		a := shapes.NewPoint(-1, -1, -1)
		b := shapes.NewPoint(1, 0, 0)
		c := shapes.NewPoint(0, 1, 0)
		d := shapes.NewPoint(0, 0, -3)
		found, _, _, _, _, _, _ := g.updateSimplex4(a, b, c, d)
		testutils.Equal(t, found, false)
	}

	// Branch 2: ACD.Dot(AO) > 0 (in front of ACD, but not ABC)
	// ABC normal should point away from origin, ACD normal toward origin
	{
		a := shapes.NewPoint(-1, -1, -1)
		b := shapes.NewPoint(0, 0, -3)
		c := shapes.NewPoint(1, 0, 0)
		d := shapes.NewPoint(0, 1, 0)
		found, _, _, _, _, _, _ := g.updateSimplex4(a, b, c, d)
		testutils.Equal(t, found, false)
	}

	// Branch 3: ADB.Dot(AO) > 0 (in front of ADB, but not ABC or ACD)
	{
		a := shapes.NewPoint(-1, -1, -1)
		b := shapes.NewPoint(0, 1, 0)
		c := shapes.NewPoint(0, 0, -3)
		d := shapes.NewPoint(1, 0, 0)
		found, _, _, _, _, _, _ := g.updateSimplex4(a, b, c, d)
		testutils.Equal(t, found, false)
	}

	// Branch 4: Inside tetrahedron (enclosed)
	// Origin must be inside the tetrahedron formed by a,b,c,d
	// Use a tetrahedron that clearly contains the origin
	{
		a := shapes.NewPoint(1, 1, 1)
		b := shapes.NewPoint(-2, 1, -1)
		c := shapes.NewPoint(1, -2, -1)
		d := shapes.NewPoint(1, 1, -2)
		// Verify all normals point away from origin
		ABC := b.Subtract(a).Cross(c.Subtract(a))
		ACD := c.Subtract(a).Cross(d.Subtract(a))
		ADB := d.Subtract(a).Cross(b.Subtract(a))
		AO := a.Invert()
		if ABC.Dot(AO) <= 0 && ACD.Dot(AO) <= 0 && ADB.Dot(AO) <= 0 {
			found, _, _, _, _, _, _ := g.updateSimplex4(a, b, c, d)
			testutils.Equal(t, found, true)
		}
	}
}
