package core

import (
	"strings"
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestAstar(t *testing.T) {
	node0 := NewNode("A", nil)
	node1 := NewNode("B", nil)
	node2 := NewNode("C", nil)
	node3 := NewNode("D", nil)
	node4 := NewNode("E", nil)
	node5 := NewNode("F", nil)
	node6 := NewNode("X", nil)

	node0.AddNeighbors(node1, node3, node2)
	node1.AddNeighbors(node0, node2, node4)
	node4.AddNeighbors(node1, node2)
	node2.AddNeighbors(node0, node4, node3, node5)
	node3.AddNeighbors(node0, node2, node5)
	node5.AddNeighbors(node3, node4, node6)
	node6.AddNeighbors(node5)

	equal := func(a, b interface{}) bool {
		p1 := a.(string)
		p2 := b.(string)
		return strings.EqualFold(p1, p2)
	}
	list := Path(equal, node0, node6, func(_, _ *Node) float64 {
		return 10
	}, func(_, _ *Node) float64 {
		return 1
	})

	var result []string
	for _, item := range GetPath(list, equal, node0, true) {
		result = append(result, item.Value().(string))
	}
	testutils.Equal(t, result, []string{"A", "D", "F", "X"})
}
