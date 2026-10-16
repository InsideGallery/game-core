// Package input reads one pointer, touch or mouse, and reports taps and drags.
package input

// Pointer reports taps, drags, hover, press and wheel movement for the first touch or the left mouse button.
type Pointer struct {
	tr     *tracker
	src    Source
	ev     event
	sample Sample
}

// NewPointer creates a Pointer on the ebiten input. slop is in logical pixels; callers pass 10 times the device scale.
func NewPointer(slop float64) *Pointer {
	return NewPointerFrom(slop, EbitenSource())
}

// NewPointerFrom creates a Pointer that reads from src.
func NewPointerFrom(slop float64, src Source) *Pointer {
	return &Pointer{tr: newTracker(slop), src: src}
}

// Update reads the input once per frame.
func (p *Pointer) Update() {
	p.sample = p.src.Read()
	p.ev = p.tr.step(p.sample.Pressed, p.sample.X, p.sample.Y)
}

// Hover returns the mouse position while no button is down. It is never ok for a touch.
func (p *Pointer) Hover() (x, y float64, ok bool) {
	return p.sample.X, p.sample.Y, !p.sample.Touch && !p.sample.Pressed
}

// Press returns the current position while a press is down and has not become a drag.
func (p *Pointer) Press() (x, y float64, ok bool) {
	return p.sample.X, p.sample.Y, p.tr.pressing()
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
	return p.sample.Wheel
}
