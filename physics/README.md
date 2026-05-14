# Physics

Import path: `github.com/InsideGallery/game-core/physics`

Package `physics` provides a simple particle-based simulation model. It uses
`geometry/shapes` for positions, forces, borders, and collision bounds.

Key exports:

- `World` holds gravity, border limits, composites, and simulation step data.
- `NewWorld` creates a world and clamps steps below one to a single-step delta.
- `World.Simulate` applies gravity, particle movement, border restraint, force
  reset, and constraint relaxation.
- `World.AddComposites` adds composite bodies to a world.
- `Particle` stores position, previous position, acceleration, and material.
- `Particle` methods apply acceleration, forces, impulses, restraint, and
  material updates.
- `Composite` groups particles and constraints.
- `Constraint` links two particles with a spring-like distance constraint.
- `Material` currently stores particle mass.
