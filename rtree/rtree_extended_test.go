package rtree

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// ---------------------------------------------------------------------------
// EntityThing – implements Entity interface for entity-based deletion tests
// ---------------------------------------------------------------------------

type EntityThing struct {
	shapes.Sphere
	ID uint32
}

func NewEntityThing(where shapes.Sphere, id uint32) *EntityThing {
	return &EntityThing{
		Sphere: where,
		ID:     id,
	}
}

func (e *EntityThing) GetID() uint32 {
	return e.ID
}

func (e *EntityThing) SetNewSpatial(s shapes.Spatial) {
	e.Sphere = s.(shapes.Sphere)
}

func (e *EntityThing) Bounds() shapes.Box {
	return e.Sphere.Bounds()
}

func (e *EntityThing) UpdateSpatial(s shapes.Spatial) {
	e.Sphere = s.(shapes.Sphere)
}

// ---------------------------------------------------------------------------
// NewRTree parameter adjustments
// ---------------------------------------------------------------------------

func TestNewRTreeMinLessThan2(t *testing.T) {
	rt := NewRTree(1, 10)
	testutils.Equal(t, rt.MinChildren, 2)
	testutils.Equal(t, rt.MaxChildren, 10)
}

func TestNewRTreeMinZero(t *testing.T) {
	rt := NewRTree(0, 10)
	testutils.Equal(t, rt.MinChildren, 2)
}

func TestNewRTreeMinGreaterThanHalfMax(t *testing.T) {
	// min=10, max=15 -> 10 > 15/2=7, so max becomes 10*2=20
	rt := NewRTree(10, 15)
	testutils.Equal(t, rt.MinChildren, 10)
	testutils.Equal(t, rt.MaxChildren, 20)
}

func TestNewRTreeValidParams(t *testing.T) {
	rt := NewRTree(5, 20)
	testutils.Equal(t, rt.MinChildren, 5)
	testutils.Equal(t, rt.MaxChildren, 20)
}

// ---------------------------------------------------------------------------
// Size and Insert
// ---------------------------------------------------------------------------

func TestInsertAndSize(t *testing.T) {
	rt := NewRTree(2, 5)
	testutils.Equal(t, rt.Size(), 0)

	rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10)))
	testutils.Equal(t, rt.Size(), 1)

	rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10)))
	testutils.Equal(t, rt.Size(), 2)

	rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(20, 20), 10)))
	testutils.Equal(t, rt.Size(), 3)
}

func TestInsertManyCausesSplits(t *testing.T) {
	rt := NewRTree(2, 4)
	for i := 0; i < 50; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*10, float64(i)*10), 5)))
	}
	testutils.Equal(t, rt.Size(), 50)
	if rt.Depth() < 2 {
		t.Fatalf("expected depth >= 2 after many inserts, got %d", rt.Depth())
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDeleteExisting(t *testing.T) {
	rt := NewRTree(2, 5)
	obj := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 10))
	rt.Insert(obj)
	testutils.Equal(t, rt.Size(), 1)

	ok := rt.Delete(obj)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, rt.Size(), 0)
}

func TestDeleteNonExisting(t *testing.T) {
	rt := NewRTree(2, 5)
	obj1 := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 10))
	obj2 := NewThing(shapes.NewSphere(shapes.NewPoint(500, 500), 10))
	rt.Insert(obj1)

	ok := rt.Delete(obj2)
	testutils.Equal(t, ok, false)
	testutils.Equal(t, rt.Size(), 1)
}

func TestDeleteFromEmpty(t *testing.T) {
	rt := NewRTree(2, 5)
	obj := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 10))
	ok := rt.Delete(obj)
	testutils.Equal(t, ok, false)
}

func TestDeleteAllObjects(t *testing.T) {
	rt := NewRTree(2, 5)
	things := make([]*Thing, 10)
	for i := 0; i < 10; i++ {
		things[i] = NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*20, float64(i)*20), 5))
		rt.Insert(things[i])
	}
	testutils.Equal(t, rt.Size(), 10)

	for _, th := range things {
		ok := rt.Delete(th)
		testutils.Equal(t, ok, true)
	}
	testutils.Equal(t, rt.Size(), 0)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate(t *testing.T) {
	rt := NewRTree(2, 5)
	obj := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 10))
	rt.Insert(obj)
	testutils.Equal(t, rt.Size(), 1)

	// Modify the object and update.
	obj.Sphere = shapes.NewSphere(shapes.NewPoint(50, 50), 10)
	rt.Update(obj)
	testutils.Equal(t, rt.Size(), 1)

	// Verify it can be found at the new location.
	results := rt.SearchIntersect(obj.Bounds(), nil)
	testutils.Equal(t, len(results), 1)
	testutils.Equal(t, results[0], obj)
}

