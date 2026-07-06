package gjkepa2d

import (
	"math"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

// MTV returns the minimum translation vector separating a from b
// (subtracting it from b's position separates the shapes) and whether they
// overlap.
//
// Sphere-vs-{Box,Sphere,convex Polyhedron} pairs use exact closed forms —
// no iterations, only the returned point allocates (either operand order;
// the sign flips with it). Any other pair falls back to GJK+EPA, whose
// corner MTVs are only as accurate as EPA's 1e-6 termination tolerance
// allows.
//
// Deterministic tie-breaks: coincident circle centers push b along +x; a
// circle center inside a box exits through the nearest face (first-in-order
// on ties: west, east, north, south); a circle center inside or exactly on
// a polygon boundary exits along the nearest edge's outward normal.
func MTV(a, b shapes.Collide) (shapes.Point, bool) {
	if s, ok := b.(shapes.Sphere); ok {
		if mtv, hit, exact := circleMTV(s.Center(), s.Radius(), a); exact {
			return mtv, hit
		}
	} else if s, ok := a.(shapes.Sphere); ok {
		if mtv, hit, exact := circleMTV(s.Center(), s.Radius(), b); exact {
			if hit {
				mtv = mtv.Invert()
			}

			return mtv, hit
		}
	}

	hit, mtv := NewGJKEPA().GJK(a, b, true)

	return mtv, hit
}

// circleMTV returns the exact MTV for a circle against another shape, and
// whether the shapes overlap. The returned vector separates the shapes when
// subtracted from the CIRCLE's position (MTV flips it when the circle is
// operand a). exact is false when the pair has no closed form and the
// caller must fall back to EPA.
func circleMTV(center shapes.Point, radius float64, other shapes.Collide) (mtv shapes.Point, hit, exact bool) {
	switch o := other.(type) {
	case shapes.Box:
		mtv, hit = circleBoxMTV(center, radius, o)

		return mtv, hit, true
	case shapes.Sphere:
		mtv, hit = circleCircleMTV(center, radius, o.Center(), o.Radius())

		return mtv, hit, true
	case shapes.Polyhedron:
		return circlePolygonMTV(center, radius, o.Vectors())
	default:
		return shapes.Point{}, false, false
	}
}

// circleBoxMTV resolves circle vs axis-aligned rectangle: clamp the center
// to the rect to get the closest point; the MTV is radial from it. A center
// inside the rect is pushed out through the nearest face.
func circleBoxMTV(center shapes.Point, radius float64, box shapes.Box) (shapes.Point, bool) {
	minP, maxP := box.Point1(), box.Point2()
	cx, cy := center.Coordinate(0), center.Coordinate(1)

	qx := math.Max(minP.Coordinate(0), math.Min(cx, maxP.Coordinate(0)))
	qy := math.Max(minP.Coordinate(1), math.Min(cy, maxP.Coordinate(1)))

	dx, dy := cx-qx, cy-qy
	distSq := dx*dx + dy*dy

	if distSq > 0 {
		if distSq >= radius*radius {
			return shapes.Point{}, false
		}

		dist := math.Sqrt(distSq)
		depth := radius - dist

		return shapes.NewPoint(-dx/dist*depth, -dy/dist*depth), true
	}

	// Center inside the rect: exit through the nearest face (deterministic first-in-order tie-break).
	exits := [4]struct{ dist, nx, ny float64 }{
		{cx - minP.Coordinate(0), -1, 0}, // west face
		{maxP.Coordinate(0) - cx, 1, 0},  // east face
		{cy - minP.Coordinate(1), 0, -1}, // north face (-y)
		{maxP.Coordinate(1) - cy, 0, 1},  // south face (+y)
	}

	best := exits[0]
	for _, e := range exits[1:] {
		if e.dist < best.dist {
			best = e
		}
	}

	depth := radius + best.dist

	return shapes.NewPoint(-best.nx*depth, -best.ny*depth), true
}

// circleCircleMTV resolves circle vs circle radially. Coincident centers
// push along +x — an arbitrary but deterministic direction.
func circleCircleMTV(c1 shapes.Point, r1 float64, c2 shapes.Point, r2 float64) (shapes.Point, bool) {
	dx := c1.Coordinate(0) - c2.Coordinate(0)
	dy := c1.Coordinate(1) - c2.Coordinate(1)
	distSq := dx*dx + dy*dy
	sum := r1 + r2

	if distSq >= sum*sum {
		return shapes.Point{}, false
	}

	if distSq == 0 {
		return shapes.NewPoint(-sum, 0), true
	}

	dist := math.Sqrt(distSq)
	depth := sum - dist

	return shapes.NewPoint(-dx/dist*depth, -dy/dist*depth), true
}

// circlePolygonMTV resolves circle vs convex polygon: radial from the
// closest boundary point when the center is outside; out through the
// nearest edge's outward normal when the center is inside or exactly on the
// boundary (where the radial direction is undefined). Degenerate polygons
// (fewer than three vertices, or all edges zero-length) report exact=false
// so the caller falls back to EPA.
func circlePolygonMTV(center shapes.Point, radius float64, verts []shapes.Point) (shapes.Point, bool, bool) {
	if len(verts) < 3 { //nolint:mnd // a polygon needs at least three vertices
		return shapes.Point{}, false, false
	}

	cx, cy := center.Coordinate(0), center.Coordinate(1)

	// Centroid orients each edge's outward normal without winding bookkeeping (convex only).
	var gx, gy float64

	for _, v := range verts {
		gx += v.Coordinate(0)
		gy += v.Coordinate(1)
	}

	gx /= float64(len(verts))
	gy /= float64(len(verts))

	inside := true
	bestDistSq := math.MaxFloat64

	var bestQX, bestQY, bestNX, bestNY float64

	for i := range verts {
		a, b := verts[i], verts[(i+1)%len(verts)]
		ax, ay := a.Coordinate(0), a.Coordinate(1)
		ex, ey := b.Coordinate(0)-ax, b.Coordinate(1)-ay

		lenSq := ex*ex + ey*ey
		if lenSq == 0 {
			continue // degenerate zero-length edge
		}

		// Unit outward normal: perpendicular to the edge, flipped away from the centroid.
		invLen := 1 / math.Sqrt(lenSq)
		nx, ny := ey*invLen, -ex*invLen

		if nx*(gx-ax)+ny*(gy-ay) > 0 {
			nx, ny = -nx, -ny
		}

		if nx*(cx-ax)+ny*(cy-ay) > 0 {
			inside = false // strictly outside this edge's half-plane
		}

		// Closest point on segment ab to the center.
		t := math.Max(0, math.Min(1, ((cx-ax)*ex+(cy-ay)*ey)/lenSq))
		qx, qy := ax+t*ex, ay+t*ey
		dSq := (cx-qx)*(cx-qx) + (cy-qy)*(cy-qy)

		if dSq < bestDistSq {
			bestDistSq, bestQX, bestQY, bestNX, bestNY = dSq, qx, qy, nx, ny
		}
	}

	if bestDistSq == math.MaxFloat64 {
		return shapes.Point{}, false, false // all edges degenerate
	}

	dist := math.Sqrt(bestDistSq)

	if !inside {
		if dist >= radius {
			return shapes.Point{}, false, true
		}

		if dist > 0 {
			depth := radius - dist

			return shapes.NewPoint(-(cx-bestQX)/dist*depth, -(cy-bestQY)/dist*depth), true, true
		}
	}

	// Center inside (or exactly on the boundary, dist == 0): the radial
	// direction is undefined or points inward, so eject through the nearest
	// edge's outward normal.
	depth := radius + dist

	return shapes.NewPoint(-bestNX*depth, -bestNY*depth), true, true
}
