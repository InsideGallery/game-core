package rtree

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

type Thing struct {
	shapes.Sphere
}

func NewThing(where shapes.Sphere) *Thing {
	return &Thing{
		Sphere: where,
	}
}

func (t *Thing) SetNewSpatial(s shapes.Spatial) {
	t.Sphere = s.(shapes.Sphere)
}

func (t *Thing) Bounds() shapes.Box {
	return t.Sphere.Bounds()
}

func (t *Thing) UpdateSpatial(s shapes.Spatial) {
	t.Sphere = s.(shapes.Sphere)
}

func TestRTree(t *testing.T) {
	rt := NewRTree(25, 50)

	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(1, 1), 100))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(0.5, 1), 40))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(100, 0), 10))
	t4 := NewThing(shapes.NewSphere(shapes.NewPoint(120, 0), 10))
	t5 := NewThing(shapes.NewSphere(shapes.NewPoint(125, 0), 20))

	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)
	rt.Insert(t4)
	rt.Insert(t5)

	results := rt.SearchIntersect(t2.Bounds(), nil)
	testutils.Equal(t, len(results), 2)
	testutils.Equal(t, results[0], t1)
	testutils.Equal(t, results[1], t2)
	results = rt.SearchIntersect(t4.Bounds(), nil)
	testutils.Equal(t, len(results), 3)
	testutils.Equal(t, results[0], t3)
	testutils.Equal(t, results[1], t4)
	testutils.Equal(t, results[2], t5)
	rt.MoveObject(t5, shapes.NewPoint(0, 20))
	results = rt.SearchIntersect(t5.Bounds(), nil)
	testutils.Equal(t, len(results), 3)
	testutils.Equal(t, results[0], t3)
	testutils.Equal(t, results[1], t4)
	testutils.Equal(t, results[2], t5)

	results = rt.SearchIntersect(t5.Bounds(), func(spatial shapes.Spatial) bool {
		return spatial.Point1().Equals(t3.Point1())
	})
	testutils.Equal(t, len(results), 2)
	testutils.Equal(t, results[0], t4)
	testutils.Equal(t, results[1], t5)
}

func TestNearest(t *testing.T) {
	rt := NewRTree(2, 5)
	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(105, 105), 5))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(205, 105), 5))
	rt.Insert(t1)
	rt.Insert(t2)
	sp, _ := rt.NearestNeighbors(10, shapes.NewPoint(154, 105), 46, nil)
	sp2, _ := rt.NearestNeighbors(10, shapes.NewPoint(156, 105), 46, nil)
	testutils.Equal(t, sp[0], t1)
	testutils.Equal(t, sp[1], t2)
	testutils.Equal(t, sp2[0], t2)
	testutils.Equal(t, sp2[1], t1)
	sp3, _ := rt.NearestNeighbors(10, shapes.NewPoint(156, 105), 44, nil)
	testutils.Equal(t, sp3[0], t2)
	sp4, _ := rt.NearestNeighbors(10, shapes.NewPoint(156, 105), 10, nil)
	testutils.Equal(t, len(sp4), 0)
	o, d := rt.NearestNeighbor(shapes.NewPoint(156, 105), nil)
	testutils.Equal(t, o, t2)
	testutils.Equal(t, d, 44.0)
}

func TestCollision(t *testing.T) {
	rt := NewRTree(2, 5)

	t1 := NewThing(shapes.NewSphere(shapes.NewPoint(1, 1), 100))
	t2 := NewThing(shapes.NewSphere(shapes.NewPoint(0.5, 1), 40))
	t3 := NewThing(shapes.NewSphere(shapes.NewPoint(100, 0), 10))
	t4 := NewThing(shapes.NewSphere(shapes.NewPoint(120, 0), 10))
	t5 := NewThing(shapes.NewSphere(shapes.NewPoint(125, 0), 20))
	t6 := NewThing(shapes.NewSphere(shapes.NewPoint(122, 2), 5))
	t7 := NewThing(shapes.NewSphere(shapes.NewPoint(0.5, 159), 40))
	t8 := NewThing(shapes.NewSphere(shapes.NewPoint(100, 120), 10))
	t9 := NewThing(shapes.NewSphere(shapes.NewPoint(120, 130), 10))
	t10 := NewThing(shapes.NewSphere(shapes.NewPoint(125, 140), 20))

	rt.Insert(t1)
	rt.Insert(t2)
	rt.Insert(t3)
	rt.Insert(t4)
	rt.Insert(t5)
	rt.Insert(t6)
	rt.Insert(t7)
	rt.Insert(t8)
	rt.Insert(t9)
	rt.Insert(t10)

	rt.Delete(t1)
	rt.Delete(t2)
	rt.Delete(t3)
	rt.Delete(t5)
	rt.Delete(t4)
	rt.Delete(t6)
	rt.Delete(t7)
	rt.Delete(t8)
	rt.Delete(t9)
	rt.Delete(t10)

	results := rt.SearchIntersect(t4.Bounds(), nil)
	testutils.Equal(t, len(results), 0)
}
