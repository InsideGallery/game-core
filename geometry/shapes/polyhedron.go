package shapes

// Polyhedron describe polygon
type Polyhedron struct {
	vectors []Point
}

// NewPolyhedron return new polygon
func NewPolyhedron(vectors ...Point) Polyhedron {
	return Polyhedron{
		vectors: vectors,
	}
}

// Vectors return all vectors
func (p Polyhedron) Vectors() []Point {
	return p.vectors
}

// Count return count of objects
func (p Polyhedron) Count() int {
	return len(p.vectors)
}

// Vector return single point
func (p Polyhedron) Vector(i int) Point {
	return p.vectors[i]
}

// ToLines return lines of polyhedron
func (p Polyhedron) ToLines() (line []Line) {
	next := 0
	for current := 0; current < p.Count(); current++ {
		next = current + 1
		if next == p.Count() {
			next = 0
		}

		vc := p.Vector(current)
		vn := p.Vector(next)

		line = append(line, NewLine(vc, vn))
	}

	return line
}

// Coordinates return 3d coordinates
func (p Polyhedron) Coordinates() [3]float64 {
	return p.Bounds().Coordinates()
}

// Get return Spatial
func (p Polyhedron) Get() Spatial {
	return p
}

// Move return new object for given diff
func (p Polyhedron) Move(diff Point) Spatial {
	vectors := make([]Point, p.Count())

	for i, point := range p.Vectors() {
		p := make([]float64, 3) // nolint:mnd
		for k, v := range point.Coordinates() {
			p[k] = v + diff.Coordinate(k)
		}

		vectors[i] = NewPoint(p...)
	}

	return NewPolyhedron(vectors...)
}

// Bounds return rectangle of object
func (p Polyhedron) Bounds() Box {
	var set bool
	minPoint := make([]float64, 3) // nolint:mnd
	maxPoint := make([]float64, 3) // nolint:mnd
	sizes := make([]float64, 3)    // nolint:mnd

	for _, point := range p.Vectors() {
		if set {
			for k, v := range point.Coordinates() {
				if v < minPoint[k] {
					minPoint[k] = v
				}

				if v > maxPoint[k] {
					maxPoint[k] = v
				}
			}
		} else {
			for k, v := range point.Coordinates() {
				minPoint[k] = v
				maxPoint[k] = v
			}
			set = true
		}
	}

	for i := range maxPoint {
		sizes[i] = maxPoint[i] - minPoint[i]
	}

	return NewBox(NewPoint(minPoint...), sizes...)
}

// Center return center point
func (p Polyhedron) Center() Point {
	ps := make([]float64, 3) // nolint:mnd
	l := float64(p.Count())

	for _, v := range p.Vectors() {
		for i := 0; i < 3; i++ {
			ps[i] += v.Coordinate(i)
		}
	}

	for i := 0; i < 3; i++ {
		ps[i] /= l
	}

	return NewPoint(ps...)
}

// Point1 return first point
func (p Polyhedron) Point1() Point {
	return p.Vector(0) // nolint:mnd
}

// Support return support point
func (p Polyhedron) Support(d Point) Point {
	furthestPoint := p.Vector(0) // nolint:mnd
	maxDot := furthestPoint.Dot(d)

	for i := 1; i < p.Count(); i++ {
		v := p.Vector(i)
		dot := v.Dot(d)

		if dot > maxDot {
			maxDot = dot
			furthestPoint = v
		}
	}

	return furthestPoint
}
