package shapes

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestNewRotatable(t *testing.T) {
	p := NewPoint(1, 2, 3)
	r := NewRotatable(p, 45.0)
	testutils.Equal(t, r.Angle, 45.0)
	testutils.Equal(t, r.Point1(), NewPoint(1, 2, 3))
}

func TestRotatable_GetAngle(t *testing.T) {
	tests := map[string]struct {
		angle float64
	}{
		"zero":     {angle: 0},
		"positive": {angle: 90},
		"negative": {angle: -45},
		"large":    {angle: 360},
		"fraction": {angle: 33.33},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			r := NewRotatable(NewPoint(0, 0, 0), tc.angle)
			testutils.Equal(t, r.GetAngle(), tc.angle)
		})
	}
}

func TestRotatable_Spatial(t *testing.T) {
	b := NewBox(NewPoint(0, 0, 0), 10, 10, 10)
	r := NewRotatable(b, 30.0)
	testutils.Equal(t, r.Bounds(), b.Bounds())
	testutils.Equal(t, r.Center(), b.Center())
}
