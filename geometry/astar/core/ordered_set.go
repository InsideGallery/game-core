package core

import (
	"github.com/InsideGallery/core/memory/set"
)

// GetPath return path
func GetPath(s set.GenericOrderedDataSet[*Node], _ func(a, b interface{}) bool, _ *Node, r bool) []*Node {
	var result []*Node
	node := s.Last()

	for node != nil {
		result = append(result, node)
		node = node.previous
	}

	if r {
		reverse(result)
	}

	return result
}

func reverse(numbers []*Node) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}
