package physics

// Composite describe complex object with multiple particles and constraints
type Composite struct {
	Particles   []*Particle
	Constraints []*Constraint
}

// NewComposite create new composite
func NewComposite() *Composite {
	return &Composite{}
}

// AddParticle add particle
func (c *Composite) AddParticle(p *Particle) {
	c.Particles = append(c.Particles, p)
}

// GetParticle get particle by index
func (c *Composite) GetParticle(i int) (*Particle, error) {
	if i > len(c.Particles) || i < 0 {
		return nil, ErrNotFoundAttachedParticle
	}

	return c.Particles[i], nil
}

// AddConstraints create constraint between two particles
func (c *Composite) AddConstraints(i1, i2 int, springConstant float64) error {
	p1, err := c.GetParticle(i1)
	if err != nil {
		return err
	}

	p2, err := c.GetParticle(i2)
	if err != nil {
		return err
	}

	c.Constraints = append(c.Constraints, NewConstraint(p1, p2, springConstant, 0))

	return nil
}

// SetMaterial set material to all particles
func (c *Composite) SetMaterial(material Material) {
	for _, p := range c.Particles {
		p.SetMaterial(material)
	}
}
