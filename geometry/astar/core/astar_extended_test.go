package core

import (
	"strings"
	"testing"

	"github.com/FrogoAI/memory/sortedset"
	"github.com/FrogoAI/set"
	"github.com/FrogoAI/testutils"
)

func TestAddNeighborsFn(t *testing.T) {
	node0 := NewNode("A", nil)
	node1 := NewNode("B", nil)
	node2 := NewNode("C", nil)
	node3 := NewNode("D", nil)

	// Use AddNeighborsFn for node0
	node0.AddNeighborsFn(func(
		visited set.GenericOrderedDataSet[*Node],
		unvisited *sortedset.SortedSet[*Node, *Node],
		n *Node,
	) []*Node {
		_ = n
		return []*Node{node1, node2}
	})

	node1.AddNeighbors(node0, node3)
	node2.AddNeighbors(node0, node3)
	node3.AddNeighbors(node1, node2)

	equal := func(a, b interface{}) bool {
		p1 := a.(string)
		p2 := b.(string)
		return strings.EqualFold(p1, p2)
	}

	list := Path(equal, node0, node3, func(_, _ *Node) float64 {
		return 10
	}, func(_, _ *Node) float64 {
		return 1
	})

	var result []string
	for _, item := range GetPath(list, equal, node0, true) {
		result = append(result, item.Value().(string))
	}
	testutils.Equal(t, result[0], "A")
	testutils.Equal(t, result[len(result)-1], "D")
}

func TestGetPathReversed(t *testing.T) {
	node0 := NewNode("A", nil)
	node1 := NewNode("B", nil)
	node2 := NewNode("C", nil)

	node0.AddNeighbors(node1)
	node1.AddNeighbors(node0, node2)
	node2.AddNeighbors(node1)

	equal := func(a, b interface{}) bool {
		return strings.EqualFold(a.(string), b.(string))
	}

	list := Path(equal, node0, node2, func(_, _ *Node) float64 {
		return 10
	}, func(_, _ *Node) float64 {
		return 1
	})

	// Test with reverse=true
	pathForward := GetPath(list, equal, node0, true)
	testutils.Equal(t, pathForward[0].Value().(string), "A")
	testutils.Equal(t, pathForward[len(pathForward)-1].Value().(string), "C")
}

func TestGetPathNotReversed(t *testing.T) {
	node0 := NewNode("A", nil)
	node1 := NewNode("B", nil)
	node2 := NewNode("C", nil)

	node0.AddNeighbors(node1)
	node1.AddNeighbors(node0, node2)
	node2.AddNeighbors(node1)

	equal := func(a, b interface{}) bool {
		return strings.EqualFold(a.(string), b.(string))
	}

	list := Path(equal, node0, node2, func(_, _ *Node) float64 {
		return 10
	}, func(_, _ *Node) float64 {
		return 1
	})

	// Test with reverse=false
	pathBackward := GetPath(list, equal, node0, false)
	testutils.Equal(t, pathBackward[0].Value().(string), "C")
	testutils.Equal(t, pathBackward[len(pathBackward)-1].Value().(string), "A")
}

func TestNodeValue(t *testing.T) {
	n := NewNode(42, nil)
	testutils.Equal(t, n.Value(), 42)
}

func TestNodeCost(t *testing.T) {
	n := NewNode("A", nil)
	n.Cost(5.0, 3.0)
	testutils.Equal(t, n.gCost, 5.0)
	testutils.Equal(t, n.fCost, 8.0)
}

func TestPathNoEqual(t *testing.T) {
	// Test Path with nil equal function - should traverse all nodes
	node0 := NewNode("A", nil)
	node1 := NewNode("B", nil)
	node2 := NewNode("C", nil)

	node0.AddNeighbors(node1)
	node1.AddNeighbors(node0, node2)
	node2.AddNeighbors(node1)

	list := Path(nil, node0, node2, func(_, _ *Node) float64 {
		return 10
	}, func(_, _ *Node) float64 {
		return 1
	})

	// With nil equal, it should visit all reachable nodes
	testutils.Equal(t, list.Count() > 0, true)
}

func TestAddNeighborsFnWithExistingUnvisitedNodes(t *testing.T) {
	node0 := NewNode("Start", nil)
	node1 := NewNode("Mid", nil)
	node2 := NewNode("End", nil)

	// Track if neighborsFn is called
	called := false
	node0.AddNeighborsFn(func(
		visited set.GenericOrderedDataSet[*Node],
		unvisited *sortedset.SortedSet[*Node, *Node],
		n *Node,
	) []*Node {
		called = true
		return []*Node{node1}
	})
	node1.AddNeighbors(node2)
	node2.AddNeighbors(node1)

	equal := func(a, b interface{}) bool {
		return a.(string) == b.(string)
	}

	Path(equal, node0, node2, func(_, _ *Node) float64 {
		return 1
	}, func(_, _ *Node) float64 {
		return 1
	})

	testutils.Equal(t, called, true)
}
