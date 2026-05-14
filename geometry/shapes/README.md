# Shapes

Import path: `github.com/InsideGallery/game-core/geometry/shapes`

Package `shapes` contains shared geometric primitives and interfaces used by
the geometry, physics, and spatial-indexing packages. Most primitives implement
`Spatial`, which exposes bounding boxes, movement, centers, and support points.

Key exports:

- `Point` provides 3D coordinate math such as distance, dot product, cross
  product, normalization, rounding, min/max, and interpolation.
- `Box`, `Sphere`, `Line`, `Triangle`, `Polyhedron`, `Ellipsoid`, and
  `MultiObject` represent common spatial shapes.
- `Border` wraps a box and reports when another spatial object is outside it.
- `Rotatable` and `RotatableBox` support rotated spatial objects.
- `Spatial` is the common interface for boundable and movable objects.
- `Collide` is implemented by objects that can report collisions.
- `DegreesToRadian`, `RadianToDegree`, `NormalizeDegrees`, `GetAngle2D`,
  `GetDiffPoint2D`, `RotatePoint`, and `RotateBy` provide angle and rotation
  helpers.
