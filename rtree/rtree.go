package rtree

import (
	"log"
	"math"
	"sort"
	"sync"

	"github.com/InsideGallery/game-core/geometry/shapes"
)

// RTree options
const (
	DefaultMinRTreeOption = 25
	DefaultMaxRTreeOption = 50
)

// RTree represents an R-tree, a balanced search tree for storing and querying shapes.Spatial objects
type RTree struct {
	MinChildren int
	MaxChildren int
	root        *node
	size        int
	height      int

	mu sync.RWMutex
}

// NewRTree returns an RTree. If the number of objects given on initialization
// is larger than max, the RTree will be initialized using the Overlap
// Minimizing Top-down bulk-loading algorithm.
func NewRTree(min, max int) *RTree {
	if min < 2 { //nolint:mnd
		min = 2 //nolint:mnd
	}

	if min > max/2 { //nolint:mnd
		max = min * 2 //nolint:mnd
	}

	rt := &RTree{
		MinChildren: min,
		MaxChildren: max,
		height:      1,
		root: &node{
			entries: []entry{},
			leaf:    true,
			level:   1,
		},
	}

	return rt
}

// Size returns the number of objects currently stored in tree.
func (r *RTree) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s := r.size

	return s
}

// Depth returns the maximum depth of tree.
func (r *RTree) Depth() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h := r.height

	return h
}

// node represents a tree node of an RTree.
type node struct {
	parent  *node
	leaf    bool
	entries []entry
	level   int // node depth in the RTree
}

// entry represents a shapes.Spatial entry record stored in a tree node.
type entry struct {
	bb    shapes.Box // bounding-box of all children of this entry
	child *node
	obj   shapes.Spatial
}

// Entity return GetID
type Entity interface {
	GetID() uint32
}

// Insert implemented per Section 3.2 of
// "R-trees: A Dynamic Index Structure for shapes.Spatial Searching" by A. Guttman,
// Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (r *RTree) Insert(obj shapes.Spatial) {
	r.mu.Lock()
	defer r.mu.Unlock()

	e := entry{obj.Bounds(), nil, obj}
	r.insert(e, 1)

	r.size++
}

// Update delete and insert object
func (r *RTree) Update(obj shapes.Spatial) {
	r.Delete(obj)
	r.Insert(obj)
}

// Collision check shapes.Spatial object on collisions
func (r *RTree) Collision(obj shapes.Spatial, filter func(spatial shapes.Spatial) bool) (objects []shapes.Spatial) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entities := r.searchIntersect([]shapes.Spatial{}, r.root, obj.Bounds(), filter)
	objects = append(objects, entities...)

	return
}

// insert adds the specified entry to the tree at the specified level.
func (r *RTree) insert(e entry, level int) {
	leaf := r.chooseNode(r.root, e, level)
	leaf.entries = append(leaf.entries, e)

	// update parent pointer if necessary
	if e.child != nil {
		e.child.parent = leaf
	}

	// split leaf if overflows
	var split *node
	if len(leaf.entries) > r.MaxChildren {
		leaf, split = leaf.split(r.MinChildren)
	}
	root, splitRoot := r.adjustTree(leaf, split)

	if splitRoot != nil {
		oldRoot := root
		r.height++
		r.root = &node{
			parent: nil,
			level:  r.height,
			entries: []entry{
				{bb: oldRoot.computeBoundingBox(), child: oldRoot},
				{bb: splitRoot.computeBoundingBox(), child: splitRoot},
			},
		}
		oldRoot.parent = r.root
		splitRoot.parent = r.root
	}
}

// chooseNode finds the node at the specified level to which e should be added.
func (r *RTree) chooseNode(n *node, e entry, level int) *node {
	if n.leaf || n.level == level {
		return n
	}

	// find the entry whose bb needs least enlargement to include obj
	diff := math.MaxFloat64
	var chosen entry

	for _, en := range n.entries {
		bb := en.bb.BoundingBox(e.bb)
		d := bb.Volume() - en.bb.Volume()

		if d < diff || (d == diff && en.bb.Volume() < chosen.bb.Volume()) {
			diff = d
			chosen = en
		}
	}

	return r.chooseNode(chosen.child, e, level)
}

