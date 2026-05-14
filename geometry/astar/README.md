# astar

Import path: `github.com/InsideGallery/game-core/geometry/astar`

Package `astar` provides a small A* path search implementation for arbitrary
weighted graphs. Callers model graph nodes by implementing `Pather`.

Key API:

- `Pather`: interface for graph nodes that can list pathable neighbors, exact
  edge costs, and heuristic estimated costs.
- `Path(from, to Pather)`: calculates a shortest path, total distance, and a
  boolean indicating whether a path was found.

