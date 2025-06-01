package rtree

import "github.com/InsideGallery/game-core/geometry/shapes"

// Moveable objects what ability to update themself
type Moveable interface {
	shapes.Spatial
	UpdateSpatial(s shapes.Spatial)
}
