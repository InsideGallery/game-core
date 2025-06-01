package gjkepa3d

import (
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// Constants
const (
	GJKMaxNumIterations = 64
	EPATolerance        = 0.0001
	EPAMaxNumFaces      = 64
	EPAMaxNumLooseEdges = 32
	EPAMaxNumIterations = 64
	dimensions          = 3
	simplexDimensions   = 2
)

// Collider describe collider
type Collider interface {
	Point1() shapes.Point
	Center() shapes.Point
	Support(d shapes.Point) shapes.Point
}

// GJKEPA contains GJK and EPA implementation
type GJKEPA struct{}

// NewGJKEPA return new GJKEPA
func NewGJKEPA() *GJKEPA {
	return &GJKEPA{}
}

// GJK returns true if two colliders are intersecting. Has optional Minimum Translation Vector output param;
// If supplied the EPA will be used to find the vector to separate coll1 from coll2
func (g *GJKEPA) GJK(coll1, coll2 Collider, calculateMTV bool) (bool, shapes.Point) {
	var a, b, c, d shapes.Point                          // Simplex: just a set of points (a is always most recently added)
	searchDir := coll1.Point1().Subtract(coll2.Point1()) // initial search direction between colliders
	// Get initial point for simplex
	c = coll2.Support(searchDir).Subtract(coll1.Support(searchDir.Invert()))
	searchDir = c.Invert() // search in direction of origin
	// Get second point for a line segment simplex
	b = coll2.Support(searchDir).Subtract(coll1.Support(searchDir.Invert()))

	if b.Dot(searchDir) < 0 {
		return false, shapes.NewPoint() // we didn't reach the origin, won't enclose it
	}

	// search perpendicular to line segment towards origin
	searchDir = c.Subtract(b.Invert()).Cross(b.Invert()).Cross(c.Subtract(b))
	if searchDir.Equals(shapes.NewPoint(0, 0, 0)) { // origin is on this line segment
		// Apparently any normal search vector will do?
		searchDir = c.Subtract(b).Cross(shapes.NewPoint(1, 0, 0)) // normal with x-axis
		if searchDir.Equals(shapes.NewPoint(0, 0, 0)) {
			searchDir = c.Subtract(b).Cross(shapes.NewPoint(0, 0, -1)) // normal with z-axis
		}
	}

	simpDim := simplexDimensions // simplex dimension

	for iterations := 0; iterations < GJKMaxNumIterations; iterations++ {
		a = coll2.Support(searchDir).Subtract(coll1.Support(searchDir.Invert()))

		if a.Dot(searchDir) < 0 {
			return false, shapes.NewPoint() // we didn't reach the origin, won't enclose it
		}

		simpDim++
		if simpDim == dimensions {
			_, b, c, d, simpDim, searchDir = g.updateSimplex3(a, b, c, d)
		} else {
			var r bool
			r, a, b, c, d, simpDim, searchDir = g.updateSimplex4(a, b, c, d)

			if r {
				var mtv shapes.Point
				if calculateMTV {
					mtv = g.EPA(a, b, c, d, coll1, coll2)
				}

				return true, mtv
			}
		}
	}

	return false, shapes.NewPoint()
}

// EPA (Expanding Polytope Algorithm) used to find the mtv of two intersecting
// colliders using the final simplex obtained with the GJK algorithm
func (g *GJKEPA) EPA(a, b, c, d shapes.Point, coll1, coll2 Collider) shapes.Point {
	var faces [EPAMaxNumFaces][4]shapes.Point // Array of faces, each with 3 verts and a normal

	// Init with final simplex from GJK
	faces[0][0] = a
	faces[0][1] = b
	faces[0][2] = c
	faces[0][3] = b.Subtract(a).Cross(c.Subtract(a)).Normalize() // ABC

	faces[1][0] = a
	faces[1][1] = c
	faces[1][2] = d
	faces[1][3] = c.Subtract(a).Cross(d.Subtract(a)).Normalize() // ACD

	faces[2][0] = a
	faces[2][1] = d
	faces[2][2] = b
	faces[2][3] = d.Subtract(a).Cross(b.Subtract(a)).Normalize() // ADB

	faces[3][0] = b
	faces[3][1] = d
	faces[3][2] = c
	faces[3][3] = d.Subtract(b).Cross(c.Subtract(b)).Normalize() // BDC
	numFaces := 4
	var closestFace int

	for iterations := 0; iterations < EPAMaxNumIterations; iterations++ {
		// Find face that's closest to origin
		minDist := faces[0][0].Dot(faces[0][3])
		closestFace = 0

		for i := 1; i < numFaces; i++ {
			dist := faces[i][0].Dot(faces[i][3])
			if dist < minDist {
				minDist = dist
				closestFace = i
			}
		}

		// search normal to face that's closest to origin
		searchDir := faces[closestFace][3]
		p := coll2.Support(searchDir).Subtract(coll1.Support(searchDir.Invert()))

		if p.Dot(searchDir)-minDist < EPATolerance {
			// Convergence (new point is not significantly further from origin)
			return faces[closestFace][3].Scale(p.Dot(searchDir))
		}

		var looseEdges [EPAMaxNumLooseEdges][2]shapes.Point
		var numLooseEdges int

		// Find all triangles that are facing p
		for i := 0; i < numFaces; i++ {
			if faces[i][3].Dot(p.Subtract(faces[i][0])) > 0 { // triangle i faces p, remove it
				// Add removed triangle's edges to loose edge list.
				// If it's already there, remove it (both triangles it belonged to are gone)
				for j := 0; j < 3; j++ { // Three edges per face
					currentEdge := [2]shapes.Point{
						faces[i][j],
						faces[i][(j+1)%3],
					}
					var foundEdge bool

					for k := 0; k < numLooseEdges; k++ { // Check if current edge is already in list
						if looseEdges[k][1] == currentEdge[0] && looseEdges[k][0] == currentEdge[1] {
							// Edge is already in the list, remove it
							// THIS ASSUMES EDGE CAN ONLY BE SHARED BY 2 TRIANGLES (which should be true)
							// THIS ALSO ASSUMES SHARED EDGE WILL BE REVERSED IN THE TRIANGLES (which
							// should be true provided every triangle is wound CCW)
							looseEdges[k][0] = looseEdges[numLooseEdges-1][0] // Overwrite current edge
							looseEdges[k][1] = looseEdges[numLooseEdges-1][1] // with last edge in list
							numLooseEdges--
							foundEdge = true

							break
						}
					}

					if !foundEdge { // add current edge to list
						if numLooseEdges >= EPAMaxNumLooseEdges {
							break
						}
						looseEdges[numLooseEdges][0] = currentEdge[0]
						looseEdges[numLooseEdges][1] = currentEdge[1]
						numLooseEdges++
					}
				}

				// Remove triangle i from list
				faces[i][0] = faces[numFaces-1][0]
				faces[i][1] = faces[numFaces-1][1]
				faces[i][2] = faces[numFaces-1][2]
				faces[i][3] = faces[numFaces-1][3]
				numFaces--
				i--
			}
		}

		// Reconstruct polytope with p added
		for i := 0; i < numLooseEdges; i++ {
			if numFaces >= EPAMaxNumFaces {
				break
			}
			faces[numFaces][0] = looseEdges[i][0]
			faces[numFaces][1] = looseEdges[i][1]
			faces[numFaces][2] = p
			faces[numFaces][3] = looseEdges[i][0].Subtract(looseEdges[i][1]).Cross(looseEdges[i][0].Subtract(p)).Normalize()

			// Check for wrong normal to maintain CCW winding
			bias := 0.000001 // in case dot result is only slightly < 0 (because origin is on face)

			if faces[numFaces][0].Dot(faces[numFaces][3])+bias < 0 {
				faces[numFaces][0], faces[numFaces][1] = faces[numFaces][1], faces[numFaces][0]
				faces[numFaces][3] = faces[numFaces][3].Invert()
			}

			numFaces++
		}
	}

	// Return most recent closest point
	return faces[closestFace][3].Scale(faces[closestFace][0].Dot(faces[closestFace][3]))
}

// updateSimplex3 triangle case
func (g *GJKEPA) updateSimplex3(a, b, c, d shapes.Point) (
	shapes.Point, shapes.Point, shapes.Point, shapes.Point, int, shapes.Point,
) {
	/* Required winding order:
	   //  b
	   //  | \
	   //  |   \
	   //  |    a
	   //  |   /
	   //  | /
	   //  c
	*/
	var searchDir shapes.Point
	n := b.Subtract(a).Cross(c.Subtract(a)) // triangle's normal
	AO := a.Invert()                        // direction to origin

	// Determine which feature is closest to origin, make that the new simplex
	simpDim := 2

	if b.Subtract(a).Cross(n).Dot(AO) > 0 { // Closest to edge AB
		c = a
		searchDir = b.Subtract(a).Cross(AO.Cross(b.Subtract(a)))

		return a, b, c, d, simpDim, searchDir
	}

	if b.Cross(c.Subtract(a)).Dot(AO) > 0 { // Closest to edge AC
		b = a
		searchDir = c.Subtract(a).Cross(AO).Cross(c.Subtract(a))

		return a, b, c, d, simpDim, searchDir
	}

	simpDim = 3

	if n.Dot(AO) > 0 { // Above triangle
		d = c
		c = b
		b = a
		searchDir = n

		return a, b, c, d, simpDim, searchDir
	}

	d = b
	b = a
	searchDir = n.Invert()

	return a, b, c, d, simpDim, searchDir
}

// updateSimplex4 tetrahedral case
func (g *GJKEPA) updateSimplex4(a, b, c, d shapes.Point) (
	bool, shapes.Point, shapes.Point, shapes.Point, shapes.Point, int, shapes.Point,
) {
	// a is peak/tip of pyramid, BCD is the base (counterclockwise winding order)

	// Get normals of three new faces
	ABC := b.Subtract(a).Cross(c.Subtract(a))
	ACD := c.Subtract(a).Cross(d.Subtract(a))
	ADB := d.Subtract(a).Cross(b.Subtract(a))

	var searchDir shapes.Point
	AO := a.Invert() // dir to origin
	simpDim := 3     // hoisting this just cause

	// Plane-test origin with 3 faces
	if ABC.Dot(AO) > 0 { // In front of ABC
		d = c
		c = b
		b = a
		searchDir = ABC

		return false, a, b, c, d, simpDim, searchDir
	}

	if ACD.Dot(AO) > 0 { // In front of ACD
		b = a
		searchDir = ACD

		return false, a, b, c, d, simpDim, searchDir
	}

	if ADB.Dot(AO) > 0 { // In front of ADB
		c = d
		d = b
		b = a
		searchDir = ADB

		return false, a, b, c, d, simpDim, searchDir
	}

	// else inside tetrahedron; enclosed!
	return true, a, b, c, d, simpDim, searchDir
}
