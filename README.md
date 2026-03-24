# game-core

Helpers and utilities for game development in Go.

## Overview

`game-core` is a shared Go library providing reusable game-development packages for InsideGallery projects. It includes geometry algorithms, physics simulation, spatial indexing, and game engine primitives.

## Installation

```bash
go get github.com/InsideGallery/game-core
```

## Packages

| Package | Description |
|---------|-------------|
| `cards` | Card game utilities |
| `engine` | Game engine core with communications, name generation, and relationship systems |
| `geometry/astar` | A* pathfinding algorithm |
| `geometry/gjkepa2d` | GJK/EPA 2D collision detection |
| `geometry/gjkepa3d` | GJK/EPA 3D collision detection |
| `geometry/hexagone` | Hexagonal grid utilities |
| `geometry/isometric` | Isometric projection |
| `geometry/quickhull` | Quickhull convex hull algorithm |
| `geometry/shapes` | Shape primitives |
| `geometry/voronoi` | Voronoi diagram generation |
| `mathutils` | Math utilities for game development |
| `physics` | Physics simulation |
| `rtree` | R-tree spatial indexing |

## Requirements

- Go 1.24.3+
- Depends on [github.com/InsideGallery/core](https://github.com/InsideGallery/core)

## License

See [LICENSE](LICENSE) for details.
