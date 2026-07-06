package gjkepa2d

import (
	"math"
	"math/rand"
	"testing"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

// circleRectMTV is the closed-form reference for circle vs axis-aligned rect:
// q = clamp(center into rect); d = center - q; hit iff |d| < r;
// mtv = -(d/|d|)*(r-|d|), so that subtracting mtv from the circle's position
// separates the shapes (same convention as GJK(rect, circle, true)).
func circleRectMTV(box shapes.Box, center shapes.Point, r float64) (shapes.Point, float64, bool) {
	qx := math.Max(box.Point1().Coordinate(0), math.Min(center.Coordinate(0), box.Point2().Coordinate(0)))
	qy := math.Max(box.Point1().Coordinate(1), math.Min(center.Coordinate(1), box.Point2().Coordinate(1)))

	dx := center.Coordinate(0) - qx
	dy := center.Coordinate(1) - qy
	dist := math.Hypot(dx, dy)

	if dist == 0 || dist >= r {
		return shapes.Point{}, dist, false
	}

	scale := (r - dist) / dist

	return shapes.NewPoint(-dx*scale, -dy*scale, 0), dist, true
}

func angleBetweenDeg(a, b shapes.Point) float64 {
	cos := a.Dot(b) / (a.Normal() * b.Normal())
	cos = math.Max(-1, math.Min(1, cos))

	return math.Acos(cos) * 180 / math.Pi
}

func assertMTV(t *testing.T, got, want shapes.Point, componentTol, maxAngleDeg float64) {
	t.Helper()

	for i := 0; i < 2; i++ {
		if diff := math.Abs(got.Coordinate(i) - want.Coordinate(i)); diff > componentTol {
			t.Fatalf("mtv coordinate %d: got %v, want %v (diff %g)", i, got, want, diff)
		}
	}

	if a := angleBetweenDeg(got, want); a > maxAngleDeg {
		t.Fatalf("mtv direction off by %g deg: got %v, want %v", a, got, want)
	}
}

func TestEPA_CircleBoxCornerMatchesClosedForm(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)

	// corner contacts land on the rounded (circular-arc) part of the
	// Minkowski boundary; EPA's 1e-6 termination tolerance on a radius-0.25
	// arc bounds the normal error at sqrt(2*1e-6/0.25) ~ 2.8e-3 rad (0.16
	// deg), i.e. components within |mtv|*2.8e-3 <= ~3.5e-4 while the
	// magnitude stays exact. 5e-4/0.2deg asserts that bound; the pre-fix
	// errors were ~1.5e-2 with direction off by several degrees
	centers := []struct {
		c   shapes.Point
		tol float64
	}{
		{shapes.NewPoint(19.60141, 4.57531), 5e-4}, // corner contact
		{shapes.NewPoint(19.62, 4.53), 5e-4},       // corner contact
		{shapes.NewPoint(19.55, 4.65), 5e-4},       // corner contact
		{shapes.NewPoint(19.65, 4.20), 1e-7},       // face control, must stay exact
	}

	for _, tc := range centers {
		circle := shapes.NewSphere(tc.c, 0.25)

		want, _, wantHit := circleRectMTV(box, tc.c, 0.25)
		if !wantHit {
			t.Fatalf("reference says no hit for center %v", tc.c)
		}

		hit, mtv := NewGJKEPA().GJK(box, circle, true)
		if !hit {
			t.Fatalf("expected hit for center %v", tc.c)
		}

		assertMTV(t, mtv, want, tc.tol, 0.2)

		if diff := math.Abs(mtv.Normal() - want.Normal()); diff > 1e-5 {
			t.Fatalf("mtv magnitude for center %v: got %g, want %g", tc.c, mtv.Normal(), want.Normal())
		}
	}
}

func TestEPA_PolytopeExpands(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	circle := shapes.NewSphere(shapes.NewPoint(19.60141, 4.57531), 0.25)

	g := NewGJKEPA()

	hit, _ := g.GJK(box, circle, true)
	if !hit {
		t.Fatal("expected hit")
	}

	if len(g.vertices) <= 3 {
		t.Fatalf("EPA polytope never expanded: %d vertices", len(g.vertices))
	}
}

func TestEPA_CircleBoxPropertyMatchesClosedForm(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	r := 0.25
	rnd := rand.New(rand.NewSource(42)) //nolint:gosec

	corner := shapes.NewPoint(19.5, 4.5)
	diag := 1 / math.Sqrt2

	// deterministic shallow (depth = 0.01*r) and deep (depth = 0.96*r) contacts
	centers := []shapes.Point{
		corner.Add(shapes.NewPoint(diag, diag).Scale(r * 0.99)), // shallow corner
		corner.Add(shapes.NewPoint(diag, diag).Scale(0.01)),     // deep corner
		shapes.NewPoint(19.5 + r*0.99, 4.0),                     // shallow face
		shapes.NewPoint(19.51, 4.0),                             // deep face
	}

	for len(centers) < 200 {
		c := shapes.NewPoint(
			18.5-r+rnd.Float64()*(1+2*r),
			3.5-r+rnd.Float64()*(1+2*r),
		)

		if _, dist, hit := circleRectMTV(box, c, r); hit && dist >= 0.01 && dist <= r*0.99 {
			centers = append(centers, c)
		}
	}

	for _, c := range centers {
		want, _, wantHit := circleRectMTV(box, c, r)
		if !wantHit {
			t.Fatalf("reference says no hit for center %v", c)
		}

		hit, mtv := NewGJKEPA().GJK(box, shapes.NewSphere(c, r), true)
		if !hit {
			t.Fatalf("expected hit for center %v", c)
		}

		// EPA terminates when the closest edge is within 1e-6 of the true
		// boundary; on a radius-r arc that caps the normal error at
		// sqrt(2*1e-6/r) ~ 2.8e-3 rad (0.16 deg), so component accuracy for
		// corner contacts is |mtv|-proportional while magnitude stays exact
		componentTol := math.Max(1e-5, want.Normal()*3e-3)
		assertMTV(t, mtv, want, componentTol, 0.2)

		if diff := math.Abs(mtv.Normal() - want.Normal()); diff > 1e-5 {
			t.Fatalf("mtv magnitude for center %v: got %g, want %g", c, mtv.Normal(), want.Normal())
		}
	}
}

func TestEPA_PolygonPolygonDeepContact(t *testing.T) {
	// axis-aligned square [0,10]x[0,10] vs a diamond centered at (10,5) with
	// half-diagonal 3. SAT overlaps: x-axis 3, y-axis 6, both diamond
	// normals ~5.66 — so the true MTV moves the square left by 3.
	square := shapes.NewPolyhedron(
		shapes.NewPoint(0, 0),
		shapes.NewPoint(10, 0),
		shapes.NewPoint(10, 10),
		shapes.NewPoint(0, 10),
	)
	diamond := shapes.NewPolyhedron(
		shapes.NewPoint(7, 5),
		shapes.NewPoint(10, 8),
		shapes.NewPoint(13, 5),
		shapes.NewPoint(10, 2),
	)

	hit, mtv := NewGJKEPA().GJK(square, diamond, true)
	if !hit {
		t.Fatal("expected hit")
	}

	assertMTV(t, mtv, shapes.NewPoint(-3, 0, 0), 1e-5, 0.1)
}
