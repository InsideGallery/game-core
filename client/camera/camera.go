package camera

import (
	"math"

	"github.com/InsideGallery/game-core/geometry/shapes"
	"github.com/hajimehoshi/ebiten/v2"
)

// Camera provides world-space transformations: pan and zoom.
// Scenes composite the World buffer via WorldMatrix; input pick uses ScreenToWorld.
type Camera struct {
	Position   shapes.Point
	ZoomFactor float64
	vpW, vpH   int
}

// NewCamera creates a camera centred at pos with no zoom.
func NewCamera(pos shapes.Point) *Camera {
	return &Camera{Position: pos}
}

// SetViewPort records the logical viewport size (usually from Layout).
func (c *Camera) SetViewPort(w, h int) {
	c.vpW = w
	c.vpH = h
}

// WorldMatrix returns the affine transform from world space to screen space:
// translate by -Position → scale by zoom → translate to viewport centre.
func (c *Camera) WorldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position.Coordinate(0), -c.Position.Coordinate(1))

	zoom := math.Pow(1.01, c.ZoomFactor)
	m.Scale(zoom, zoom)

	m.Translate(float64(c.vpW)*0.5, float64(c.vpH)*0.5)
	return m
}

// ScreenToWorld converts screen coordinates to world coordinates.
func (c *Camera) ScreenToWorld(sx, sy float64) (float64, float64) {
	inv := c.WorldMatrix()
	if inv.IsInvertible() {
		inv.Invert()
		return inv.Apply(sx, sy)
	}
	return math.NaN(), math.NaN()
}

// WorldToScreen converts world coordinates to screen coordinates.
func (c *Camera) WorldToScreen(wx, wy float64) (float64, float64) {
	m := c.WorldMatrix()
	return m.Apply(wx, wy)
}
