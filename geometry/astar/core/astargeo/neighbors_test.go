package astargeo

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
	"github.com/InsideGallery/game-core/geometry/astar/core"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

func TestGetNeighbors(t *testing.T) {
	xs, ys := 8, 12
	xd, yd := 2, 1

	xi, yi, wi, hi := 2, 8, 1, 7
	rtree.Insert(shapes.NewBox(shapes.NewPoint(float64(xi), float64(yi)), float64(wi), float64(hi)))
	xi, yi, wi, hi = 8, 1, 3, 6
	rtree.Insert(shapes.NewBox(shapes.NewPoint(float64(xi), float64(yi)), float64(wi), float64(hi)))

	to := core.NewNode(shapes.NewPoint(float64(xd), float64(yd)), nil)

	fromPos := shapes.NewPoint(float64(xs), float64(ys))
	from := core.NewNode(fromPos, nil)
	from.AddNeighborsFn(GetNeighbors)

	equal := func(a, b interface{}) bool {
		p1 := a.(shapes.Point)
		p2 := b.(shapes.Point)
		return p1.Equals(p2)
	}
	nodes := core.Path(equal, from, to, func(node1, node2 *core.Node) float64 {
		p1 := node1.Value().(shapes.Point)
		p2 := node2.Value().(shapes.Point)
		return p1.ManhattanDistance(p2)
	}, func(node1, node2 *core.Node) float64 {
		p1 := node1.Value().(shapes.Point)
		p2 := node2.Value().(shapes.Point)
		return p1.ManhattanDistance(p2)
	})

	var i int
	for range core.GetPath(nodes, equal, from, true) {
		i++
	}
	testutils.Equal(t, i, 17)
}
