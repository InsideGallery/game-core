package gjkepa2d

import (
	"math"
	"math/rand"
	"testing"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

func assertPoint(t *testing.T, got, want shapes.Point, tol float64) {
	t.Helper()

	for i := 0; i < 2; i++ {
		if diff := math.Abs(got.Coordinate(i) - want.Coordinate(i)); diff > tol {
			t.Fatalf("coordinate %d: got %v, want %v (diff %g)", i, got, want, diff)
		}
	}
}

func TestMTV_CircleBoxMatchesClosedForm(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)

	centers := []shapes.Point{
		shapes.NewPoint(19.60141, 4.57531), // corner contact
		shapes.NewPoint(19.62, 4.53),       // corner contact
		shapes.NewPoint(19.55, 4.65),       // corner contact
		shapes.NewPoint(19.65, 4.20),       // face contact
	}

	for _, c := range centers {
		want, _, wantHit := circleRectMTV(box, c, 0.25)
		if !wantHit {
			t.Fatalf("reference says no hit for center %v", c)
		}

		mtv, hit := MTV(box, shapes.NewSphere(c, 0.25))
		if !hit {
			t.Fatalf("expected hit for center %v", c)
		}

		assertPoint(t, mtv, want, 1e-12)
	}
}

func TestMTV_OperandOrderFlipsSign(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	circle := shapes.NewSphere(shapes.NewPoint(19.60141, 4.57531), 0.25)

	forward, hitF := MTV(box, circle)
	reverse, hitR := MTV(circle, box)

	if !hitF || !hitR {
		t.Fatalf("expected hit both ways, got %v/%v", hitF, hitR)
	}

	assertPoint(t, reverse, forward.Invert(), 1e-12)
}

func TestMTV_CircleBoxNoOverlap(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(0, 0), 10, 10)

	for _, c := range []shapes.Point{
		shapes.NewPoint(15, 5),      // clear of the east face
		shapes.NewPoint(11, 5),      // exactly touching (dist == r) is not overlap
		shapes.NewPoint(10.8, 10.8), // clear of the corner
		shapes.NewPoint(-3.4, -3.4), // clear of the opposite corner
	} {
		if mtv, hit := MTV(box, shapes.NewSphere(c, 1)); hit {
			t.Fatalf("expected no hit for center %v, got mtv %v", c, mtv)
		}
	}
}

func TestMTV_CircleInsideBoxEjectsThroughNearestFace(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(0, 0), 10, 10)
	r := 1.0

	cases := []struct {
		name   string
		center shapes.Point
		want   shapes.Point // subtracting from the circle must move it out the nearest face
	}{
		{"nearest west", shapes.NewPoint(2, 5), shapes.NewPoint(3, 0)},
		{"nearest east", shapes.NewPoint(8, 5), shapes.NewPoint(-3, 0)},
		{"nearest north", shapes.NewPoint(5, 2), shapes.NewPoint(0, 3)},
		{"nearest south", shapes.NewPoint(5, 8), shapes.NewPoint(0, -3)},
		{"exact center ties break west", shapes.NewPoint(5, 5), shapes.NewPoint(6, 0)},
		{"west/north tie breaks west", shapes.NewPoint(2, 2), shapes.NewPoint(3, 0)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mtv, hit := MTV(box, shapes.NewSphere(tc.center, r))
			if !hit {
				t.Fatal("expected hit")
			}

			assertPoint(t, mtv, tc.want, 1e-12)
		})
	}
}

func TestMTV_CircleCircle(t *testing.T) {
	a := shapes.NewSphere(shapes.NewPoint(100, 100), 10)
	b := shapes.NewSphere(shapes.NewPoint(115, 100), 10)

	// b overlaps a by 5 along +x; subtracting mtv from b must push it to +x
	mtv, hit := MTV(a, b)
	if !hit {
		t.Fatal("expected hit")
	}

	assertPoint(t, mtv, shapes.NewPoint(-5, 0), 1e-12)

	// coincident centers push b along +x, deterministically
	c := shapes.NewSphere(shapes.NewPoint(100, 100), 10)

	mtv, hit = MTV(a, c)
	if !hit {
		t.Fatal("expected hit for coincident centers")
	}

	assertPoint(t, mtv, shapes.NewPoint(-20, 0), 1e-12)

	// separated
	d := shapes.NewSphere(shapes.NewPoint(200, 200), 10)
	if _, hit = MTV(a, d); hit {
		t.Fatal("expected no hit")
	}
}

