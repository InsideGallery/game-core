// Package scroll provides a clamped vertical scroll offset driven by drag and wheel.
package scroll

// View is a vertical scroll state. Offset is the content distance above the viewport top.
type View struct {
	Offset    float64
	ViewportH float64
	ContentH  float64
}

func (v *View) maxOffset() float64 {
	return max(0, v.ContentH-v.ViewportH)
}

// Apply scrolls by a drag and a wheel step, then clamps. Dragging up (negative dragDY) scrolls down.
func (v *View) Apply(dragDY, wheel, wheelStep float64) {
	v.Offset -= dragDY
	v.Offset -= wheel * wheelStep
	v.Clamp()
}

// Clamp keeps Offset within [0, max(0, ContentH-ViewportH)].
func (v *View) Clamp() {
	v.Offset = min(max(v.Offset, 0), v.maxOffset())
}

// AtEnd reports whether the viewport is within margin of the content end.
func (v *View) AtEnd(margin float64) bool {
	return v.Offset >= v.ContentH-v.ViewportH-margin
}

// Visible reports whether the item at top (content coordinates) with height crosses the viewport.
func (v *View) Visible(top, height float64) bool {
	return top+height > v.Offset && top < v.Offset+v.ViewportH
}
