package input

import "math"

// event is the result of one tracker step.
type event struct {
	tap    bool
	tapX   float64
	tapY   float64
	drag   bool
	dx, dy float64
}

// tracker is the pure tap and drag state machine.
type tracker struct {
	slop     float64
	down     bool
	dragging bool
	startX   float64
	startY   float64
	lastX    float64
	lastY    float64
}

func newTracker(slop float64) *tracker {
	return &tracker{slop: slop}
}

// step advances the machine by one frame with the current pressed state and position.
func (t *tracker) step(pressed bool, x, y float64) event {
	var ev event

	switch {
	case pressed && !t.down:
		t.down = true
		t.dragging = false
		t.startX, t.startY = x, y
		t.lastX, t.lastY = x, y
	case pressed:
		if !t.dragging && math.Hypot(x-t.startX, y-t.startY) > t.slop {
			t.dragging = true
		}

		if t.dragging {
			ev.drag = true
			ev.dx, ev.dy = x-t.lastX, y-t.lastY
		}

		t.lastX, t.lastY = x, y
	case t.down:
		if !t.dragging {
			ev.tap = true
			ev.tapX, ev.tapY = t.lastX, t.lastY
		}

		t.down = false
		t.dragging = false
	}

	return ev
}