// ---------------------------------------------------------------------------
// SearchIntersect
// ---------------------------------------------------------------------------

func TestSearchIntersectReturnsMatching(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(100, 100), 10))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 5))
	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)

	// Search around origin – should find t1 and t3, not t2.
	searchBox := shapes.NewBox(shapes.NewPoint(-5, -5), 20, 20)
	results := rt.SearchIntersect(searchBox, nil)
	testutils.Equal(t, len(results), 2)
}

func TestSearchIntersectWithFilter(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(5, 5), 5))
	rt.Insert(t1)
	rt.Insert(t2)

	searchBox := shapes.NewBox(shapes.NewPoint(-15, -15), 30, 30)
	// Filter out t1.
	results := rt.SearchIntersect(searchBox, func(s shapes.Spatial) bool {
		return s == t1
	})
	testutils.Equal(t, len(results), 1)
	testutils.Equal(t, results[0], t2)
}

func TestSearchIntersectNoResults(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 5))
	rt.Insert(t1)

	searchBox := shapes.NewBox(shapes.NewPoint(1000, 1000), 10, 10)
	results := rt.SearchIntersect(searchBox, nil)
	testutils.Equal(t, len(results), 0)
}

// ---------------------------------------------------------------------------
// Collision
// ---------------------------------------------------------------------------

func TestCollisionReturnsIntersecting(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 20))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 20))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(500, 500), 5))
	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)

	results := rt.Collision(t1, nil)
	// t1 and t2 overlap, t3 does not.
	if len(results) < 1 {
		t.Fatal("expected at least 1 collision result")
	}
	// t1 itself should also be in the collision set.
	foundT1 := false
	foundT2 := false
	for _, r := range results {
		if r == t1 {
			foundT1 = true
		}
		if r == t2 {
			foundT2 = true
		}
	}
	testutils.Equal(t, foundT1, true)
	testutils.Equal(t, foundT2, true)
}

func TestCollisionWithFilter(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 20))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 20))
	rt.Insert(t1)
	rt.Insert(t2)

	// Filter out self.
	results := rt.Collision(t1, func(s shapes.Spatial) bool {
		return s == t1
	})
	foundT1 := false
	for _, r := range results {
		if r == t1 {
			foundT1 = true
		}
	}
	testutils.Equal(t, foundT1, false)
}

// ---------------------------------------------------------------------------
// MoveObject
// ---------------------------------------------------------------------------

func TestMoveObject(t *testing.T) {
	rt := NewRTree(2, 5)
	obj := NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10))
	rt.Insert(obj)
	testutils.Equal(t, rt.Size(), 1)

	rt.MoveObject(obj, shapes.NewPoint(50, 50))
	testutils.Equal(t, rt.Size(), 1)

	// Should be found at new location.
	results := rt.SearchIntersect(obj.Bounds(), nil)
	testutils.Equal(t, len(results), 1)
	testutils.Equal(t, results[0], obj)
}

func TestMoveObjectUpdatesPosition(t *testing.T) {
	rt := NewRTree(2, 5)
	obj := NewThing(shapes.NewSphere(shapes.NewPoint(10, 20), 5))
	rt.Insert(obj)

	rt.MoveObject(obj, shapes.NewPoint(100, 200))
	// After move, the object's internal position should have changed.
	coords := obj.Sphere.Coordinates()
	testutils.Equal(t, coords[0], 110.0)
	testutils.Equal(t, coords[1], 220.0)
}

// ---------------------------------------------------------------------------
// Depth
// ---------------------------------------------------------------------------

func TestDepthIncreasesWithInserts(t *testing.T) {
	rt := NewRTree(2, 3)
	initialDepth := rt.Depth()

	// Insert enough objects to force splits.
	for i := 0; i < 30; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*50, float64(i)*50), 5)))
	}
	testutils.Equal(t, rt.Size(), 30)
	if rt.Depth() <= initialDepth {
		t.Fatalf("expected depth to increase, initial=%d, current=%d", initialDepth, rt.Depth())
	}
}

