# astar core

Import path: `github.com/InsideGallery/game-core/geometry/astar/core`

Package `core` provides an alternate A* implementation built around mutable
`Node` values and ordered sets. It supports static neighbors and dynamic
neighbor generation.

Key API:

- `Node`: graph node with a stored value, path cost state, previous node, and
  neighbor configuration.
- `NewNode(value, previous, neighbors...)`: creates a node and optionally sets
  its previous node and neighbors.
- `Node.AddNeighbors(...)`: replaces the node's static neighbors.
- `Node.AddNeighborsFn(...)`: configures dynamic neighbor discovery.
- `Path(equal, from, to, heuristic, edgeWeight)`: explores nodes from `from`
  toward `to` using the supplied heuristic and edge weight functions.
- `GetPath(...)`: reconstructs a path from the visited node set.