// adjustTree splits overflowing nodes and propagates the changes upwards.
func (r *RTree) adjustTree(n, nn *node) (*node, *node) {
	// Let the caller handle root adjustments.
	if n == r.root {
		return n, nn
	}

	// Re-size the bounding box of n to account for lower-level changes.
	en := n.getEntry()
	en.bb = n.computeBoundingBox()

	// If nn is nil, then we're just propagating changes upwards.
	if nn == nil {
		return r.adjustTree(n.parent, nil)
	}

	// Otherwise, these are two nodes resulting from a split.
	// n was reused as the "left" node, but we need to add nn to n.parent.
	enn := entry{nn.computeBoundingBox(), nn, nil}
	n.parent.entries = append(n.parent.entries, enn)

	// If the new entry overflows the parent, split the parent and propagate.
	if len(n.parent.entries) > r.MaxChildren {
		return r.adjustTree(n.parent.split(r.MinChildren))
	}

	// Otherwise keep propagating changes upwards.
	return r.adjustTree(n.parent, nil)
}

// getEntry returns a pointer to the entry for the node n from n's parent.
func (n *node) getEntry() *entry {
	var e *entry

	for i := range n.parent.entries {
		if n.parent.entries[i].child == n {
			e = &n.parent.entries[i]
			break
		}
	}

	return e
}

// computeBoundingBox finds the MBR of the children of n.
func (n *node) computeBoundingBox() (bb shapes.Box) {
	for _, e := range n.entries {
		bb = bb.BoundingBox(e.bb)
	}

	return
}

// split splits a node into two groups while attempting to minimize the
// bounding-box area of the resulting groups.
func (n *node) split(minGroupSize int) (left, right *node) {
	// find the initial split
	l, r := n.pickSeeds()
	leftSeed, rightSeed := n.entries[l], n.entries[r]

	// get the entries to be divided between left and right
	remaining := append(n.entries[:l], n.entries[l+1:r]...) //nolint:gocritic
	remaining = append(remaining, n.entries[r+1:]...)

	// setup the new split nodes, but re-use n as the left node
	left = n
	left.entries = []entry{leftSeed}
	right = &node{
		parent:  n.parent,
		leaf:    n.leaf,
		level:   n.level,
		entries: []entry{rightSeed},
	}

	if rightSeed.child != nil {
		rightSeed.child.parent = right
	}

	if leftSeed.child != nil {
		leftSeed.child.parent = left
	}

	// distribute all of n's old entries into left and right.
	for len(remaining) > 0 {
		next := pickNext(left, right, remaining)
		e := remaining[next]

		if len(remaining)+len(left.entries) <= minGroupSize { //nolint:gocritic
			assign(e, left)
		} else if len(remaining)+len(right.entries) <= minGroupSize {
			assign(e, right)
		} else {
			assignGroup(e, left, right)
		}

		remaining = append(remaining[:next], remaining[next+1:]...)
	}

	return
}

// getAllBoundingBoxes traverses tree populating slice of bounding boxes of non-leaf nodes.
func (n *node) getAllBoundingBoxes() []shapes.Box {
	var rects []shapes.Box

	if n.leaf {
		return rects
	}

	for _, e := range n.entries {
		if e.child == nil {
			return rects
		}
		rectsInter := append(e.child.getAllBoundingBoxes(), e.bb)
		rects = append(rects, rectsInter...)
	}

	return rects
}

func assign(e entry, group *node) {
	if e.child != nil {
		e.child.parent = group
	}
	group.entries = append(group.entries, e)
}