func TestMTV_CirclePolygon(t *testing.T) {
	square := shapes.NewPolyhedron(
		shapes.NewPoint(0, 0),
		shapes.NewPoint(10, 0),
		shapes.NewPoint(10, 10),
		shapes.NewPoint(0, 10),
	)

	cases := []struct {
		name    string
		center  shapes.Point
		r       float64
		want    shapes.Point
		wantHit bool
	}{
		{"outside radial from east edge", shapes.NewPoint(11, 5), 2, shapes.NewPoint(-1, 0), true},
		// center (11,11) vs corner (10,10): dist=sqrt2, depth=2-sqrt2,
		// per-component mtv = -(depth/sqrt2) = -(sqrt2-1)
		{"outside radial from corner", shapes.NewPoint(11, 11), 2,
			shapes.NewPoint(-(math.Sqrt2 - 1), -(math.Sqrt2 - 1)), true},
		{"flush on east edge ejects along its normal", shapes.NewPoint(10, 5), 1, shapes.NewPoint(-1, 0), true},
		{"inside ejects through nearest edge normal", shapes.NewPoint(9, 5), 1, shapes.NewPoint(-2, 0), true},
		{"separated", shapes.NewPoint(15, 5), 2, shapes.Point{}, false},
		{"touching is not overlap", shapes.NewPoint(12, 5), 2, shapes.Point{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mtv, hit := MTV(square, shapes.NewSphere(tc.center, tc.r))
			if hit != tc.wantHit {
				t.Fatalf("hit: got %v, want %v", hit, tc.wantHit)
			}

			if tc.wantHit {
				assertPoint(t, mtv, tc.want, 1e-12)
			}
		})
	}
}

func TestMTV_DegeneratePolyhedronFallsBackToEPA(t *testing.T) {
	// two-vertex "polygon" has no closed form; the EPA fallback treats it as
	// a segment and must still resolve the contact
	segment := shapes.NewPolyhedron(shapes.NewPoint(0, 0), shapes.NewPoint(10, 0))
	circle := shapes.NewSphere(shapes.NewPoint(5, 0.5), 1)

	mtv, hit := MTV(segment, circle)
	if !hit {
		t.Fatal("expected hit")
	}

	// circle center is 0.5 above the segment: push it up by depth 0.5
	assertPoint(t, mtv, shapes.NewPoint(0, -0.5), 1e-5)
}

func TestMTV_NonSpherePairFallsBackToEPA(t *testing.T) {
	// same square-vs-diamond deep contact as TestEPA_PolygonPolygonDeepContact
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

	mtv, hit := MTV(square, diamond)
	if !hit {
		t.Fatal("expected hit")
	}

	assertPoint(t, mtv, shapes.NewPoint(-3, 0), 1e-5)
}

func TestMTV_PropertyAgreesWithEPA(t *testing.T) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	r := 0.25
	rnd := rand.New(rand.NewSource(7)) //nolint:gosec

	tested := 0
	for tested < 200 {
		c := shapes.NewPoint(
			18.5-r+rnd.Float64()*(1+2*r),
			3.5-r+rnd.Float64()*(1+2*r),
		)

		if _, dist, hit := circleRectMTV(box, c, r); !hit || dist < 0.01 || dist > r*0.99 {
			continue
		}
		tested++

		circle := shapes.NewSphere(c, r)

		exact, hitExact := MTV(box, circle)
		hitEPA, epa := NewGJKEPA().GJK(box, circle, true)

		if !hitExact || !hitEPA {
			t.Fatalf("hit disagreement for center %v: exact=%v epa=%v", c, hitExact, hitEPA)
		}

		// agreement is bounded by EPA's 1e-6 termination tolerance on the
		// corner arcs: magnitude within 1e-5, components within
		// |mtv|*2.8e-3, direction within 0.2 deg (see epa_corner_test.go)
		if diff := math.Abs(exact.Normal() - epa.Normal()); diff > 1e-5 {
			t.Fatalf("magnitude disagreement for center %v: exact %g, epa %g", c, exact.Normal(), epa.Normal())
		}

		assertPoint(t, epa, exact, math.Max(1e-5, exact.Normal()*3e-3))

		if a := angleBetweenDeg(epa, exact); a > 0.2 {
			t.Fatalf("direction disagreement %g deg for center %v", a, c)
		}
	}
}

func BenchmarkMTVCircleBoxAnalytic(b *testing.B) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	circle := shapes.NewSphere(shapes.NewPoint(19.60141, 4.57531), 0.25)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		MTV(box, circle)
	}
}

func BenchmarkMTVCircleBoxEPA(b *testing.B) {
	box := shapes.NewBox(shapes.NewPoint(18.5, 3.5), 1, 1)
	circle := shapes.NewSphere(shapes.NewPoint(19.60141, 4.57531), 0.25)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		NewGJKEPA().GJK(box, circle, true)
	}
}

func BenchmarkMTVCirclePolygonAnalytic(b *testing.B) {
	square := shapes.NewPolyhedron(
		shapes.NewPoint(0, 0),
		shapes.NewPoint(10, 0),
		shapes.NewPoint(10, 10),
		shapes.NewPoint(0, 10),
	)
	circle := shapes.NewSphere(shapes.NewPoint(10.5, 10.5), 1)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		MTV(square, circle)
	}
}
