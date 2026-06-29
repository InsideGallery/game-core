package camera_test

import (
	"math"
	"testing"

	"github.com/InsideGallery/game-core/client/camera"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		pos        shapes.Point
		zoom       float64
		vpW, vpH   int
		wx, wy     float64
	}{
		{"no transform", shapes.NewPoint(0, 0), 0, 800, 600, 100, 200},
		{"pan", shapes.NewPoint(50, 30), 0, 800, 600, 100, 200},
		{"zoom in", shapes.NewPoint(0, 0), 10, 800, 600, 100, 200},
		{"pan+zoom", shapes.NewPoint(50, 30), 5, 1024, 768, -40, 60},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := camera.NewCamera(tc.pos)
			c.ZoomFactor = tc.zoom
			c.SetViewPort(tc.vpW, tc.vpH)

			sx, sy := c.WorldToScreen(tc.wx, tc.wy)
			gotX, gotY := c.ScreenToWorld(sx, sy)

			const eps = 1e-9
			if math.Abs(gotX-tc.wx) > eps || math.Abs(gotY-tc.wy) > eps {
				t.Errorf("round-trip: want (%v,%v) got (%v,%v)", tc.wx, tc.wy, gotX, gotY)
			}
		})
	}
}

func TestScreenToWorldWithZeroZoom(t *testing.T) {
	c := camera.NewCamera(shapes.NewPoint(0, 0))
	c.SetViewPort(800, 600)

	// At zoom=0 and position=(0,0), screen centre maps to world (0,0)
	wx, wy := c.ScreenToWorld(400, 300)
	const eps = 1e-9
	if math.Abs(wx) > eps || math.Abs(wy) > eps {
		t.Errorf("centre: want (0,0) got (%v,%v)", wx, wy)
	}
}
