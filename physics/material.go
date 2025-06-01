package physics

// Material describe material
type Material struct {
	Mass float64
}

// NewMaterial return new material
func NewMaterial(mass float64) Material {
	return Material{
		Mass: mass,
	}
}
