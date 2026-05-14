# game-core

Helpers and utilities for game development in Go.

## Overview

`game-core` is a shared Go library providing reusable game-development packages for InsideGallery projects. It
includes geometry algorithms, physics simulation, spatial indexing, and game engine primitives.

This module is a library only. It has no `main` package and is intended to be imported by applications and other
libraries.

## Installation

```bash
go get github.com/InsideGallery/game-core
```

## Packages

| Package | Description |
|---------|-------------|
| `cards` | Card game utilities |
| `engine` | Game engine core |
| `engine/communications` | Communication components and systems |
| `engine/names` | Randomized name generation |
| `engine/relations` | Parent/child relationship components |
| `geometry/astar` | A* pathfinding algorithm |
| `geometry/astar/core` | Core A* node and path search primitives |
| `geometry/astar/core/astargeo` | Geometry-backed A* neighbor lookup |
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

Each package directory contains its own `README.md` with package-specific import paths and API notes.

## Requirements

- Go 1.26.3+
- Direct dependencies include:
  - [github.com/InsideGallery/core](https://github.com/InsideGallery/core) v1.2.1
  - [github.com/FrogoAI/memory](https://github.com/FrogoAI/memory)
  - [github.com/FrogoAI/multiproc](https://github.com/FrogoAI/multiproc)
  - [github.com/FrogoAI/set](https://github.com/FrogoAI/set)
  - [github.com/FrogoAI/testutils](https://github.com/FrogoAI/testutils)

## Verification

Run the full test and lint suite before publishing changes:

```bash
go test ./...
go test -race -count=1 ./...
golangci-lint run ./...
golangci-lint run --fix ./...
```

## License

See [LICENSE](LICENSE) for details.