// ---------------------------------------------------------------------------
// NearestNeighbor
// ---------------------------------------------------------------------------

func TestNearestNeighborFindsClosest(t *testing.T) {
	rt := NewRTree(2, 5)
	close := NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 5))
	far := NewThing(shapes.NewSphere(shapes.NewPoint(100, 100), 5))
	rt.Insert(close)
	rt.Insert(far)

	obj, dist := rt.NearestNeighbor(shapes.NewPoint(12, 12), nil)
	testutils.Equal(t, obj, close)
	if dist > 10 {
		t.Fatalf("expected small distance, got %f", dist)
	}
}

func TestNearestNeighborWithFilter(t *testing.T) {
	rt := NewRTree(2, 5)
	close := NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 5))
	far := NewThing(shapes.NewSphere(shapes.NewPoint(100, 100), 5))
	rt.Insert(close)
	rt.Insert(far)

	// Filter out the closest one.
	obj, _ := rt.NearestNeighbor(shapes.NewPoint(12, 12), func(s shapes.Spatial) bool {
		return s == close
	})
	testutils.Equal(t, obj, far)
}

func TestNearestNeighborEmptyTree(t *testing.T) {
	rt := NewRTree(2, 5)
	obj, _ := rt.NearestNeighbor(shapes.NewPoint(0, 0), nil)
	if obj != nil {
		t.Fatalf("expected nil for empty tree, got %v", obj)
	}
}

// ---------------------------------------------------------------------------
// NearestNeighbors
// ---------------------------------------------------------------------------

func TestNearestNeighborsFindsKNearest(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 0), 5))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(20, 0), 5))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(100, 0), 5))
	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)

	objs, dists := rt.NearestNeighbors(2, shapes.NewPoint(0, 0), 200, nil)
	testutils.Equal(t, len(objs), 2)
	testutils.Equal(t, len(dists), 2)
	testutils.Equal(t, objs[0], t1)
	testutils.Equal(t, objs[1], t2)
}

func TestNearestNeighborsMaxDistance(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 0), 5))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(1000, 0), 5))
	rt.Insert(t1)
	rt.Insert(t2)

	objs, _ := rt.NearestNeighbors(10, shapes.NewPoint(0, 0), 20, nil)
	testutils.Equal(t, len(objs), 1)
	testutils.Equal(t, objs[0], t1)
}

func TestNearestNeighborsEmptyTree(t *testing.T) {
	rt := NewRTree(2, 5)
	objs, dists := rt.NearestNeighbors(5, shapes.NewPoint(0, 0), 100, nil)
	testutils.Equal(t, len(objs), 0)
	testutils.Equal(t, len(dists), 0)
}

func TestNearestNeighborsWithFilter(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 0), 5))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(20, 0), 5))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(30, 0), 5))
	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)

	// Filter out t1.
	objs, _ := rt.NearestNeighbors(2, shapes.NewPoint(0, 0), 200, func(s shapes.Spatial) bool {
		return s == t1
	})
	testutils.Equal(t, len(objs), 2)
	testutils.Equal(t, objs[0], t2)
	testutils.Equal(t, objs[1], t3)
}

// ---------------------------------------------------------------------------
// GetAllBoundingBoxes
// ---------------------------------------------------------------------------

func TestGetAllBoundingBoxesEmpty(t *testing.T) {
	rt := NewRTree(2, 5)
	boxes := rt.GetAllBoundingBoxes()
	testutils.Equal(t, len(boxes), 0)
}

func TestGetAllBoundingBoxesNonEmpty(t *testing.T) {
	rt := NewRTree(2, 3)
	// Insert enough to cause splits and create non-leaf nodes.
	for i := 0; i < 20; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*50, float64(i)*50), 5)))
	}
	boxes := rt.GetAllBoundingBoxes()
	if len(boxes) == 0 {
		t.Fatal("expected bounding boxes from non-leaf nodes")
	}
}

// ---------------------------------------------------------------------------
// Entity-based deletion (objects implementing GetID)
// ---------------------------------------------------------------------------

