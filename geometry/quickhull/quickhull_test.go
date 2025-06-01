package quickhull

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// Simple 2D test (square, all points on a plane)
func TestConvexHull2DSquare(t *testing.T) {
	// Construct a square as 'point cloud' that looks roughly like this.
	// All points are on a single plane (Z=1).
	// E is on the border but should not be part of the Convex Hull.
	// F should be inside the Convex Hull.
	//
	// C - - - - - - D
	// |             |
	// |             |
	// E      F      |
	// |             |
	// |             |
	// A - - - - - - B
	//
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 1),
		shapes.NewPoint(0, 5, 1),
		shapes.NewPoint(0, 10, 1),
		shapes.NewPoint(10, 0, 1),
		shapes.NewPoint(10, 10, 1),
		shapes.NewPoint(5, 5, 1), // This point is inside the hull
	}

	expectedHull := []shapes.Point{
		shapes.NewPoint(10, 0, 1),
		shapes.NewPoint(0, 0, 1),
		shapes.NewPoint(0, 10, 1),
		shapes.NewPoint(10, 10, 1),
	}

	testutils.Equal(t, expectedHull, convexHull(pointCloud).Vertices)
}

// Simple 2D test (triangle, all points on a plane)
func TestConvexHull2DTriangle(t *testing.T) {
	// Construct a triangular as 'point cloud' that looks roughly like this.
	// All points are on a single plane (Z=1).
	// D should be inside the Convex Hull.
	//
	//         C
	//       /   \
	//     /   D   \
	//   /           \
	// A - - - - - - - B
	//
	pointCloud := []shapes.Point{
		shapes.NewPoint(1, 2, 1),
		shapes.NewPoint(4, 7, 1),
		shapes.NewPoint(7, 2, 1),
		shapes.NewPoint(4, 4, 1), // This point is inside the hull
	}

	expectedHull := []shapes.Point{
		shapes.NewPoint(7, 2, 1),
		shapes.NewPoint(4, 7, 1),
		shapes.NewPoint(1, 2, 1),
	}

	testutils.Equal(t, expectedHull, convexHull(pointCloud).Vertices)
}

// Simple 3D test (one point in a box)
func TestConvexHull3D(t *testing.T) {
	// Construct a point cloud that looks roughly like this.
	// A should be inside the Convex Hull.
	/*
		@ + + + + + + + + + + + @
		+\                      +\
		+ \                     + \
		+  \                    +  \
		+   \                   +   \
		+    @ + + + + + + + + +++ + @
		+    +                  +    +
		+    +                  +    +
		+    +                  +    +
		+    +         A        +    +
		+    +                  +    +
		+    +                  +    +
		@ + +++ + + + + + + + + @    +
		 \   +                   \   +
		  \  +                    \  +
		   \ +                     \ +
		    \+                      \+
		     @ + + + + + + + + + + + @
	*/
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 10, 10),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(10, 0, 10),
		shapes.NewPoint(10, 10, 0),
		shapes.NewPoint(10, 10, 10),
		shapes.NewPoint(5, 5, 5),
	}

	actual := convexHull(pointCloud).Vertices

	expected := []shapes.Point{
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(10, 10, 0),
		shapes.NewPoint(10, 10, 10),
		shapes.NewPoint(0, 10, 10),
		shapes.NewPoint(10, 0, 10),
	}

	testutils.Equal(t, 8, len(actual))
	testutils.Equal(t, expected, convexHull(pointCloud).Vertices)
}

func TestHalfEdgeOutput(t *testing.T) {
	var pointCloud []shapes.Point

	for i := 0; i < 1000; i++ {
		pointCloud = append(pointCloud, shapes.NewPoint(
			FastRandFloat64(10000, -1, 1),
			FastRandFloat64(10000, -1, 1),
			FastRandFloat64(10000, -1, 1),
			FastRandFloat64(10000, -1, 1),
		))
	}

	twoOrNegativeTwo := func(i, and int) float64 {
		if i&and > 0 {
			return -2
		}
		return 2
	}

	for i := 0; i < 8; i++ {
		pointCloud = append(pointCloud, shapes.NewPoint(
			twoOrNegativeTwo(i, 1),
			twoOrNegativeTwo(i, 2),
			twoOrNegativeTwo(i, 4),
		))
	}

	mesh := new(QuickHull).ConvexHullAsMesh(pointCloud, 0)

	testutils.Equal(t, 12, len(mesh.Faces))
	testutils.Equal(t, 36, len(mesh.HalfEdges))
	testutils.Equal(t, 8, len(mesh.Vertices))
}

func TestPlanes(t *testing.T) {
	m := shapes.NewPoint(1, 0, 0)
	n := shapes.NewPoint(2, 0, 0)
	p := newPlane(m, n)

	dist := signedDistanceToPlane(shapes.NewPoint(3, 0, 0), p)
	testutils.Equal(t, 1.0, dist)

	dist = signedDistanceToPlane(shapes.NewPoint(1, 0, 0), p)
	testutils.Equal(t, -1.0, dist)

	m = shapes.NewPoint(2, 0, 0)
	p = newPlane(m, n)

	dist = signedDistanceToPlane(shapes.NewPoint(6, 0, 0), p)

	testutils.Equal(t, 8.0, dist)
}

func convexHull(pointCloud []shapes.Point) ConvexHull {
	return new(QuickHull).ConvexHull(pointCloud, true, false, 0)
}
