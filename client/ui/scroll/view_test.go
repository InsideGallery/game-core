package scroll

import "testing"

func TestApply(t *testing.T) {
	v := View{ViewportH: 100, ContentH: 300}
	v.Apply(-50, 0, 0)
	if v.Offset != 50 {
		t.Fatalf("offset = %v, want 50", v.Offset)
	}
	v.Apply(-500, 0, 0)
	if v.Offset != 200 {
		t.Fatalf("offset = %v, want 200", v.Offset)
	}
	v.Apply(0, 1, 30)
	if v.Offset != 170 {
		t.Fatalf("wheel offset = %v, want 170", v.Offset)
	}
	v.Apply(500, 0, 0)
	if v.Offset != 0 {
		t.Fatalf("offset = %v, want 0", v.Offset)
	}
}

func TestClampShortContent(t *testing.T) {
	v := View{Offset: 40, ViewportH: 100, ContentH: 50}
	v.Clamp()
	if v.Offset != 0 {
		t.Fatalf("offset = %v, want 0", v.Offset)
	}
}

func TestAtEnd(t *testing.T) {
	v := View{Offset: 190, ViewportH: 100, ContentH: 300}
	if !v.AtEnd(20) {
		t.Fatal("AtEnd(20) = false")
	}
	if v.AtEnd(5) {
		t.Fatal("AtEnd(5) = true")
	}
}

func TestVisible(t *testing.T) {
	v := View{Offset: 100, ViewportH: 100, ContentH: 500}
	cases := []struct {
		top, h float64
		want   bool
	}{{0, 100, false}, {0, 101, true}, {150, 20, true}, {200, 10, false}, {199, 10, true}}
	for _, c := range cases {
		if got := v.Visible(c.top, c.h); got != c.want {
			t.Errorf("Visible(%v,%v) = %v", c.top, c.h, got)
		}
	}
}
