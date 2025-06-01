package shapes

type Rotatable struct {
	Spatial
	Angle float64
}

func NewRotatable(spatial Spatial, angle float64) Rotatable {
	return Rotatable{
		Spatial: spatial,
		Angle:   angle,
	}
}

func (r *Rotatable) GetAngle() float64 {
	return r.Angle
}
