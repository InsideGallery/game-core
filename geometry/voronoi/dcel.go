package voronoi

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// DCEL stores the state of the data structure and provides methods for linking of three sets of
// objects: vertecies, edges and faces.
type DCEL struct {
	Vertices  []*Vertex
	Faces     []*Face
	HalfEdges []*HalfEdge
}

// Vertex represents a node in the DCEL structure. Each vertex has 2D coordinates and a pointer
// to an arbitrary half edge that has this vertex as its target (origin). Annotations (user data)
// can be stored in the Data field.
type Vertex struct {
	X, Y     int
	HalfEdge *HalfEdge
	Data     interface{}
}

// Face represents a subdivision of the plane. Each face has a pointer to one of the half edges
// at its boundary. Faces can have user specified IDs and annotations.
type Face struct {
	HalfEdge *HalfEdge
	ID       int64
	Data     interface{}
}

// HalfEdge represents one of the half-edges in an edge pair. Each half-edge has a pointer to its
// target vertex (origin), the face to which it belongs, its twin edge (a reversed half-edge, pointing
// to a neighbour face) and pointers to the next and previous half-edges at the boundary of its face.
// Half-edges can also store user data.
type HalfEdge struct {
	Target *Vertex
	Face   *Face
	Twin   *HalfEdge
	Next   *HalfEdge
	Prev   *HalfEdge
	Data   interface{}
}

// String implement string function
func (v *Vertex) String() string {
	return fmt.Sprintf("{Vertex %p; X,Y: %d,%d; Edge: %p}", v, v.X, v.Y, v.HalfEdge)
}

// String implement string function
func (f *Face) String() string {
	return fmt.Sprintf("{Face #%d %p}", f.ID, f)
}

// String implement string function
func (he *HalfEdge) String() string {
	faceID := "nil"
	if he.Face != nil {
		faceID = "#" + strconv.FormatInt(he.Face.ID, 10) //nolint:mnd
	}

	return fmt.Sprintf("{Edge %p; Target: %d,%d; Twin: %p; Face: %s}", he, he.Target.X, he.Target.Y, he.Twin, faceID)
}

// IsClosed returns true if both half-edges in the pair have a target vertex.
func (he *HalfEdge) IsClosed() bool {
	return he.Target != nil && he.Twin != nil && he.Twin.Target != nil
}

// NewDCEL creates a new DCEL data structure.
func NewDCEL() *DCEL {
	return &DCEL{}
}

// NewFace creates a new face and stores it in the DCEL structure.
func (d *DCEL) NewFace() *Face {
	face := &Face{}
	d.Faces = append(d.Faces, face)

	return face
}

// NewVertex creates a new vertex with the given coordinates and stores it in the structure.
func (d *DCEL) NewVertex(x, y int) *Vertex {
	vertex := &Vertex{
		X: x,
		Y: y,
	}
	d.Vertices = append(d.Vertices, vertex)

	return vertex
}

// NewHalfEdge creates a new half-edge starting at the given vertex and stores it in the structure.
func (d *DCEL) NewHalfEdge(face *Face, vertex *Vertex, twin *HalfEdge) *HalfEdge {
	halfEdge := &HalfEdge{
		Face:   face,
		Target: vertex,
		Twin:   twin,
	}

	// Link new half-edge to the one pointed by the face counter-clockwise
	if face.HalfEdge != nil {
		halfEdge.Prev = face.HalfEdge.Prev
		if halfEdge.Prev != nil {
			halfEdge.Prev.Next = halfEdge
		}
		face.HalfEdge.Prev = halfEdge
		halfEdge.Next = face.HalfEdge
	}

	// Set the new half-edge as the one pointed by the face
	face.HalfEdge = halfEdge

	// If the vertex is not yet linked with an edge, link to this one
	if vertex != nil && vertex.HalfEdge == nil {
		vertex.HalfEdge = halfEdge
	}

	d.HalfEdges = append(d.HalfEdges, halfEdge)

	return halfEdge
}

// NewEdge creates a pair of half-edges, one of them starting at the given vertex.
func (d *DCEL) NewEdge(face1, face2 *Face, vertex *Vertex) (*HalfEdge, *HalfEdge) {
	halfEdge := d.NewHalfEdge(face1, vertex, nil)
	twin := d.NewHalfEdge(face2, nil, halfEdge)
	halfEdge.Twin = twin

	return halfEdge, twin
}

