# gjkepa3d

Import path: `github.com/InsideGallery/game-core/geometry/gjkepa3d`

Package `gjkepa3d` implements 3D collision checks using GJK and EPA over
colliders backed by `geometry/shapes` points. It can report intersection and,
when requested, a minimum translation vector.

Key API:

- `Collider`: interface for 3D support-mapping collision shapes.
- `GJKEPA`: collision checker for 3D GJK/EPA.
- `NewGJKEPA()`: creates a collision checker.
- `(*GJKEPA).GJK(coll1, coll2, calculateMTV)`: checks intersection and
  optionally calculates an MTV through EPA.
- `(*GJKEPA).EPA(...)`: expands the final simplex to estimate penetration
  depth and direction.
- `RoundedGJKEPA(shapeA, shapeB)`: collision check with rounded MTV output and
  same-center handling.
- `SolveKinematicBody(...)`: separates colliding shapes and reflects input
  forces using elasticity values.

