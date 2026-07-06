# gjkepa2d

Import path: `github.com/InsideGallery/game-core/geometry/gjkepa2d`

Package `gjkepa2d` implements 2D collision checks using GJK and EPA over
`geometry/shapes` collision types. It can report whether two shapes intersect
and, when requested, a minimum translation vector.

Key API:

- `MTV(a, b)`: minimum translation vector separating `a` from `b`
  (subtracting it from `b`'s position separates the shapes) plus an overlap
  flag. Sphere-vs-{Box, Sphere, convex Polyhedron} pairs — either operand
  order — use exact closed forms; every other pair falls back to GJK+EPA.
  Deterministic tie-breaks: coincident circle centers push `b` along +x; a
  circle center inside a box exits through the nearest face (ties resolve
  west, east, north, south); a circle center inside or exactly on a polygon
  boundary exits along the nearest edge's outward normal.
- `GJKEPA`: collision checker state for GJK/EPA.
- `NewGJKEPA()`: creates a collision checker.
- `(*GJKEPA).GJK(shapeA, shapeB, calculateMTV)`: checks intersection and
  optionally calculates an MTV through EPA.
- `(*GJKEPA).EPA(...)`: expands the final simplex to estimate penetration
  depth and direction.
- `RoundedGJKEPA(shapeA, shapeB)`: collision check with rounded MTV output and
  same-center handling.
- `SolveKinematicBody(...)`: separates colliding shapes and reflects input
  forces using elasticity values.

## MTV accuracy and performance

The closed forms in `MTV` are exact (float64 rounding only). The EPA path
converges to its `1e-6` termination tolerance: MTV magnitude is accurate to
~`1e-6`, but on curved Minkowski-boundary regions (e.g. circle-vs-box corner
arcs, radius r) the direction error is bounded by `sqrt(2e-6/r)` radians —
about 0.16° for r = 0.25 — so corner MTV components can be off by
`|mtv| * 2.8e-3`. Prefer `MTV` for sphere pairs; use the EPA path when no
closed form exists.

Benchmarks (go1.x, AMD64, corner-contact circle vs box, `-benchmem`):

```
BenchmarkMTVCircleBoxAnalytic-16      6782956    174.8 ns/op     80 B/op    2 allocs/op
BenchmarkMTVCircleBoxEPA-16             17257    80967 ns/op  10120 B/op   67 allocs/op
BenchmarkMTVCirclePolygonAnalytic-16  4934608    283.9 ns/op     56 B/op    2 allocs/op
```

The analytic path is ~460x faster than EPA for the same contact and does not
allocate beyond the returned point.

