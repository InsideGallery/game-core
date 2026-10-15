// Package input reads one pointer, touch or mouse, and reports taps and drags.
package input

import "github.com/hajimehoshi/ebiten/v2"

// Pointer reports taps, drags and wheel movement for the first touch or the left mouse button.
type Pointer struct {
	tr    *tracker
	ev    event
	wheel float64
	ids   []ebiten.TouchID
}

// NewPointer creates a Pointer. slop is in logical pixels; callers pass 10 times the device scale.
func NewPointer(slop float64) *Pointer {
	return &Pointer{tr: newTracker(slop)}
}

// Update reads the input once per frame.
func (p *Pointer) Update() {
	p.ids = ebiten.AppendTouchIDs(p.ids[:0])

	if len(p.ids) > 0 {
		x, y := ebiten.TouchPosition(p.ids[0])
		p.ev = p.tr.step(true, float64(x), float64(y))
	} else {
		x, y := ebiten.CursorPosition()
		p.ev = p.tr.step(ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft), float64(x), float64(y))
	}

	_, p.wheel = ebiten.Wheel()
}

// Tap returns the position of a tap that ended this frame.
func (p *Pointer) Tap() (x, y float64, ok bool) {
	return p.ev.tapX, p.ev.tapY, p.ev.tap
}

// Drag returns the movement since the last frame while a press is dragging.
func (p *Pointer) Drag() (dx, dy float64, active bool) {
	return p.ev.dx, p.ev.dy, p.ev.drag
}

// Wheel returns the vertical mouse wheel delta of this frame.
func (p *Pointer) Wheel() float64 {
	return p.wheel
}
