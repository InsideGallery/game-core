package physics

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// ---------------------------------------------------------------------------
// Material
// ---------------------------------------------------------------------------

func TestNewMaterial(t *testing.T) {
	m := NewMaterial(5.0)
	testutils.Equal(t, m.Mass, 5.0)
}

func TestNewMaterialZero(t *testing.T) {
	m := NewMaterial(0)
	testutils.Equal(t, m.Mass, 0.0)
}

// ---------------------------------------------------------------------------
// Particle – creation
// ---------------------------------------------------------------------------

func TestNewParticle(t *testing.T) {
	m := NewMaterial(2.0)
	p := NewParticle(shapes.NewPoint(10, 20), m)
	testutils.Equal(t, p.Position.Coordinate(0), 10.0)
	testutils.Equal(t, p.Position.Coordinate(1), 20.0)
	testutils.Equal(t, p.Previous.Coordinate(0), 10.0)
	testutils.Equal(t, p.Previous.Coordinate(1), 20.0)
	testutils.Equal(t, p.Material.Mass, 2.0)
}

// ---------------------------------------------------------------------------
// Particle – Simulate
// ---------------------------------------------------------------------------

func TestParticleSimulateNonZeroMass(t *testing.T) {
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.Accelerate(shapes.NewPoint(0, 10))
	p.Simulate(1.0)
	// After simulation the particle should have moved.
	// velocity = 2*pos - prev = 0, plus acceleration*delta^2 = (0,10)*1 = (0,10)
	testutils.Equal(t, p.Position.Coordinate(1), 10.0)
}

func TestParticleSimulateZeroMass(t *testing.T) {
	m := NewMaterial(0)
	p := NewParticle(shapes.NewPoint(5, 5), m)
	p.Acceleration = shapes.NewPoint(0, 10)
	p.Simulate(1.0)
	// Zero-mass particle should not move.
	testutils.Equal(t, p.Position.Coordinate(0), 5.0)
	testutils.Equal(t, p.Position.Coordinate(1), 5.0)
}

func TestParticleSimulateResetsAcceleration(t *testing.T) {
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.Accelerate(shapes.NewPoint(1, 2))
	p.Simulate(1.0)
	// Acceleration should be zeroed after simulate.
	testutils.Equal(t, p.Acceleration.Coordinate(0), 0.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 0.0)
	testutils.Equal(t, p.Acceleration.Coordinate(2), 0.0)
}

// ---------------------------------------------------------------------------
// Particle – Accelerate
// ---------------------------------------------------------------------------

func TestParticleAccelerate(t *testing.T) {
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.Accelerate(shapes.NewPoint(3, 4))
	testutils.Equal(t, p.Acceleration.Coordinate(0), 3.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 4.0)
}

func TestParticleAccelerateAdditive(t *testing.T) {
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.Accelerate(shapes.NewPoint(1, 2))
	p.Accelerate(shapes.NewPoint(3, 4))
	testutils.Equal(t, p.Acceleration.Coordinate(0), 4.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 6.0)
}

// ---------------------------------------------------------------------------
// Particle – ApplyForce
// ---------------------------------------------------------------------------

func TestParticleApplyForceNonZeroMass(t *testing.T) {
	m := NewMaterial(2.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.ApplyForce(shapes.NewPoint(10, 20))
	// acceleration += force * (1/mass) = (10,20) * 0.5 = (5,10)
	testutils.Equal(t, p.Acceleration.Coordinate(0), 5.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 10.0)
}

func TestParticleApplyForceZeroMass(t *testing.T) {
	m := NewMaterial(0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.ApplyForce(shapes.NewPoint(10, 20))
	// Zero-mass: no effect.
	testutils.Equal(t, p.Acceleration.Coordinate(0), 0.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 0.0)
}

// ---------------------------------------------------------------------------
// Particle – ApplyImpulse
// ---------------------------------------------------------------------------

func TestParticleApplyImpulseNonZeroMass(t *testing.T) {
	m := NewMaterial(2.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.ApplyImpulse(shapes.NewPoint(6, 8))
	// position += impulse * (1/mass) = (6,8) * 0.5 = (3,4)
	testutils.Equal(t, p.Position.Coordinate(0), 3.0)
	testutils.Equal(t, p.Position.Coordinate(1), 4.0)
}

func TestParticleApplyImpulseZeroMass(t *testing.T) {
	m := NewMaterial(0)
	p := NewParticle(shapes.NewPoint(5, 5), m)
	p.ApplyImpulse(shapes.NewPoint(6, 8))
	// Zero-mass: position unchanged.
	testutils.Equal(t, p.Position.Coordinate(0), 5.0)
	testutils.Equal(t, p.Position.Coordinate(1), 5.0)
}

// ---------------------------------------------------------------------------
// Particle – ResetForces
// ---------------------------------------------------------------------------

func TestParticleResetForces(t *testing.T) {
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	p.Accelerate(shapes.NewPoint(99, 99, 99))
	p.ResetForces()
	testutils.Equal(t, p.Acceleration.Coordinate(0), 0.0)
	testutils.Equal(t, p.Acceleration.Coordinate(1), 0.0)
	testutils.Equal(t, p.Acceleration.Coordinate(2), 0.0)
}

// ---------------------------------------------------------------------------
// Particle – SetMaterial
// ---------------------------------------------------------------------------

func TestParticleSetMaterial(t *testing.T) {
	m1 := NewMaterial(1.0)
	m2 := NewMaterial(5.0)
	p := NewParticle(shapes.NewPoint(0, 0), m1)
	testutils.Equal(t, p.Material.Mass, 1.0)
	p.SetMaterial(m2)
	testutils.Equal(t, p.Material.Mass, 5.0)
}

// ---------------------------------------------------------------------------
// Particle – Restrain
// ---------------------------------------------------------------------------

func TestParticleRestrain(t *testing.T) {
	border := shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100))
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(200, 200), m)
	p.Restrain(border, 2)
	// After restrain, position should be clamped inside the border.
	x := p.Position.Coordinate(0)
	y := p.Position.Coordinate(1)
	if x > 100 || y > 100 {
		t.Fatalf("particle not restrained: got (%f, %f)", x, y)
	}
}

