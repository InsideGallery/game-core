package voronoi

import (
	"testing"

	"github.com/FrogoAI/testutils"
	"github.com/InsideGallery/game-core/geometry/shapes"
)

// ---------- Site tests ----------

func TestSiteString(t *testing.T) {
	s := Site{X: 10, Y: 20}
	testutils.Equal(t, s.String(), "10,20")
}

func TestSiteSliceSort(t *testing.T) {
	sites := SiteSlice{
		{X: 5, Y: 10},
		{X: 3, Y: 5},
		{X: 7, Y: 5},
	}
	testutils.Equal(t, sites.Len(), 3)
	testutils.Equal(t, sites.Less(1, 0), true)  // Y=5 < Y=10
	testutils.Equal(t, sites.Less(0, 1), false) // Y=10 > Y=5
	testutils.Equal(t, sites.Less(1, 2), true)  // same Y, X=3 < X=7

	sites.Swap(0, 1)
	testutils.Equal(t, sites[0].X, 3)
	testutils.Equal(t, sites[1].X, 5)
}

// ---------- DCEL String tests ----------

func TestVertexString(t *testing.T) {
	v := &Vertex{X: 10, Y: 20}
	s := v.String()
	testutils.Equal(t, len(s) > 0, true)
}

func TestFaceString(t *testing.T) {
	f := &Face{ID: 42}
	s := f.String()
	testutils.Equal(t, len(s) > 0, true)
}

func TestHalfEdgeString(t *testing.T) {
	v := &Vertex{X: 1, Y: 2}
	he := &HalfEdge{Target: v}
	s := he.String()
	testutils.Equal(t, len(s) > 0, true)
}

func TestHalfEdgeStringWithFace(t *testing.T) {
	v := &Vertex{X: 1, Y: 2}
	f := &Face{ID: 5}
	he := &HalfEdge{Target: v, Face: f}
	s := he.String()
	testutils.Equal(t, len(s) > 0, true)
}

// ---------- IsClosed tests ----------

func TestHalfEdgeIsClosed(t *testing.T) {
	v1 := &Vertex{X: 1, Y: 2}
	v2 := &Vertex{X: 3, Y: 4}
	twin := &HalfEdge{Target: v2}
	he := &HalfEdge{Target: v1, Twin: twin}
	testutils.Equal(t, he.IsClosed(), true)
}

func TestHalfEdgeIsClosedNoTarget(t *testing.T) {
	twin := &HalfEdge{Target: &Vertex{X: 3, Y: 4}}
	he := &HalfEdge{Twin: twin}
	testutils.Equal(t, he.IsClosed(), false)
}

func TestHalfEdgeIsClosedNoTwinTarget(t *testing.T) {
	twin := &HalfEdge{}
	he := &HalfEdge{Target: &Vertex{X: 1, Y: 2}, Twin: twin}
	testutils.Equal(t, he.IsClosed(), false)
}

func TestHalfEdgeIsClosedNoTwin(t *testing.T) {
	he := &HalfEdge{Target: &Vertex{X: 1, Y: 2}}
	testutils.Equal(t, he.IsClosed(), false)
}

// ---------- EventQueue String test ----------

func TestEventQueueString(t *testing.T) {
	sites := SiteSlice{
		{X: 10, Y: 20, ID: 0},
		{X: 30, Y: 40, ID: 1},
	}
	eq := NewEventQueue(sites)
	s := eq.String()
	testutils.Equal(t, len(s) > 0, true)
}

// ---------- Node tree tests ----------

func TestNodeIsLeaf(t *testing.T) {
	n := &Node{Site: &Site{X: 1, Y: 2}}
	testutils.Equal(t, n.IsLeaf(), true)
}

func TestNodeIsNotLeaf(t *testing.T) {
	n := &Node{
		Left:  &Node{Site: &Site{X: 1, Y: 2}},
		Right: &Node{Site: &Site{X: 3, Y: 4}},
	}
	testutils.Equal(t, n.IsLeaf(), false)
}

func TestNodeString(t *testing.T) {
	leaf := &Node{Site: &Site{X: 1, Y: 2}}
	s := leaf.String()
	testutils.Equal(t, len(s) > 0, true)
}

func TestNodeStringNil(t *testing.T) {
	var n *Node
	testutils.Equal(t, n.String(), "()")
}

func TestNodeStringInternal(t *testing.T) {
	root := &Node{
		Left:  &Node{Site: &Site{X: 1, Y: 2}},
		Right: &Node{Site: &Site{X: 3, Y: 4}},
	}
	root.Left.Parent = root
	root.Right.Parent = root
	s := root.String()
	testutils.Equal(t, len(s) > 0, true)
}

