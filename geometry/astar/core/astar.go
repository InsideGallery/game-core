package core

import (
	"math"

	"github.com/InsideGallery/core/memory/comparator"
	"github.com/InsideGallery/core/memory/set"
	"github.com/InsideGallery/core/memory/sortedset"
)

// https://isaaccomputerscience.org/concepts/dsa_search_a_star

// Heuristic function to calculate cost
type Heuristic func(node1, node2 *Node) float64

// Node node description
type Node struct {
	value       interface{}
	gCost       float64
	fCost       float64
	previous    *Node
	neighbors   []*Node
	neighborsFn func(set.GenericOrderedDataSet[*Node], *sortedset.SortedSet[*Node, *Node], *Node) []*Node
}

// NewNode return new node
func NewNode(value interface{}, previous *Node, neighbors ...*Node) *Node {
	node := &Node{
		previous: previous,
		value:    value,
	}
	node.AddNeighbors(neighbors...)

	return node
}

func (n *Node) AddNeighbors(neighbors ...*Node) *Node {
	for _, ng := range neighbors {
		ng.gCost = math.MaxFloat64
		ng.fCost = math.MaxFloat64
		ng.previous = nil
	}

	n.neighbors = neighbors

	return n
}

func (n *Node) AddNeighborsFn(
	neighbors func(set.GenericOrderedDataSet[*Node],
		*sortedset.SortedSet[*Node, *Node],
		*Node,
	) []*Node,
) *Node {
	n.neighborsFn = func(
		visited set.GenericOrderedDataSet[*Node],
		unvisited *sortedset.SortedSet[*Node, *Node],
		node *Node,
	) []*Node {
		nodes := neighbors(visited, unvisited, node)

		for _, ng := range nodes {
			if !unvisited.Contains(ng) {
				ng.gCost = math.MaxFloat64
				ng.fCost = math.MaxFloat64
				ng.previous = nil
			}
		}

		return nodes
	}

	return n
}

// Neighbors return graph neighbors
func (n *Node) Neighbors(visited set.GenericOrderedDataSet[*Node], unvisited *sortedset.SortedSet[*Node, *Node]) []*Node {
	if n.neighborsFn != nil {
		return n.neighborsFn(visited, unvisited, n)
	}

	return n.neighbors
}

// Value return node value
func (n *Node) Value() interface{} {
	return n.value
}

// Cost calculate cost of node
func (n *Node) Cost(gCost, heuristic float64) {
	n.gCost = gCost
	n.fCost = n.gCost + heuristic
}

func Path(
	equal func(a, b interface{}) bool,
	from *Node, to *Node, heuristic Heuristic, edgeWeight Heuristic,
) set.GenericOrderedDataSet[*Node] {
	from.Cost(0, heuristic(from, to))
	unvisited := sortedset.NewSortedSet[*Node, *Node](func(a, b interface{}) int {
		v1 := a.(*Node)
		v2 := b.(*Node)

		return comparator.Float64Comparator(v1.fCost, v2.fCost)
	})
	unvisited.Upsert(from, from)
	visited := set.NewGenericOrderedDataSet[*Node]()

	for unvisited.GetCount() > 0 {
		res := unvisited.PeekMin().Value()
		if equal != nil && equal(res.Value(), to.Value()) {
			visited.Add(res)
			break
		}

		neighbors := res.Neighbors(visited, unvisited)
		for _, n := range neighbors {
			if visited.Contains(n) {
				continue
			}

			newGScore := res.gCost + edgeWeight(res, n)
			if newGScore < n.gCost {
				n.Cost(newGScore, heuristic(n, to))
				n.previous = res
				unvisited.Upsert(n, n)
			}
		}

		visited.Add(res)
		unvisited.Remove(res)
	}

	return visited
}
