package astargeo

import (
	"github.com/InsideGallery/core/memory/set"
	"github.com/InsideGallery/core/memory/sortedset"
	"github.com/InsideGallery/game-core/geometry/astar/core"
	"github.com/InsideGallery/game-core/geometry/shapes"
	rtree2 "github.com/InsideGallery/game-core/rtree"
)

var (
	points = map[int]map[int]*core.Node{}
	rtree  = rtree2.NewRTree(25, 50) //nolint:mnd
)

// GetNeighbors return Neighbors
func GetNeighbors(
	_ set.GenericOrderedDataSet[*core.Node],
	_ *sortedset.SortedSet[*core.Node, *core.Node],
	from *core.Node,
) (nodes []*core.Node) {
	pos := from.Value().(shapes.Point)
	var step float64 = 1

	for i := 0; i < 4; i++ {
		switch i {
		case 0: //nolint:mnd
			pos = shapes.NewPoint(pos.Coordinate(0)-step, pos.Coordinate(1))
		case 1: //nolint:mnd
			pos = shapes.NewPoint(pos.Coordinate(0)+step, pos.Coordinate(1))
		case 2: //nolint:mnd
			pos = shapes.NewPoint(pos.Coordinate(0), pos.Coordinate(1)-step)
		case 3: //nolint:mnd
			pos = shapes.NewPoint(pos.Coordinate(0), pos.Coordinate(1)+step)
		}

		x, y := pos.Coordinate(0), pos.Coordinate(1)

		if obj := rtree.Collision(pos, nil); len(obj) > 0 {
			continue
		}

		if node, ok := points[int(x)][int(y)]; ok {
			nodes = append(nodes, node)
		} else {
			n := core.NewNode(pos, nil)

			if _, ok := points[int(x)]; !ok {
				points[int(x)] = map[int]*core.Node{}
			}

			points[int(x)][int(y)] = n
			n.AddNeighborsFn(GetNeighbors)
			nodes = append(nodes, n)
		}
	}

	return nodes
}
