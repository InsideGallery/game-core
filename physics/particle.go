package physics

import (
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// Particle describe particle
type Particle struct {
	Material Material

	Position     shapes.Point
	Previous     shapes.Point
	Velocity     shapes.Point
	Acceleration shapes.Point
}

// NewParticle create new particle
func NewParticle(position shapes.Point, material Material) *Particle {
	return &Particle{
		Position: position,
		Previous: position,
		Material: material,
	}
}

// Simulate simulate world
func (p *Particle) Simulate(worldDelta float64) {
	if p.Material.Mass == 0 {
		return
	}

	p.Velocity = p.Position.Scale(2).Subtract(p.Previous) // nolint:mnd
	p.Previous = p.Position
	p.Position = p.Velocity.Add(p.Acceleration.Scale(worldDelta * worldDelta)) //nolint:mnd
	p.Velocity = p.Position.Subtract(p.Previous)
	p.Acceleration = shapes.NewPoint(0, 0, 0)
}

// Accelerate apply acceleration
func (p *Particle) Accelerate(rate shapes.Point) {
	p.Acceleration = p.Acceleration.Add(rate)
}

// ApplyForce apply force
func (p *Particle) ApplyForce(force shapes.Point) {
	if p.Material.Mass == 0 {
		return
	}
	p.Acceleration = p.Acceleration.Add(force.Scale(1 / p.Material.Mass))
}

// ApplyImpulse immediately change position
func (p *Particle) ApplyImpulse(impulse shapes.Point) {
	if p.Material.Mass == 0 {
		return
	}
	p.Position = p.Position.Add(impulse.Scale(1 / p.Material.Mass))
}

// ResetForces reset acceleration
func (p *Particle) ResetForces() {
	p.Acceleration = shapes.NewPoint(0, 0, 0)
}

// Restrain return particle into border on collision
func (p *Particle) Restrain(border shapes.Border, dimensions int) {
	_, points := border.Collision(p.Position, dimensions)
	p.Position = shapes.NewPoint(points...)
}

// SetMaterial set material
func (p *Particle) SetMaterial(material Material) {
	p.Material = material
}
