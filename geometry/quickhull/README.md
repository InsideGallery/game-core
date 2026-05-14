# Quickhull

Import path: `github.com/InsideGallery/game-core/geometry/quickhull`

Package `quickhull` builds 3D convex hulls from point clouds using the
Quickhull algorithm. Inputs and outputs use `shapes.Point` values from the
geometry shapes package.

Key exports:

- `QuickHull` runs hull generation.
- `QuickHull.ConvexHull` returns a `ConvexHull` with vertices and triangle
  indices. It can return clockwise or counter-clockwise triangles and can use
  original point-cloud indices.
- `QuickHull.ConvexHullAsMesh` returns a half-edge mesh representation.
- `ConvexHull.Triangles` expands indexed hull data into triangle point triples.
- `HalfEdgeMesh`, `HalfEdge`, and `Face` describe the mesh output.
- `FastRand` and `FastRandFloat64` provide package-local random helpers.
