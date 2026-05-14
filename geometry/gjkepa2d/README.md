# gjkepa2d

Import path: `github.com/InsideGallery/game-core/geometry/gjkepa2d`

Package `gjkepa2d` implements 2D collision checks using GJK and EPA over
`geometry/shapes` collision types. It can report whether two shapes intersect
and, when requested, a minimum translation vector.

Key API:

- `GJKEPA`: collision checker state for GJK/EPA.
- `NewGJKEPA()`: creates a collision checker.
- `(*GJKEPA).GJK(shapeA, shapeB, calculateMTV)`: checks intersection and
  optionally calculates an MTV through EPA.
- `(*GJKEPA).EPA(...)`: expands the final simplex to estimate penetration
  depth and direction.
- `RoundedGJKEPA(shapeA, shapeB)`: collision check with rounded MTV output and
  same-center handling.
- `SolveKinematicBody(...)`: separates colliding shapes and reflects input
  forces using elasticity values.