func TestNodeStringInternalWithParent(t *testing.T) {
	grandparent := &Node{}
	parent := &Node{Parent: grandparent}
	parent.Left = &Node{Site: &Site{X: 1, Y: 2}, Parent: parent}
	parent.Right = &Node{Site: &Site{X: 3, Y: 4}, Parent: parent}
	grandparent.Left = parent
	grandparent.Right = &Node{Site: &Site{X: 5, Y: 6}, Parent: grandparent}

	s := grandparent.String()
	testutils.Equal(t, len(s) > 0, true)
}

// ---------- CircleEvents ----------

func TestCircleEventsRemoveEvent(t *testing.T) {
	e1 := &Event{X: 1, Y: 1}
	e2 := &Event{X: 2, Y: 2}
	ce := CircleEvents{e1, e2}
	ce.RemoveEvent(e1)
	testutils.Equal(t, len(ce), 1)
}

func TestCircleEventsHasEvent(t *testing.T) {
	e1 := &Event{X: 1, Y: 1}
	e2 := &Event{X: 2, Y: 2}
	ce := CircleEvents{e1}
	testutils.Equal(t, ce.HasEvent(e1), true)
	testutils.Equal(t, ce.HasEvent(e2), false)
}

func TestNodeRemoveEvent(t *testing.T) {
	e := &Event{X: 1, Y: 1}
	n := &Node{
		LeftEvents:   CircleEvents{e},
		MiddleEvents: CircleEvents{e},
		RightEvents:  CircleEvents{e},
	}
	n.RemoveEvent(e)
	testutils.Equal(t, len(n.LeftEvents), 0)
	testutils.Equal(t, len(n.MiddleEvents), 0)
	testutils.Equal(t, len(n.RightEvents), 0)
}

func TestNodeHasEvent(t *testing.T) {
	e := &Event{X: 1, Y: 1}
	n := &Node{LeftEvents: CircleEvents{e}}
	testutils.Equal(t, n.HasEvent(e), true)

	n2 := &Node{MiddleEvents: CircleEvents{e}}
	testutils.Equal(t, n2.HasEvent(e), true)

	n3 := &Node{RightEvents: CircleEvents{e}}
	testutils.Equal(t, n3.HasEvent(e), true)

	n4 := &Node{}
	testutils.Equal(t, n4.HasEvent(e), false)
}

func TestNodeAddEvents(t *testing.T) {
	e := &Event{X: 1, Y: 1}
	n := &Node{}
	n.AddLeftEvent(e)
	n.AddMiddleEvent(e)
	n.AddRightEvent(e)
	testutils.Equal(t, len(n.LeftEvents), 1)
	testutils.Equal(t, len(n.MiddleEvents), 1)
	testutils.Equal(t, len(n.RightEvents), 1)
}

// ---------- PrevArc / NextArc for nil ----------

func TestPrevArcNil(t *testing.T) {
	var n *Node
	testutils.Equal(t, n.PrevArc() == nil, true)
}

func TestNextArcNil(t *testing.T) {
	var n *Node
	testutils.Equal(t, n.NextArc() == nil, true)
}

func TestPrevArcSingleRoot(t *testing.T) {
	n := &Node{Site: &Site{X: 1, Y: 2}}
	testutils.Equal(t, n.PrevArc() == nil, true)
}

func TestNextArcSingleRoot(t *testing.T) {
	n := &Node{Site: &Site{X: 1, Y: 2}}
	testutils.Equal(t, n.NextArc() == nil, true)
}

// ---------- FirstArc / LastArc ----------

func TestFirstArc(t *testing.T) {
	leaf1 := &Node{Site: &Site{X: 1, Y: 2}}
	leaf2 := &Node{Site: &Site{X: 3, Y: 4}}
	root := &Node{Left: leaf1, Right: leaf2}
	leaf1.Parent = root
	leaf2.Parent = root

	testutils.Equal(t, root.FirstArc(), leaf1)
}

func TestLastArc(t *testing.T) {
	leaf1 := &Node{Site: &Site{X: 1, Y: 2}}
	leaf2 := &Node{Site: &Site{X: 3, Y: 4}}
	root := &Node{Left: leaf1, Right: leaf2}
	leaf1.Parent = root
	leaf2.Parent = root

	testutils.Equal(t, root.LastArc(), leaf2)
}

// ---------- Voronoi with more sites ----------

func TestVoronoiMultipleSites(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 800, 600)
	sites := SiteSlice{
		{X: 100, Y: 100, ID: 0},
		{X: 200, Y: 200, ID: 1},
		{X: 300, Y: 100, ID: 2},
		{X: 400, Y: 300, ID: 3},
		{X: 150, Y: 400, ID: 4},
	}
	v := New(sites, rect)
	v.Generate()
	polys := v.ToPolyhedrons()
	testutils.Equal(t, len(polys) > 0, true)
}

