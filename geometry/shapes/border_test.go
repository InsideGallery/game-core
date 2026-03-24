package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestNewBorder(t *testing.T) {
	b := NewBorder(NewBox(NewPoint(0, 0, 0), 100, 100, 100))
	testutils.Equal(t, b.Point1(), NewPoint(0, 0, 0))
	testutils.Equal(t, b.Point2(), NewPoint(100, 100, 100))
}

func TestBorder_Get(t *testing.T) {
	b := NewBorder(NewBox(NewPoint(0, 0), 100, 100))
	got := b.Get()
	_, ok := got.(Border)
	testutils.Equal(t, ok, true)
}

func TestBorder_Collision_Point(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj        Spatial
		dims       int
		wantCol    bool
		wantNonZero bool
	}{
		"point inside": {
			obj:        NewPoint(50, 50),
			dims:       2,
			wantCol:    false,
			wantNonZero: false,
		},
		"point outside left": {
			obj:        NewPoint(-5, 50),
			dims:       2,
			wantCol:    true,
			wantNonZero: true,
		},
		"point outside right": {
			obj:        NewPoint(105, 50),
			dims:       2,
			wantCol:    true,
			wantNonZero: true,
		},
		"point outside top": {
			obj:        NewPoint(50, -5),
			dims:       2,
			wantCol:    true,
			wantNonZero: true,
		},
		"point outside bottom": {
			obj:        NewPoint(50, 105),
			dims:       2,
			wantCol:    true,
			wantNonZero: true,
		},
		"point on boundary": {
			obj:        NewPoint(0, 0),
			dims:       2,
			wantCol:    false,
			wantNonZero: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, depth := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
			if tc.wantNonZero {
				hasNonZero := false
				for _, v := range depth {
					if v != 0 {
						hasNonZero = true
						break
					}
				}
				testutils.Equal(t, hasNonZero, true)
			}
		})
	}
}

func TestBorder_Collision_Sphere(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj     Spatial
		dims    int
		wantCol bool
	}{
		"sphere inside": {
			obj:     NewSphere(NewPoint(50, 50), 10),
			dims:    2,
			wantCol: false,
		},
		"sphere outside left": {
			obj:     NewSphere(NewPoint(5, 50), 10),
			dims:    2,
			wantCol: true,
		},
		"sphere outside right": {
			obj:     NewSphere(NewPoint(95, 50), 10),
			dims:    2,
			wantCol: true,
		},
		"sphere outside top": {
			obj:     NewSphere(NewPoint(50, 5), 10),
			dims:    2,
			wantCol: true,
		},
		"sphere outside bottom": {
			obj:     NewSphere(NewPoint(50, 95), 10),
			dims:    2,
			wantCol: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, _ := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
		})
	}
}

func TestBorder_Collision_Box(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj     Spatial
		dims    int
		wantCol bool
	}{
		"box inside": {
			obj:     NewBox(NewPoint(10, 10), 20, 20),
			dims:    2,
			wantCol: false,
		},
		"box outside left": {
			obj:     NewBox(NewPoint(-5, 10), 20, 20),
			dims:    2,
			wantCol: true,
		},
		"box outside right": {
			obj:     NewBox(NewPoint(85, 10), 20, 20),
			dims:    2,
			wantCol: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, _ := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
		})
	}
}

func TestBorder_Collision_Line(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj     Spatial
		dims    int
		wantCol bool
	}{
		"line inside": {
			obj:     NewLine(NewPoint(10, 10), NewPoint(90, 90)),
			dims:    2,
			wantCol: false,
		},
		"line p1 outside left": {
			obj:     NewLine(NewPoint(-5, 50), NewPoint(50, 50)),
			dims:    2,
			wantCol: true,
		},
		"line p1 outside right": {
			obj:     NewLine(NewPoint(105, 50), NewPoint(50, 50)),
			dims:    2,
			wantCol: true,
		},
		"line p2 outside left": {
			obj:     NewLine(NewPoint(50, 50), NewPoint(-5, 50)),
			dims:    2,
			wantCol: true,
		},
		"line p2 outside right": {
			obj:     NewLine(NewPoint(50, 50), NewPoint(105, 50)),
			dims:    2,
			wantCol: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, _ := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
		})
	}
}

func TestBorder_Collision_Polyhedron(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj     Spatial
		dims    int
		wantCol bool
	}{
		"polyhedron inside": {
			obj:     NewPolyhedron(NewPoint(10, 10), NewPoint(90, 10), NewPoint(90, 90), NewPoint(10, 90)),
			dims:    2,
			wantCol: false,
		},
		"polyhedron partially outside": {
			obj:     NewPolyhedron(NewPoint(-5, 10), NewPoint(50, 10), NewPoint(50, 90), NewPoint(-5, 90)),
			dims:    2,
			wantCol: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, _ := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
		})
	}
}

func TestBorder_Collision_MultiObject(t *testing.T) {
	border := NewBorder(NewBox(NewPoint(0, 0), 100, 100))

	tests := map[string]struct {
		obj     Spatial
		dims    int
		wantCol bool
	}{
		"multiobject inside": {
			obj:     NewMultiObject(NewPoint(50, 50), NewBox(NewPoint(10, 10), 20, 20)),
			dims:    2,
			wantCol: false,
		},
		"multiobject with outside element": {
			obj:     NewMultiObject(NewPoint(-5, 50), NewBox(NewPoint(10, 10), 20, 20)),
			dims:    2,
			wantCol: true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			col, _ := border.Collision(tc.obj, tc.dims)
			testutils.Equal(t, col, tc.wantCol)
		})
	}
}
