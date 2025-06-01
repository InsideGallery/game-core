package hexagone

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestAxialToCube(t *testing.T) {
	testcases := map[string]struct {
		axial Axial
		cube  Cube
	}{
		"axial;0;0": {
			axial: NewAxial(0, 0),
			cube:  NewCube(0, 0, 0),
		},
		"axial;10;10": {
			axial: NewAxial(10, 10),
			cube:  NewCube(10, -20, 10),
		},
		"axial;10;0": {
			axial: NewAxial(10, 0),
			cube:  NewCube(10, -10, 0),
		},
		"axial;0;10": {
			axial: NewAxial(0, 10),
			cube:  NewCube(0, -10, 10),
		},
		"axial;-10;10": {
			axial: NewAxial(-10, 10),
			cube:  NewCube(-10, 0, 10),
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.axial.ToCube(), test.cube)
		})
	}
}

func TestCubeToAxial(t *testing.T) {
	testcases := map[string]struct {
		axial Axial
		cube  Cube
	}{
		"axial;0;0": {
			axial: NewAxial(0, 0),
			cube:  NewCube(0, 0, 0),
		},
		"axial;10;10": {
			axial: NewAxial(10, 10),
			cube:  NewCube(10, -20, 10),
		},
		"axial;10;0": {
			axial: NewAxial(10, 0),
			cube:  NewCube(10, -10, 0),
		},
		"axial;0;10": {
			axial: NewAxial(0, 10),
			cube:  NewCube(0, -10, 10),
		},
		"axial;-10;10": {
			axial: NewAxial(-10, 10),
			cube:  NewCube(-10, 0, 10),
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.cube.ToAxis(), test.axial)
		})
	}
}

func TestCubeDistance(t *testing.T) {
	testcases := map[string]struct {
		cube   Cube
		cube2  Cube
		result int
	}{
		"cube;0;0;0": {
			cube:   NewCube(0, 0, 0),
			cube2:  NewCube(0, 0, 0),
			result: 0,
		},
		"cube;10;0;0": {
			cube:   NewCube(10, 0, 0),
			cube2:  NewCube(0, 0, 0),
			result: 5,
		},
		"cube;0;10;0": {
			cube:   NewCube(0, 10, 0),
			cube2:  NewCube(0, 0, 0),
			result: 5,
		},
		"cube;0;0;10": {
			cube:   NewCube(0, 0, 10),
			cube2:  NewCube(0, 0, 0),
			result: 5,
		},
		"cube;10;10;0": {
			cube:   NewCube(10, 10, 0),
			cube2:  NewCube(0, 0, 0),
			result: 10,
		},
		"cube;0;10;10": {
			cube:   NewCube(0, 10, 10),
			cube2:  NewCube(0, 0, 0),
			result: 10,
		},
		"cube;10;10;10": {
			cube:   NewCube(10, 10, 10),
			cube2:  NewCube(0, 0, 0),
			result: 15,
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.cube.Distance(test.cube2), test.result)
		})
	}
}

func TestAxialDistance(t *testing.T) {
	testcases := map[string]struct {
		axial  Axial
		axial2 Axial
		result int
	}{
		"axial;0;0": {
			axial:  NewAxial(0, 0),
			axial2: NewAxial(0, 0),
			result: 0,
		},
		"axial;10;0": {
			axial:  NewAxial(10, 0),
			axial2: NewAxial(0, 0),
			result: 10,
		},
		"axial;0;10": {
			axial:  NewAxial(0, 10),
			axial2: NewAxial(0, 0),
			result: 10,
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.axial.Distance(test.axial2), test.result)
		})
	}
}

func TestAxialToPosition(t *testing.T) {
	testcases := map[string]struct {
		axial Axial
		point shapes.Point
	}{
		"axial;0;0": {
			axial: NewAxial(0, 0),
			point: shapes.NewPoint(0, 0),
		},
		"axial;10;10": {
			axial: NewAxial(10, 10),
			point: shapes.NewPoint(25.980762113533157, 15),
		},
		"axial;-10;-10": {
			axial: NewAxial(-10, -10),
			point: shapes.NewPoint(-25.980762113533157, -15),
		},
		"axial;1;0": {
			axial: NewAxial(1, 0),
			point: shapes.NewPoint(1.7320508075688772, 0),
		},
	}

	for name, test := range testcases {
		test := test
		t.Run(name, func(t *testing.T) {
			testutils.Equal(t, test.axial.ToPosition(), test.point)
		})
	}
}
