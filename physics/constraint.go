package physics

// Constraint describe relation between two particles
type Constraint struct {
	Particle1 *Particle
	Particle2 *Particle
	Target    float64
	Stiff     float64
	Damp      float64
}

// NewConstraint return new constraint
func NewConstraint(p1, p2 *Particle, springConstant, distanceConstraint float64) *Constraint {
	if springConstant > 1 {
		springConstant = 1
	}

	if distanceConstraint == 0 {
		distanceConstraint = p2.Position.Distance(p1.Position)
	}

	return &Constraint{
		Particle1: p1,
		Particle2: p2,
		Stiff:     springConstant,
		Target:    distanceConstraint,
	}
}

// Relax add relax forces
func (c *Constraint) Relax() {
	D := c.Particle2.Position.Subtract(c.Particle1.Position)
	F := D.Normalize().Scale(0.5 * c.Stiff * (D.Normal() - c.Target))

	if c.Particle1.Material.Mass != 0 && c.Particle2.Material.Mass == 0 { //nolint:gocritic
		c.Particle1.ApplyImpulse(F.Scale(2)) // nolint:mnd
	} else if c.Particle1.Material.Mass == 0 && c.Particle2.Material.Mass != 0 {
		c.Particle2.ApplyImpulse(F.Invert().Scale(2)) // nolint:mnd
	} else {
		c.Particle1.ApplyImpulse(F)
		c.Particle2.ApplyImpulse(F.Invert())
	}
}