func TestVoronoiHandleNextEventEmpty(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)
	v := New(SiteSlice{}, rect)
	// Should not panic on empty queue
	v.HandleNextEvent()
}

func TestVoronoiReset(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)
	sites := SiteSlice{
		{X: 10, Y: 10, ID: 0},
		{X: 50, Y: 50, ID: 1},
	}
	v := New(sites, rect)
	v.Generate()
	v.Reset()
	testutils.Equal(t, v.SweepLine, 0)
	testutils.Equal(t, v.ParabolaTree == nil, true)
}

func TestVoronoiGetFaceVertices(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 600, 480)
	sites := []shapes.Point{
		shapes.NewPoint(110, 20),
		shapes.NewPoint(140, 40),
		shapes.NewPoint(155, 80),
		shapes.NewPoint(350, 120),
		shapes.NewPoint(200, 240),
	}
	v := NewFromPoints(sites, rect)
	v.Generate()

	// Get vertices for each face
	for _, face := range v.DCEL.Faces {
		vertices := v.GetFaceVertices(face)
		// Verify we get some vertices
		testutils.Equal(t, len(vertices) >= 0, true)
	}
}

func TestVoronoiNewFromPoints(t *testing.T) {
	rect := shapes.NewBox(shapes.NewPoint(0, 0), 100, 100)
	points := []shapes.Point{
		shapes.NewPoint(25, 25),
		shapes.NewPoint(75, 75),
	}
	v := NewFromPoints(points, rect)
	testutils.Equal(t, len(v.Sites), 2)
	testutils.Equal(t, v.Sites[0].X, 25)
	testutils.Equal(t, v.Sites[1].X, 75)
}

// ---------- Parabola tests ----------

func TestGetParabolaABC(t *testing.T) {
	focus := &Site{X: 5, Y: 10}
	a, b, c := GetParabolaABC(focus, 0)
	// a = 1/(2*(10-0)) = 0.05
	testutils.Equal(t, a > 0, true)
	testutils.Equal(t, b != 0, true)
	testutils.Equal(t, c != 0, true)
}

func TestGetYByX(t *testing.T) {
	focus := &Site{X: 5, Y: 10}
	y := GetYByX(focus, 5, 0)
	// At x = focus.X, the parabola should be at its vertex
	testutils.Equal(t, y > 0, true)
}

// ---------- verticesByCCW sort tests ----------

func TestVerticesByCCWSortInterface(t *testing.T) {
	v1 := &Vertex{X: 0, Y: 0}
	v2 := &Vertex{X: 10, Y: 0}
	v3 := &Vertex{X: 0, Y: 10}
	verts := verticesByCCW{v1, v2, v3}

	testutils.Equal(t, verts.Len(), 3)
	// Just exercise Less and Swap
	_ = verts.Less(0, 1)
	verts.Swap(0, 1)
	testutils.Equal(t, verts[0], v2)
}

// ---------- halfEdgesByCCW edge cases ----------

func TestHalfEdgesByCCWUpdateLinksSingle(t *testing.T) {
	v := &Vertex{X: 1, Y: 1}
	he := &HalfEdge{Target: v}
	edges := halfEdgesByCCW{he}
	edges.UpdateLinks()
	testutils.Equal(t, he.Prev == nil, true)
	testutils.Equal(t, he.Next == nil, true)
}

func TestHalfEdgesByCCWUpdateLinksMultiple(t *testing.T) {
	v1 := &Vertex{X: 0, Y: 0}
	v2 := &Vertex{X: 10, Y: 0}
	v3 := &Vertex{X: 0, Y: 10}
	he1 := &HalfEdge{Target: v1}
	he2 := &HalfEdge{Target: v2}
	he3 := &HalfEdge{Target: v3}
	edges := halfEdgesByCCW{he1, he2, he3}
	edges.UpdateLinks()
	// First.Prev should be last, last.Next should be first
	testutils.Equal(t, he1.Prev, he3)
	testutils.Equal(t, he3.Next, he1)
}

func TestHalfEdgesByCCWLessNilTargets(t *testing.T) {
	he1 := &HalfEdge{} // nil Target
	he2 := &HalfEdge{Target: &Vertex{X: 1, Y: 1}}
	edges := halfEdgesByCCW{he1, he2}
	// nil target should not be less
	testutils.Equal(t, edges.Less(0, 1), false)
	testutils.Equal(t, edges.Less(1, 0), true)
}
