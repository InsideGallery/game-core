package quickhull

import (
	"testing"

	"github.com/FrogoAI/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// ---------- ConvexHull.Triangles ----------

func TestConvexHullTriangles(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(5, 5, 5),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	triangles := hull.Triangles()
	testutils.Equal(t, len(triangles) > 0, true)
	// Each triangle should have 3 points
	for _, tri := range triangles {
		testutils.Equal(t, len(tri), 3)
	}
}

// ---------- ConvexHull with original indices ----------

func TestConvexHullWithOriginalIndices(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(5, 5, 5),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, true, 0)
	testutils.Equal(t, len(hull.Vertices), len(pointCloud))
	testutils.Equal(t, len(hull.Indices) > 0, true)
}

// ---------- ConvexHull CW vs CCW ----------

func TestConvexHullCW(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(5, 5, 5),
	}
	hullCCW := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	hullCW := new(QuickHull).ConvexHull(pointCloud, false, false, 0)
	// Both should have the same number of indices
	testutils.Equal(t, len(hullCCW.Indices), len(hullCW.Indices))
}

// ---------- FastRand tests ----------

func TestFastRandFloat64(t *testing.T) {
	val := FastRandFloat64(1000, -1, 1)
	testutils.Equal(t, val >= -1 && val <= 1, true)
}

func TestFastRandFloat64ZeroPrecision(t *testing.T) {
	val := FastRandFloat64(0, -1, 1)
	testutils.Equal(t, val, 0.0)
}

func TestFastRandFloat64NegativePrecision(t *testing.T) {
	val := FastRandFloat64(-1, -1, 1)
	testutils.Equal(t, val, 0.0)
}

func TestFastRand(t *testing.T) {
	val := FastRand(0, 100)
	testutils.Equal(t, val >= 0 && val < 100, true)
}

func TestFastRandMinGreaterThanMax(t *testing.T) {
	val := FastRand(100, 0)
	testutils.Equal(t, val, 0)
}

// ---------- Empty point cloud ----------

func TestConvexHullEmptyPointCloud(t *testing.T) {
	hull := new(QuickHull).ConvexHull([]shapes.Point{}, true, false, 0)
	testutils.Equal(t, len(hull.Indices), 0)
}

// ---------- HalfEdgeMesh output ----------

func TestConvexHullAsMeshSimple(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(0, 10, 0),
		shapes.NewPoint(0, 0, 10),
		shapes.NewPoint(5, 5, 5),
	}
	mesh := new(QuickHull).ConvexHullAsMesh(pointCloud, 0)
	testutils.Equal(t, len(mesh.Faces) > 0, true)
	testutils.Equal(t, len(mesh.HalfEdges) > 0, true)
	testutils.Equal(t, len(mesh.Vertices) > 0, true)
}

// ---------- Degenerate cases ----------

func TestConvexHullSinglePoint(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(5, 5, 5),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	testutils.Equal(t, len(hull.Indices) > 0, true)
}

func TestConvexHullTwoPoints(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 10, 10),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	testutils.Equal(t, len(hull.Indices) > 0, true)
}

func TestConvexHullThreePoints(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 0),
		shapes.NewPoint(10, 0, 0),
		shapes.NewPoint(0, 10, 0),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	testutils.Equal(t, len(hull.Indices) > 0, true)
}

// ---------- Coplanar points (triggers planar case) ----------

func TestConvexHullCoplanar(t *testing.T) {
	pointCloud := []shapes.Point{
		shapes.NewPoint(0, 0, 5),
		shapes.NewPoint(10, 0, 5),
		shapes.NewPoint(0, 10, 5),
		shapes.NewPoint(10, 10, 5),
		shapes.NewPoint(5, 5, 5),
	}
	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	testutils.Equal(t, len(hull.Indices) > 0, true)
}

// ---------- Plane and Ray utility ----------

func TestPlaneIsPointOnPositiveSide(t *testing.T) {
	n := shapes.NewPoint(0, 0, 1)
	p := newPlane(n, shapes.NewPoint(0, 0, 0))
	// Point above the plane
	testutils.Equal(t, p.isPointOnPositiveSide(shapes.NewPoint(0, 0, 1)), true)
	// Point below the plane
	testutils.Equal(t, p.isPointOnPositiveSide(shapes.NewPoint(0, 0, -1)), false)
}

func TestSignedDistanceToPlane(t *testing.T) {
	n := shapes.NewPoint(0, 0, 1)
	p := newPlane(n, shapes.NewPoint(0, 0, 5))
	d := signedDistanceToPlane(shapes.NewPoint(0, 0, 10), p)
	testutils.Equal(t, d, 5.0)
}

func TestSquaredDistanceBetweenPointAndRay(t *testing.T) {
	r := newRay(shapes.NewPoint(0, 0, 0), shapes.NewPoint(1, 0, 0))
	p := shapes.NewPoint(0, 5, 0)
	d := squaredDistanceBetweenPointAndRay(p, r)
	testutils.Equal(t, d, 25.0)
}

// ---------- HalfEdge disable/isDisabled ----------

func TestHalfEdgeDisable(t *testing.T) {
	he := HalfEdge{EndVertex: 5}
	testutils.Equal(t, he.isDisabled(), false)
	he.disable()
	testutils.Equal(t, he.isDisabled(), true)
}

// ---------- meshBuilderFace disable/isDisabled ----------

func TestMeshBuilderFaceDisable(t *testing.T) {
	f := meshBuilderFace{halfEdgeIndex: 3}
	testutils.Equal(t, f.isDisabled(), false)
	f.disable()
	testutils.Equal(t, f.isDisabled(), true)
}

func TestNewMeshBuilderFace(t *testing.T) {
	f := newMeshBuilderFace()
	testutils.Equal(t, f.isDisabled(), true)
}

// ---------- Large random cloud ----------

func TestConvexHullLargeCloud(t *testing.T) {
	var pointCloud []shapes.Point
	for i := 0; i < 500; i++ {
		pointCloud = append(pointCloud, shapes.NewPoint(
			FastRandFloat64(10000, -10, 10),
			FastRandFloat64(10000, -10, 10),
			FastRandFloat64(10000, -10, 10),
		))
	}
	// Add extremes to ensure bounding box
	pointCloud = append(
		pointCloud,
		shapes.NewPoint(20, 0, 0),
		shapes.NewPoint(-20, 0, 0),
		shapes.NewPoint(0, 20, 0),
		shapes.NewPoint(0, -20, 0),
		shapes.NewPoint(0, 0, 20),
		shapes.NewPoint(0, 0, -20),
	)

	hull := new(QuickHull).ConvexHull(pointCloud, true, false, 0)
	testutils.Equal(t, len(hull.Indices) > 0, true)
	testutils.Equal(t, len(hull.Vertices) > 0, true)
}
