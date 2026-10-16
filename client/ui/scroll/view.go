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

// Bar returns the scroll bar of the current offset in viewport coordinates (0 is the viewport top).
// ok is false when nothing scrolls. The length is at least minLen and at most ViewportH.
func (v *View) Bar(minLen float64) (top, length float64, ok bool) {
	if v.ViewportH <= 0 || v.ContentH <= v.ViewportH {
		return 0, 0, false
	}

	length = min(max(minLen, v.ViewportH*v.ViewportH/v.ContentH), v.ViewportH)
	top = (v.ViewportH - length) * v.Offset / (v.ContentH - v.ViewportH)

	return top, length, true
}
