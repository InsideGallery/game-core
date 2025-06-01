package physics

import "github.com/InsideGallery/game-core/geometry/shapes"

// World describe particle world
type World struct {
	shapes.Border
	Gravity    shapes.Point
	Composites []*Composite
	Step       float64
	Delta      float64 // Delta time (1.0 / time Step)
}

// NewWorld return new world with gravity
func NewWorld(border shapes.Border, gravity shapes.Point, step float64) *World {
	w := &World{
		Border:  border,
		Gravity: gravity,
	}

	if step < 1 {
		w.Step = 1
		w.Delta = 1
	} else {
		w.Step = step
		w.Delta = 1 / step
	}

	return w
}

// Simulate simulate physics
func (w *World) Simulate(steps, dimensions int) {
	for i := 0; i < steps; i++ {
		for _, c := range w.Composites {
			for _, p := range c.Particles {
				p.Accelerate(w.Gravity)
				p.Simulate(w.Delta)
				p.Restrain(w.Border, dimensions)
				p.ResetForces()
			}

			for _, c := range c.Constraints {
				c.Relax()
			}
		}
	}
}

// AddComposites add composite
func (w *World) AddComposites(c ...*Composite) {
	w.Composites = append(w.Composites, c...)
}