// assignGroup chooses one of two groups to which a node should be added.
func assignGroup(e entry, left, right *node) {
	leftBB := left.computeBoundingBox()
	rightBB := right.computeBoundingBox()
	leftEnlarged := leftBB.BoundingBox(e.bb)
	rightEnlarged := rightBB.BoundingBox(e.bb)

	// first, choose the group that needs the least enlargement
	leftDiff := leftEnlarged.Volume() - leftBB.Volume()
	rightDiff := rightEnlarged.Volume() - rightBB.Volume()

	if diff := leftDiff - rightDiff; diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group that has smaller area
	if diff := leftBB.Volume() - rightBB.Volume(); diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group with fewer entries
	if diff := len(left.entries) - len(right.entries); diff <= 0 {
		assign(e, left)
		return
	}

	assign(e, right)
}

// pickSeeds chooses two child entries of n to start a split.
func (n *node) pickSeeds() (int, int) {
	left, right := 0, 1
	maxWastedSpace := -1.0

	for i, e1 := range n.entries {
		for j, e2 := range n.entries[i+1:] {
			d := e1.bb.BoundingBox(e2.bb).Volume() - e1.bb.Volume() - e2.bb.Volume()
			if d > maxWastedSpace {
				maxWastedSpace = d
				left, right = i, j+i+1
			}
		}
	}

	return left, right
}

// pickNext chooses an entry to be added to an entry group.
func pickNext(left, right *node, entries []entry) (next int) {
	maxDiff := -1.0
	leftBB := left.computeBoundingBox()
	rightBB := right.computeBoundingBox()

	for i, e := range entries {
		d1 := leftBB.BoundingBox(e.bb).Volume() - leftBB.Volume()
		d2 := rightBB.BoundingBox(e.bb).Volume() - rightBB.Volume()
		d := math.Abs(d1 - d2)

		if d > maxDiff {
			maxDiff = d
			next = i
		}
	}

	return
}

// MoveObject move object
func (r *RTree) MoveObject(obj Moveable, v shapes.Point) {
	r.Delete(obj)
	obj.UpdateSpatial(obj.Move(v))
	r.Insert(obj)
}

// Delete removes an object from the tree
func (r *RTree) Delete(obj shapes.Spatial) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	n := r.findLeaf(r.root, obj)
	if n == nil {
		return false
	}

	ind := -1

	for i, e := range n.entries {
		if obj == e.obj {
			ind = i
			break
		}

		s, ok := e.obj.(Entity)
		s2, ok2 := obj.(Entity)

		if !ok || !ok2 {
			continue
		}

		if s.GetID() == s2.GetID() {
			ind = i
			break
		}
	}

	if ind < 0 {
		return false
	}

	n.entries = append(n.entries[:ind], n.entries[ind+1:]...)

	r.condenseTree(n)

	r.size--

	if !r.root.leaf && len(r.root.entries) == 1 {
		r.root = r.root.entries[0].child
	}

	r.height = r.root.level

	return true
}

// findLeaf finds the leaf node containing obj.
func (r *RTree) findLeaf(n *node, obj shapes.Spatial) *node {
	if n.leaf {
		return n
	}

	// if not leaf, search all candidate subtrees
	for _, e := range n.entries {
		if e.bb.ContainsRectangle(obj.Bounds()) {
			leaf := r.findLeaf(e.child, obj)
			if leaf == nil {
				continue
			}

			// check if the leaf actually contains the object
			for _, leafEntry := range leaf.entries {
				if leafEntry.obj == obj {
					return leaf
				}

				s, ok := leafEntry.obj.(Entity)
				if !ok {
					continue
				}

				s2, ok := obj.(Entity)
				if !ok {
					continue
				}

				if s.GetID() == s2.GetID() {
					return leaf
				}
			}
		}
	}

	return nil
}

