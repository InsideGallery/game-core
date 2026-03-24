package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestNewMultiObject(t *testing.T) {
	p := NewPoint(1, 2, 3)
	b := NewBox(NewPoint(10, 20, 30), 5, 5, 5)
	m := NewMultiObject(p, b)
	testutils.Equal(t, m.Count(), 2)
}

func TestMultiObject_Point1(t *testing.T) {
	p := NewPoint(5, 10, 15)
	b := NewBox(NewPoint(1, 2, 3), 10, 10, 10)
	m := NewMultiObject(p, b)
	// Point1 returns Bounds().Point1(), which is the bounding box min
	// BoundingBox starts from a zero Box, so min includes origin (0,0,0)
	p1 := m.Point1()
	testutils.Equal(t, p1.Coordinate(0), 0.0)
	testutils.Equal(t, p1.Coordinate(1), 0.0)
	testutils.Equal(t, p1.Coordinate(2), 0.0)
}

func TestMultiObject_Objects(t *testing.T) {
	p := NewPoint(1, 2, 3)
	b := NewBox(NewPoint(4, 5, 6), 1, 1, 1)
	m := NewMultiObject(p, b)
	objs := m.Objects()
	testutils.Equal(t, len(objs), 2)
}

func TestMultiObject_Count(t *testing.T) {
	tests := map[string]struct {
		objects []Spatial
		want    int
	}{
		"empty": {
			objects: nil,
			want:    0,
		},
		"one": {
			objects: []Spatial{NewPoint(1, 2, 3)},
			want:    1,
		},
		"three": {
			objects: []Spatial{NewPoint(1, 2, 3), NewPoint(4, 5, 6), NewBox(NewPoint(0, 0), 5, 5)},
			want:    3,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			m := NewMultiObject(tc.objects...)
			testutils.Equal(t, m.Count(), tc.want)
		})
	}
}

func TestMultiObject_Coordinates(t *testing.T) {
	p := NewPoint(5, 10, 0)
	b := NewBox(NewPoint(1, 2, 0), 10, 10, 0)
	m := NewMultiObject(p, b)
	coords := m.Coordinates()
	// BoundingBox starts from a zero Box, so min includes origin (0,0,0)
	testutils.Equal(t, coords[0], 0.0)
	testutils.Equal(t, coords[1], 0.0)
	testutils.Equal(t, coords[2], 0.0)
}

func TestMultiObject_Object(t *testing.T) {
	p := NewPoint(1, 2, 3)
	b := NewBox(NewPoint(4, 5, 6), 1, 1, 1)
	m := NewMultiObject(p, b)
	obj0 := m.Object(0)
	testutils.Equal(t, obj0.Point1(), NewPoint(1, 2, 3))
	obj1 := m.Object(1)
	testutils.Equal(t, obj1.Point1(), NewPoint(4, 5, 6))
}

func TestMultiObject_Get(t *testing.T) {
	m := NewMultiObject(NewPoint(1, 2, 3))
	got := m.Get()
	_, ok := got.(MultiObject)
	testutils.Equal(t, ok, true)
}

func TestMultiObject_Move(t *testing.T) {
	p := NewPoint(1, 2, 3)
	b := NewBox(NewPoint(10, 20, 30), 5, 5, 5)
	m := NewMultiObject(p, b)
	moved := m.Move(NewPoint(10, 10, 10)).(MultiObject)
	testutils.Equal(t, moved.Count(), 2)
	// First object should be moved point
	testutils.Equal(t, moved.Object(0).Point1(), NewPoint(11, 12, 13))
	// Second object should be moved box
	testutils.Equal(t, moved.Object(1).Point1(), NewPoint(20, 30, 40))
}

func TestMultiObject_Bounds(t *testing.T) {
	p := NewPoint(0, 0, 0)
	b := NewBox(NewPoint(5, 5, 5), 10, 10, 10)
	m := NewMultiObject(p, b)
	bounds := m.Bounds()
	testutils.Equal(t, bounds.Point1(), NewPoint(0, 0, 0))
	testutils.Equal(t, bounds.Point2(), NewPoint(15, 15, 15))
}

func TestMultiObject_Center(t *testing.T) {
	p1 := NewPoint(0, 0, 0)
	p2 := NewPoint(10, 10, 10)
	m := NewMultiObject(p1, p2)
	c := m.Center()
	// Center of the two sub-centers: (0,0,0) and (10,10,10) => average (5,5,5)
	testutils.Equal(t, c.Coordinate(0), 5.0)
	testutils.Equal(t, c.Coordinate(1), 5.0)
	testutils.Equal(t, c.Coordinate(2), 5.0)
}

func TestMultiObject_Support(t *testing.T) {
	// Use objects that implement Collide (have Support method)
	s1 := NewSphere(NewPoint(0, 0, 0), 5)
	s2 := NewSphere(NewPoint(10, 0, 0), 5)
	m := NewMultiObject(s1, s2)
	// Support in positive x direction should come from s2
	p := m.Support(NewPoint(1, 0, 0))
	testutils.Equal(t, p, NewPoint(15, 0, 0))
}

func TestMultiObject_Support_NoCollideObjects(t *testing.T) {
	// Points don't implement Collide interface (they do via Support method)
	// Actually Point does have Support. Let's use a mix.
	p := NewPoint(5, 0, 0)
	b := NewBox(NewPoint(10, 0, 0), 10, 10, 10)
	m := NewMultiObject(p, b)
	sp := m.Support(NewPoint(1, 0, 0))
	// Box Support in +x should return Point2.x = 20
	testutils.Equal(t, sp, NewPoint(20, 0, 0))
}
