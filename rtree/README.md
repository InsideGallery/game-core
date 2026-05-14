# RTree

Import path: `github.com/InsideGallery/game-core/rtree`

Package `rtree` provides an R-tree spatial index for `shapes.Spatial` objects.
It supports insertion, deletion, updates, intersection searches, collision
queries, and nearest-neighbor lookups.

Key exports:

- `NewRTree` creates an index with minimum and maximum child counts.
- `DefaultMinRTreeOption` and `DefaultMaxRTreeOption` provide default child
  count settings.
- `RTree.Insert`, `Update`, `Delete`, `Size`, and `Depth` manage indexed
  objects and tree state.
- `RTree.SearchIntersect` and `Collision` find objects whose bounds intersect a
  box or another spatial object.
- `RTree.NearestNeighbor` and `NearestNeighbors` query closest indexed objects.
- `RTree.MoveObject` updates objects that implement `Moveable`.
- `RTree.GetAllBoundingBoxes` returns stored bounding boxes.
- `Entity` allows identity-based deletion for indexed objects with IDs.
- `Moveable` combines `shapes.Spatial` with `UpdateSpatial`.