// condenseTree deletes underflowing nodes and propagates the changes upwards.
func (r *RTree) condenseTree(n *node) {
	var deleted []*node

	for n != r.root {
		if len(n.entries) < r.MinChildren {
			// remove n from parent entries
			var entries []entry
			for _, e := range n.parent.entries {
				if e.child != n {
					entries = append(entries, e)
				}
			}

			if len(n.parent.entries) == len(entries) {
				log.Println("Incorrect entries size")
			}

			n.parent.entries = entries

			// only add n to deleted if it still has children
			if len(n.entries) > 0 {
				deleted = append(deleted, n)
			}
		} else {
			// just a child entry deletion, no underflow
			n.getEntry().bb = n.computeBoundingBox()
		}
		n = n.parent
	}

	for _, n := range deleted {
		// reinsert entry so that it will remain at the same level as before
		e := entry{n.computeBoundingBox(), n, nil}
		r.insert(e, n.level+1)
	}
}

// SearchIntersect returns all objects that intersect the specified rectangle.
func (r *RTree) SearchIntersect(bb shapes.Box, filter func(spatial shapes.Spatial) bool) []shapes.Spatial {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := r.searchIntersect([]shapes.Spatial{}, r.root, bb, filter)

	return result
}

func (r *RTree) searchIntersect(results []shapes.Spatial, n *node, bb shapes.Box,
	filter func(spatial shapes.Spatial) bool,
) []shapes.Spatial {
	for _, e := range n.entries {
		if _, r := e.bb.Intersect(bb); !r {
			continue
		}

		if !n.leaf {
			results = r.searchIntersect(results, e.child, bb, filter)
			continue
		}

		if filter != nil && filter(e.obj) {
			continue
		}

		results = append(results, e.obj)
	}

	return results
}

// GetAllBoundingBoxes returning slice of bounding boxes by traversing tree.
func (r *RTree) GetAllBoundingBoxes() []shapes.Box {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var rects []shapes.Box

	if r.root != nil {
		rects = r.root.getAllBoundingBoxes()
	}

	return rects
}

// NearestNeighbor returns the closest object to the specified point
func (r *RTree) NearestNeighbor(p shapes.Point, filter func(spatial shapes.Spatial) bool) (shapes.Spatial, float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.nearestNeighbor(p, r.root, math.MaxFloat64, nil, filter)
}

// utilities for sorting slices of entries
type entrySlice struct {
	entries   []entry
	distances []float64
}

// Len return len of entries
func (s entrySlice) Len() int { return len(s.entries) }

// Swap swap for sort
func (s entrySlice) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
	s.distances[i], s.distances[j] = s.distances[j], s.distances[i]
}

// Less return true if i less of j
func (s entrySlice) Less(i, j int) bool {
	return s.distances[i] < s.distances[j]
}

func sortEntries(p shapes.Point, entries []entry) ([]entry, []float64) {
	sorted := make([]entry, len(entries))
	distances := make([]float64, len(entries))

	return sortPreselectedEntries(p, entries, sorted, distances)
}

func sortPreselectedEntries(p shapes.Point, entries, sorted []entry, distances []float64) ([]entry, []float64) {
	sorted = sorted[:len(entries)]
	distances = distances[:len(entries)]

	for i := 0; i < len(entries); i++ {
		sorted[i] = entries[i]
		distances[i] = p.MinDistance(entries[i].bb)
	}

	sort.Sort(entrySlice{sorted, distances})

	return sorted, distances
}

func pruneEntries(p shapes.Point, entries []entry, minDistances []float64) []entry {
	minMinMaxDist := math.MaxFloat64

	for i := range entries {
		minMaxDist := p.MinMaxDistance(entries[i].bb)
		if minMaxDist < minMinMaxDist {
			minMinMaxDist = minMaxDist
		}
	}

	// remove all entries with minDist > minMinMaxDist
	var pruned []entry

	for i := range entries {
		if minDistances[i] <= minMinMaxDist {
			pruned = append(pruned, entries[i])
		}
	}

	return pruned
}

