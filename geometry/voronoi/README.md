# Voronoi

Import path: `github.com/InsideGallery/game-core/geometry/voronoi`

Package `voronoi` generates bounded 2D Voronoi diagrams from sites or shape
points. It implements Fortune's algorithm and stores the resulting diagram in a
DCEL structure.

Key exports:

- `New` creates a `Voronoi` generator from `SiteSlice` and bounds.
- `NewFromPoints` creates a generator from `[]shapes.Point` and bounds.
- `Voronoi.Generate` runs the algorithm.
- `Voronoi.ToPolyhedrons` converts generated cells to `shapes.Polyhedron`
  values.
- `DCEL`, `Vertex`, `Face`, and `HalfEdge` represent the diagram topology.
- `Site` and `SiteSlice` describe generator points.
- `EventQueue`, `Event`, and `Node` expose the event queue and beach-line
  structures used by the algorithm.
