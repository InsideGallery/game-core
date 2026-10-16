package input

import "testing"

func TestTrackerTap(t *testing.T) {
	tr := newTracker(10)

	for _, s := range []struct{ x, y float64 }{{5, 5}, {8, 8}} {
		if ev := tr.step(true, s.x, s.y); ev.tap || ev.drag {
			t.Fatalf("unexpected event while held: %+v", ev)
		}
	}

	ev := tr.step(false, 0, 0)
	if !ev.tap || ev.tapX != 8 || ev.tapY != 8 || ev.drag {
		t.Fatalf("want tap at 8,8 and no drag, got %+v", ev)
	}

	if ev := tr.step(false, 0, 0); ev.tap || ev.drag {
		t.Fatalf("idle frame reported %+v", ev)
	}
}

func TestTrackerDrag(t *testing.T) {
	tr := newTracker(10)

	if ev := tr.step(true, 0, 0); ev.tap || ev.drag {
		t.Fatalf("press: %+v", ev)
	}

	ev := tr.step(true, 0, 30)
	if !ev.drag || ev.dx != 0 || ev.dy != 30 {
		t.Fatalf("want drag 0,30 got %+v", ev)
	}

	ev = tr.step(true, 0, 50)
	if !ev.drag || ev.dx != 0 || ev.dy != 20 {
		t.Fatalf("want drag 0,20 got %+v", ev)
	}

	if ev = tr.step(false, 0, 50); ev.tap || ev.drag {
		t.Fatalf("release must not tap: %+v", ev)
	}
}

func TestTrackerPressing(t *testing.T) {
	tr := newTracker(10)

	if tr.pressing() {
		t.Fatal("idle tracker pressing")
	}

	tr.step(true, 0, 0)

	if !tr.pressing() {
		t.Fatal("held press not pressing")
	}

	tr.step(true, 0, 30)

	if tr.pressing() {
		t.Fatal("drag still pressing")
	}

	tr.step(false, 0, 0)

	if tr.pressing() {
		t.Fatal("released still pressing")
	}
}