func pruneEntriesMinDist(d float64, entries []entry, minDistances []float64) []entry {
	var i int

	for ; i < len(entries); i++ {
		if minDistances[i] > d {
			break
		}
	}

	return entries[:i]
}

func (r *RTree) nearestNeighbor(p shapes.Point, n *node, d float64,
	nearest shapes.Spatial, filter func(spatial shapes.Spatial) bool,
) (shapes.Spatial, float64) {
	if n.leaf {
		for _, e := range n.entries {
			if filter != nil && filter(e.obj) {
				continue
			}

			dist := math.Sqrt(p.MinDistance(e.bb))
			if dist < d {
				d = dist
				nearest = e.obj
			}
		}
	} else {
		branches, distances := sortEntries(p, n.entries)
		branches = pruneEntries(p, branches, distances)

		for _, e := range branches {
			subNearest, dist := r.nearestNeighbor(p, e.child, d, nearest, filter)
			if dist < d {
				d = dist
				nearest = subNearest
			}
		}
	}

	return nearest, d
}

// NearestNeighbors gets the closest Spatials to the Point.
func (r *RTree) NearestNeighbors(k int, p shapes.Point,
	maxDistance float64, filter func(spatial shapes.Spatial) bool,
) ([]shapes.Spatial, []float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// preallocate the buffers for sortings the branches. At each level of the
	// tree, we slide the buffer by the number of entries in the node.
	maxBufferSize := r.MaxChildren * r.Depth()
	branches := make([]entry, maxBufferSize)
	branchDistances := make([]float64, maxBufferSize)

	// allocate the buffers for the results
	distances := make([]float64, 0, k)
	objs := make([]shapes.Spatial, 0, k)

	return r.nearestNeighbors(k, p, r.root, distances, objs, branches, branchDistances, maxDistance, filter)
}

// insert obj into nearest and return the first k elements in increasing order.
func insertNearest(k int, distances []float64, nearest []shapes.Spatial,
	distance float64, obj shapes.Spatial, maxDistance float64,
) ([]float64, []shapes.Spatial) {
	i := sort.SearchFloat64s(distances, distance)
	for i < len(nearest) && distance >= distances[i] {
		i++
	}

	if i >= k {
		return distances, nearest
	}

	if distance > maxDistance {
		return distances, nearest
	}

	// no resize since cap = k
	if len(nearest) < k {
		distances = append(distances, 0)
		nearest = append(nearest, nil)
	}

	left, right := distances[:i], distances[i:len(distances)-1]
	copy(distances, left)
	copy(distances[i+1:], right)
	distances[i] = distance

	leftObjs, rightObjs := nearest[:i], nearest[i:len(nearest)-1]
	copy(nearest, leftObjs)
	copy(nearest[i+1:], rightObjs)
	nearest[i] = obj

	return distances, nearest
}

func (r *RTree) nearestNeighbors(
	k int, p shapes.Point, n *node,
	distances []float64, nearest []shapes.Spatial, b []entry,
	bd []float64, maxDistance float64,
	filter func(spatial shapes.Spatial) bool,
) ([]shapes.Spatial, []float64) {
	if n.leaf {
		for _, e := range n.entries {
			if filter != nil && filter(e.obj) {
				continue
			}

			distance := p.MinDistance(e.bb)
			distance = math.Sqrt(distance)
			distances, nearest = insertNearest(k, distances, nearest, distance, e.obj, maxDistance)
		}
	} else {
		branches, branchDists := sortPreselectedEntries(p, n.entries, b, bd)
		// only prune if buffer has k elements
		if l := len(distances); l >= k {
			branches = pruneEntriesMinDist(distances[l-1], branches, branchDists)
		}

		for _, e := range branches {
			nearest, distances = r.nearestNeighbors(k, p, e.child, distances,
				nearest, b[len(n.entries):], bd[len(n.entries):], maxDistance, filter)
		}
	}

	return nearest, distances
}
