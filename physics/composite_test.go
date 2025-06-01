package physics

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestComposite(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(-640, -480), 640, 480)), shapes.NewPoint(0, 2), 4)
	rope := NewComposite()
	w.AddComposites(rope)
	material := NewMaterial(1)
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 50), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 90), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 130), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 170), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 210), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 250), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 290), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 330), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 370), material))
	rope.AddParticle(NewParticle(shapes.NewPoint(0, 410), material))
	testutils.Equal(t, rope.AddConstraints(0, 1, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(1, 2, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(2, 3, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(3, 4, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(4, 5, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(5, 6, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(6, 7, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(7, 8, 1.0), nil)
	testutils.Equal(t, rope.AddConstraints(8, 9, 1.0), nil)
	w.Simulate(4, 2)
}
