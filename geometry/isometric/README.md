# Isometric

Import path: `github.com/InsideGallery/game-core/geometry/isometric`

Package `isometric` provides coordinate conversion helpers for isometric tile
maps and screens. It works with `shapes.Point` values and supports both fixed
tile-size conversions and conversions using per-axis tile dimensions.

Key exports:

- `FromOrthographic` and `ToOrthographic` convert between orthographic and
  isometric coordinates.
- `FromOrthographicWithStaticSize` and `ToOrthographicWithStaticSize` provide
  fixed tile-size variants.
- `PixelToScreenCoords`, `ScreenToPixelCoords`, and `TileToScreenCoords`
  convert between pixel, tile, and screen coordinates with height support.
- `ProjectISO`, `Nearness`, and `Closer` help project and order points in
  isometric space.
