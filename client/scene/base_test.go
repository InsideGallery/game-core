package scene

import "testing"

func TestViewportNeedsResize(t *testing.T) {
	tests := []struct {
		name     string
		curW     int
		curH     int
		newW     int
		newH     int
		hasWorld bool
		want     bool
	}{
		{"no world always needs resize", 800, 600, 800, 600, false, true},
		{"same dims with world — skip", 800, 600, 800, 600, true, false},
		{"width increased", 800, 600, 1024, 600, true, true},
		{"width decreased", 1024, 600, 800, 600, true, true},
		{"height increased", 800, 600, 800, 768, true, true},
		{"height decreased", 800, 768, 800, 600, true, true},
		{"both dims changed", 800, 600, 1024, 768, true, true},
		{"zero cur → first allocation without world", 0, 0, 800, 600, false, true},
		{"zero cur with world (impossible in practice) — dims differ", 0, 0, 800, 600, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := viewportNeedsResize(tc.curW, tc.curH, tc.newW, tc.newH, tc.hasWorld)
			if got != tc.want {
				t.Errorf("viewportNeedsResize(%d,%d,%d,%d,hasWorld=%v)=%v, want %v",
					tc.curW, tc.curH, tc.newW, tc.newH, tc.hasWorld, got, tc.want)
			}
		})
	}
}