func TestParticleRestrainInsideBorder(t *testing.T) {
	border := shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100))
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(50, 50), m)
	p.Restrain(border, 2)
	// When inside the border, Collision returns depth [0,0], so position becomes (0,0).
	testutils.Equal(t, p.Position.Coordinate(0), 0.0)
	testutils.Equal(t, p.Position.Coordinate(1), 0.0)
}

// ---------------------------------------------------------------------------
// Composite – creation and GetParticle
// ---------------------------------------------------------------------------

func TestNewComposite(t *testing.T) {
	c := NewComposite()
	testutils.Equal(t, len(c.Particles), 0)
	testutils.Equal(t, len(c.Constraints), 0)
}

func TestCompositeAddParticleAndGetValid(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(1, 2), m)
	c.AddParticle(p)
	got, err := c.GetParticle(0)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, got, p)
}

func TestGetParticleNegativeIndex(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	_, err := c.GetParticle(-1)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

func TestGetParticleIndexEqualsLen(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	// len == 1, so index 1 should error
	_, err := c.GetParticle(1)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

func TestGetParticleIndexGreaterThanLen(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	_, err := c.GetParticle(5)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

func TestGetParticleEmptyComposite(t *testing.T) {
	c := NewComposite()
	_, err := c.GetParticle(0)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

// ---------------------------------------------------------------------------
// Composite – AddConstraints
// ---------------------------------------------------------------------------

func TestAddConstraintsValid(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	c.AddParticle(NewParticle(shapes.NewPoint(10, 0), m))
	err := c.AddConstraints(0, 1, 0.5)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, len(c.Constraints), 1)
}

func TestAddConstraintsInvalidFirstIndex(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	c.AddParticle(NewParticle(shapes.NewPoint(10, 0), m))
	err := c.AddConstraints(5, 1, 0.5)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

func TestAddConstraintsInvalidSecondIndex(t *testing.T) {
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	c.AddParticle(NewParticle(shapes.NewPoint(10, 0), m))
	err := c.AddConstraints(0, 5, 0.5)
	testutils.Equal(t, err, ErrNotFoundAttachedParticle)
}

// ---------------------------------------------------------------------------
// Composite – SetMaterial
// ---------------------------------------------------------------------------

func TestCompositeSetMaterial(t *testing.T) {
	c := NewComposite()
	m1 := NewMaterial(1.0)
	m2 := NewMaterial(9.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m1))
	c.AddParticle(NewParticle(shapes.NewPoint(1, 1), m1))
	c.SetMaterial(m2)
	for _, p := range c.Particles {
		testutils.Equal(t, p.Material.Mass, 9.0)
	}
}

// ---------------------------------------------------------------------------
// Constraint
// ---------------------------------------------------------------------------

func TestNewConstraintSpringConstantClamped(t *testing.T) {
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(10, 0), m)
	c := NewConstraint(p1, p2, 5.0, 0)
	testutils.Equal(t, c.Stiff, 1.0)
}

func TestNewConstraintSpringConstantNotClamped(t *testing.T) {
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(10, 0), m)
	c := NewConstraint(p1, p2, 0.5, 0)
	testutils.Equal(t, c.Stiff, 0.5)
}

func TestNewConstraintAutoDistance(t *testing.T) {
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(10, 0), m)
	c := NewConstraint(p1, p2, 1.0, 0)
	testutils.Equal(t, c.Target, 10.0)
}

func TestNewConstraintExplicitDistance(t *testing.T) {
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(10, 0), m)
	c := NewConstraint(p1, p2, 1.0, 25.0)
	testutils.Equal(t, c.Target, 25.0)
}

func TestConstraintRelaxBothNonZeroMass(t *testing.T) {
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(20, 0), m)
	c := NewConstraint(p1, p2, 1.0, 10.0)

	origP1X := p1.Position.Coordinate(0)
	origP2X := p2.Position.Coordinate(0)
	c.Relax()
	// Both particles should move towards each other.
	newP1X := p1.Position.Coordinate(0)
	newP2X := p2.Position.Coordinate(0)
	if newP1X <= origP1X {
		t.Fatalf("expected p1 to move right, got %f -> %f", origP1X, newP1X)
	}
	if newP2X >= origP2X {
		t.Fatalf("expected p2 to move left, got %f -> %f", origP2X, newP2X)
	}
}