func TestEntityBasedDeletion(t *testing.T) {
	rt := NewRTree(2, 5)
	e1 := NewEntityThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10), 1)
	e2 := NewEntityThing(shapes.NewSphere(shapes.NewPoint(50, 50), 10), 2)
	rt.Insert(e1)
	rt.Insert(e2)
	testutils.Equal(t, rt.Size(), 2)

	// Create a different instance with the same ID - should be deleted by ID match.
	e1Copy := NewEntityThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10), 1)
	ok := rt.Delete(e1Copy)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, rt.Size(), 1)
}

func TestEntityBasedDeletionNoMatch(t *testing.T) {
	rt := NewRTree(2, 5)
	e1 := NewEntityThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10), 1)
	rt.Insert(e1)

	// Different ID and different location.
	e3 := NewEntityThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10), 99)
	ok := rt.Delete(e3)
	testutils.Equal(t, ok, false)
	testutils.Equal(t, rt.Size(), 1)
}

// ---------------------------------------------------------------------------
// Mixed Entity and non-Entity objects
// ---------------------------------------------------------------------------

func TestMixedEntityAndNonEntity(t *testing.T) {
	rt := NewRTree(2, 5)
	e1 := NewEntityThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10), 1)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10))
	rt.Insert(e1)
	rt.Insert(t1)
	testutils.Equal(t, rt.Size(), 2)

	// Delete the non-entity by reference.
	ok := rt.Delete(t1)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, rt.Size(), 1)

	// Delete the entity.
	ok = rt.Delete(e1)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, rt.Size(), 0)
}

// ---------------------------------------------------------------------------
// Large-scale insert/delete/search stress test
// ---------------------------------------------------------------------------

func TestLargeScaleOperations(t *testing.T) {
	rt := NewRTree(2, 5)
	n := 100
	things := make([]*Thing, n)
	for i := 0; i < n; i++ {
		things[i] = NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*10, float64(i)*10), 5))
		rt.Insert(things[i])
	}
	testutils.Equal(t, rt.Size(), n)

	// Delete half.
	for i := 0; i < n/2; i++ {
		ok := rt.Delete(things[i])
		testutils.Equal(t, ok, true)
	}
	testutils.Equal(t, rt.Size(), n/2)

	// SearchIntersect should still work.
	results := rt.SearchIntersect(things[n-1].Bounds(), nil)
	if len(results) == 0 {
		t.Fatal("expected to find last inserted object")
	}
}

// ---------------------------------------------------------------------------
// Leaf-only tree (few objects, no splits)
// ---------------------------------------------------------------------------

func TestGetAllBoundingBoxesLeafOnly(t *testing.T) {
	rt := NewRTree(2, 5)
	rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10)))
	rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(10, 10), 10)))
	// With only 2 objects, tree is leaf-only, so no non-leaf bounding boxes.
	boxes := rt.GetAllBoundingBoxes()
	testutils.Equal(t, len(boxes), 0)
}

// ---------------------------------------------------------------------------
// NearestNeighbor with deep tree (exercises sortEntries, pruneEntries, and
// the non-leaf branch of nearestNeighbor)
// ---------------------------------------------------------------------------

func TestNearestNeighborDeepTree(t *testing.T) {
	rt := NewRTree(2, 3)
	// Insert many spread-out objects to build a deep tree.
	for i := 0; i < 50; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*100, float64(i)*100), 5)))
	}
	if rt.Depth() < 2 {
		t.Fatalf("expected deep tree, got depth %d", rt.Depth())
	}
	// Find nearest to origin - should be the first inserted object.
	obj, dist := rt.NearestNeighbor(shapes.NewPoint(0, 0), nil)
	if obj == nil {
		t.Fatal("expected to find nearest neighbor in deep tree")
	}
	if dist > 10 {
		t.Fatalf("expected small distance, got %f", dist)
	}
}

func TestNearestNeighborDeepTreeWithFilter(t *testing.T) {
	rt := NewRTree(2, 3)
	things := make([]*Thing, 30)
	for i := 0; i < 30; i++ {
		things[i] = NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*50, 0), 5))
		rt.Insert(things[i])
	}
	// Filter out the nearest, should get the second nearest.
	obj, _ := rt.NearestNeighbor(shapes.NewPoint(0, 0), func(s shapes.Spatial) bool {
		return s == things[0]
	})
	if obj == nil {
		t.Fatal("expected non-nil result")
	}
	testutils.Equal(t, obj, things[1])
}