// CloseTwins adds a vertex to the specified edges.
func (v *Voronoi) CloseTwins(list []*HalfEdge, vertex *Vertex) {
	for i := 0; i < len(list); i++ {
		he := list[i]
		if he.Twin != nil && he.Twin.Target == nil {
			he.Twin.Target = vertex
		} else if he.Target == nil {
			he.Target = vertex
		}
	}
}

// halfEdgesByCCW implements a slice of half-edges that sort in counter-clockwise order.
type halfEdgesByCCW []*HalfEdge

func (s halfEdgesByCCW) Len() int {
	return len(s)
}

func (s halfEdgesByCCW) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s halfEdgesByCCW) Less(i, j int) bool {
	if s[i].Target == nil {
		return false
	} else if s[j].Target == nil {
		return true
	}

	// Find center of polygon
	var sumX int64
	var sumY int64
	var cnt int

	for _, v := range s {
		if v.Target != nil {
			sumX += int64(v.Target.X)
			sumY += int64(v.Target.Y)
			cnt++
		}
	}

	centerX := float64(sumX) / float64(cnt)
	centerY := float64(sumY) / float64(cnt)

	// Sort counter-clockwise
	a1 := math.Atan2(float64(s[i].Target.Y)-centerY, float64(s[i].Target.X)-centerX)
	a2 := math.Atan2(float64(s[j].Target.Y)-centerY, float64(s[j].Target.X)-centerX)

	return a1 >= a2
}

func (s halfEdgesByCCW) UpdateLinks() {
	for i := 0; i < len(s); i++ {
		if i > 0 {
			s[i].Prev, s[i-1].Next = s[i-1], s[i]
		}

		if i < len(s)-1 {
			s[i].Next, s[i+1].Prev = s[i+1], s[i]
		}
	}

	if len(s) == 1 {
		s[0].Prev, s[0].Next = nil, nil
	} else if len(s) > 1 {
		s[0].Prev, s[len(s)-1].Next = s[len(s)-1], s[0]
	}
}

// ReorderFaceEdges reorders face half-edges in a clockwise way, while also removing duplicates.
func (v *Voronoi) ReorderFaceEdges(face *Face) {
	var edges []*HalfEdge
	edge := face.HalfEdge

	for edge != nil {
		edges = append(edges, edge)
		edge = edge.Next

		if edge == face.HalfEdge {
			break
		}
	}

	sort.Sort(halfEdgesByCCW(edges))
	halfEdgesByCCW(edges).UpdateLinks()
}

// GetFaceHalfEdges returns the half-edges that form the boundary of a face (cell).
func (v *Voronoi) GetFaceHalfEdges(face *Face) []*HalfEdge {
	v.ReorderFaceEdges(face)

	var edges []*HalfEdge
	edge := face.HalfEdge

	for edge != nil {
		edges = append(edges, edge)

		edge = edge.Next
		if edge == face.HalfEdge {
			break
		}
	}

	return edges
}

// verticesByCCW implements a slice of vertices that sort in counter-clockwise order.
type verticesByCCW []*Vertex

func (s verticesByCCW) Len() int {
	return len(s)
}

func (s verticesByCCW) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s verticesByCCW) Less(i, j int) bool {
	// Find center of polygon
	var sumX float64
	var sumY float64

	for _, v := range s {
		sumX += float64(v.X)
		sumY += float64(v.Y)
	}

	centerX := sumX / float64(len(s))
	centerY := sumY / float64(len(s))

	// Sort counter-clockwise
	a1 := math.Atan2(float64(s[i].Y)-centerY, float64(s[i].X)-centerX)
	a2 := math.Atan2(float64(s[j].Y)-centerY, float64(s[j].X)-centerX)

	return a1 >= a2
}

// GetFaceVertices returns the vertices that form the boundary of a face (cell),
// sorted in counter-clockwise order.
func (v *Voronoi) GetFaceVertices(face *Face) []*Vertex {
	var vertices []*Vertex
	exists := make(map[string]bool)
	edge := face.HalfEdge

	for edge != nil {
		if edge.Target != nil {
			id := fmt.Sprintf("%v", edge.Target)
			if !exists[id] {
				exists[id] = true
				vertices = append(vertices, edge.Target)
			}
		}

		if edge.Twin != nil && edge.Twin.Target != nil {
			id := fmt.Sprintf("%v", edge.Twin.Target)
			if !exists[id] {
				exists[id] = true
				vertices = append(vertices, edge.Twin.Target)
			}
		}

		edge = edge.Next
		if edge == face.HalfEdge {
			break
		}
	}

	sort.Sort(verticesByCCW(vertices))

	return vertices
}
