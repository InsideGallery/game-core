package shapes

import "math"

// MultiObject contains different object to describe another object
type MultiObject struct {
	objects []Spatial
}

// NewMultiObject return new multi-object
func NewMultiObject(objects ...Spatial) MultiObject {
	return MultiObject{
		objects: objects,
	}
}

// Point1 return point1
func (m MultiObject) Point1() Point {
	return m.Bounds().Point1()
}

// Objects return all objects
func (m MultiObject) Objects() []Spatial {
	return m.objects
}

// Count return count of objects
func (m MultiObject) Count() int {
	return len(m.objects)
}

// Coordinates return 3d coordinates
func (m MultiObject) Coordinates() [3]float64 {
	return m.Bounds().Coordinates()
}

// Object return single object
func (m MultiObject) Object(i int) Spatial {
	return m.objects[i]
}

// Get return Spatial
func (m MultiObject) Get() Spatial {
	return m
}

// Move return new object for given diff
func (m MultiObject) Move(diff Point) Spatial {
	entities := make([]Spatial, m.Count())
	for i, o := range m.Objects() {
		entities[i] = o.Move(diff)
	}

	return NewMultiObject(entities...)
}

// Bounds return rectangle of object
func (m MultiObject) Bounds() Box {
	var r Box
	for _, o := range m.objects {
		r = r.BoundingBox(o.Bounds())
	}

	return r
}

// Center return center point
func (m MultiObject) Center() Point {
	p := make([]Point, m.Count())
	for k, v := range m.Objects() {
		p[k] = v.Center()
	}

	return NewPolyhedron(p...).Center()
}

// Support return support point
func (m MultiObject) Support(d Point) Point {
	var maxPoint Point
	maxDot := -math.MaxFloat64

	for _, o := range m.Objects() {
		if v, ok := o.(Collide); ok {
			p := v.Support(d)
			dot := p.Dot(d)

			if dot > maxDot {
				maxDot = dot
				maxPoint = p
			}
		}
	}

	return maxPoint
}