// ---------------------------------------------------------------------------
// NearestNeighbors with deep tree (exercises sortPreselectedEntries,
// pruneEntriesMinDist, and non-leaf branch of nearestNeighbors)
// ---------------------------------------------------------------------------

func TestNearestNeighborsDeepTree(t *testing.T) {
	rt := NewRTree(2, 3)
	things := make([]*Thing, 50)
	for i := 0; i < 50; i++ {
		things[i] = NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*100, 0), 5))
		rt.Insert(things[i])
	}
	if rt.Depth() < 2 {
		t.Fatalf("expected deep tree, got depth %d", rt.Depth())
	}
	// Ask for 3 nearest to origin.
	objs, dists := rt.NearestNeighbors(3, shapes.NewPoint(0, 0), 10000, nil)
	testutils.Equal(t, len(objs), 3)
	testutils.Equal(t, len(dists), 3)
	// The first should be the closest.
	testutils.Equal(t, objs[0], things[0])
	testutils.Equal(t, objs[1], things[1])
	testutils.Equal(t, objs[2], things[2])
}

func TestNearestNeighborsDeepTreeMaxDistPrune(t *testing.T) {
	rt := NewRTree(2, 3)
	for i := 0; i < 50; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*100, 0), 5)))
	}
	// With a small maxDistance, only nearby objects should be returned.
	objs, _ := rt.NearestNeighbors(10, shapes.NewPoint(0, 0), 50, nil)
	if len(objs) > 3 {
		t.Fatalf("expected few results with small maxDistance, got %d", len(objs))
	}
}

// ---------------------------------------------------------------------------
// Entity-based findLeaf paths (exercises the Entity branch in findLeaf)
// ---------------------------------------------------------------------------

func TestEntityFindLeafDeepTree(t *testing.T) {
	rt := NewRTree(2, 3)
	entities := make([]*EntityThing, 30)
	for i := 0; i < 30; i++ {
		entities[i] = NewEntityThing(shapes.NewSphere(shapes.NewPoint(float64(i)*50, float64(i)*50), 5), uint32(i+1))
		rt.Insert(entities[i])
	}
	if rt.Depth() < 2 {
		t.Fatalf("expected deep tree for entity test")
	}
	// Delete by ID match using a copy (different pointer, same ID).
	entityCopy := NewEntityThing(shapes.NewSphere(shapes.NewPoint(250, 250), 5), 6) // matches entities[5]
	ok := rt.Delete(entityCopy)
	testutils.Equal(t, ok, true)
	testutils.Equal(t, rt.Size(), 29)
}

// ---------------------------------------------------------------------------
// condenseTree with underflow cascading (delete many to trigger re-insert)
// ---------------------------------------------------------------------------

func TestCondenseTreeUnderflow(t *testing.T) {
	rt := NewRTree(2, 4)
	things := make([]*Thing, 20)
	for i := 0; i < 20; i++ {
		things[i] = NewThing(shapes.NewSphere(shapes.NewPoint(float64(i)*30, float64(i)*30), 5))
		rt.Insert(things[i])
	}
	// Delete enough objects from one area to trigger underflow.
	for i := 0; i < 15; i++ {
		ok := rt.Delete(things[i])
		testutils.Equal(t, ok, true)
	}
	testutils.Equal(t, rt.Size(), 5)
	// Remaining objects should still be findable.
	for i := 15; i < 20; i++ {
		results := rt.SearchIntersect(things[i].Bounds(), nil)
		found := false
		for _, r := range results {
			if r == things[i] {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected to find thing[%d] after condense", i)
		}
	}
}

// ---------------------------------------------------------------------------
// assignGroup tie-breaking paths (area equality, entry count)
// ---------------------------------------------------------------------------

func TestInsertIdenticalObjectsTriggerAssignGroupTieBreaks(t *testing.T) {
	rt := NewRTree(2, 3)
	// Insert many objects at identical positions to trigger tie-breaking in assignGroup.
	for i := 0; i < 20; i++ {
		rt.Insert(NewThing(shapes.NewSphere(shapes.NewPoint(0, 0), 10)))
	}
	testutils.Equal(t, rt.Size(), 20)
	// All should still be findable.
	results := rt.SearchIntersect(shapes.NewBox(shapes.NewPoint(-15, -15), 30, 30), nil)
	testutils.Equal(t, len(results), 20)
}
