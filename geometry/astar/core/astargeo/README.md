# astargeo

Import path: `github.com/InsideGallery/game-core/geometry/astar/core/astargeo`

Package `astargeo` provides geometry-oriented neighbor generation for the
`geometry/astar/core` package. It treats node values as `shapes.Point` values
and generates four cardinal neighbors at a step size of one unit.

Key API:

- `GetNeighbors(...)`: dynamic neighbor function compatible with
  `core.Node.AddNeighborsFn`. It reuses generated point nodes and skips points
  that collide with the package's internal R-tree.

