# hexagone

Import path: `github.com/InsideGallery/game-core/geometry/hexagone`

Package `hexagone` provides axial and cube coordinate helpers for hexagonal
grids.

Key API:

- `Axial`: axial hex coordinate with `Col` and `Row`.
- `NewAxial(col, row)`: creates an axial coordinate.
- `Axial.ToCube()`: converts axial coordinates to cube coordinates.
- `Axial.ToPosition()`: converts axial coordinates to a 2D point position.
- `Axial.Distance(...)`: returns hex distance between axial coordinates.
- `Cube`: cube hex coordinate with `X`, `Y`, and `Z`.
- `NewCube(x, y, z)`: creates a cube coordinate.
- `Cube.ToAxis()`: converts cube coordinates to axial coordinates.
- `Cube.Distance(...)`: returns hex distance between cube coordinates.