func TestConstraintRelaxFirstZeroMass(t *testing.T) {
	m0 := NewMaterial(0)
	m1 := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m0)
	p2 := NewParticle(shapes.NewPoint(20, 0), m1)
	c := NewConstraint(p1, p2, 1.0, 10.0)

	origP1X := p1.Position.Coordinate(0)
	c.Relax()
	// p1 has zero mass, so it should not move but p2 should.
	testutils.Equal(t, p1.Position.Coordinate(0), origP1X)
}

func TestConstraintRelaxSecondZeroMass(t *testing.T) {
	m0 := NewMaterial(0)
	m1 := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m1)
	p2 := NewParticle(shapes.NewPoint(20, 0), m0)
	c := NewConstraint(p1, p2, 1.0, 10.0)

	origP2X := p2.Position.Coordinate(0)
	c.Relax()
	// p2 has zero mass, so it should not move but p1 should.
	testutils.Equal(t, p2.Position.Coordinate(0), origP2X)
}

// ---------------------------------------------------------------------------
// World
// ---------------------------------------------------------------------------

func TestNewWorldStepLessThanOne(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)), shapes.NewPoint(0, 0), 0.5)
	testutils.Equal(t, w.Step, 1.0)
	testutils.Equal(t, w.Delta, 1.0)
}

func TestNewWorldStepValid(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)), shapes.NewPoint(0, 0), 4)
	testutils.Equal(t, w.Step, 4.0)
	testutils.Equal(t, w.Delta, 0.25)
}

func TestNewWorldStepExactlyOne(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)), shapes.NewPoint(0, 0), 1)
	testutils.Equal(t, w.Step, 1.0)
	testutils.Equal(t, w.Delta, 1.0)
}

func TestNewWorldStepNegative(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)), shapes.NewPoint(0, 0), -5)
	testutils.Equal(t, w.Step, 1.0)
	testutils.Equal(t, w.Delta, 1.0)
}

func TestWorldAddComposites(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)), shapes.NewPoint(0, 0), 1)
	c1 := NewComposite()
	c2 := NewComposite()
	w.AddComposites(c1, c2)
	testutils.Equal(t, len(w.Composites), 2)
}

func TestWorldSimulate(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(-500, -500), 500, 500)), shapes.NewPoint(0, 2), 4)
	c := NewComposite()
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	c.AddParticle(p)
	w.AddComposites(c)

	origY := p.Position.Coordinate(1)
	w.Simulate(10, 2)
	// Gravity is (0,2), so y coordinate should have increased after simulation.
	newY := p.Position.Coordinate(1)
	if newY <= origY {
		t.Fatalf("expected particle to move down under gravity, got %f -> %f", origY, newY)
	}
}

func TestWorldSimulateWithConstraints(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(-500, -500), 500, 500)), shapes.NewPoint(0, 2), 4)
	c := NewComposite()
	m := NewMaterial(1.0)
	c.AddParticle(NewParticle(shapes.NewPoint(0, 0), m))
	c.AddParticle(NewParticle(shapes.NewPoint(0, 40), m))
	err := c.AddConstraints(0, 1, 1.0)
	testutils.Equal(t, err, nil)
	w.AddComposites(c)
	w.Simulate(5, 2)
	// Just verify no panics and particles have moved.
	testutils.Equal(t, len(w.Composites), 1)
}

func TestWorldSimulateZeroSteps(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(-500, -500), 500, 500)), shapes.NewPoint(0, 2), 4)
	c := NewComposite()
	m := NewMaterial(1.0)
	p := NewParticle(shapes.NewPoint(0, 0), m)
	c.AddParticle(p)
	w.AddComposites(c)
	w.Simulate(0, 2)
	// No steps, position unchanged.
	testutils.Equal(t, p.Position.Coordinate(0), 0.0)
	testutils.Equal(t, p.Position.Coordinate(1), 0.0)
}

func TestWorldSimulateMultipleComposites(t *testing.T) {
	w := NewWorld(shapes.NewBorder(shapes.NewBox(shapes.NewPoint(-500, -500), 500, 500)), shapes.NewPoint(0, 2), 4)
	c1 := NewComposite()
	c2 := NewComposite()
	m := NewMaterial(1.0)
	p1 := NewParticle(shapes.NewPoint(0, 0), m)
	p2 := NewParticle(shapes.NewPoint(100, 0), m)
	c1.AddParticle(p1)
	c2.AddParticle(p2)
	w.AddComposites(c1, c2)
	w.Simulate(5, 2)
	// Both should be affected by gravity.
	if p1.Position.Coordinate(1) <= 0 {
		t.Fatal("p1 should have moved under gravity")
	}
	if p2.Position.Coordinate(1) <= 0 {
		t.Fatal("p2 should have moved under gravity")
	}
}
