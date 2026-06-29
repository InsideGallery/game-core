package core

import (
	"context"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type stubSystem struct {
	name    string
	updated bool
	drawn   bool
}

func (s *stubSystem) Update(_ context.Context) error {
	s.updated = true

	return nil
}

func (s *stubSystem) Draw(_ context.Context, _ *ebiten.Image) {
	s.drawn = true
}

func TestNewSystems(t *testing.T) {
	s := NewSystems()
	if s == nil {
		t.Fatal("expected non-nil systems")
	}

	if len(s.Get()) != 0 {
		t.Fatalf("expected 0 systems, got %d", len(s.Get()))
	}
}

func TestAddAndGet(t *testing.T) {
	s := NewSystems()
	a := &stubSystem{name: "a"}
	b := &stubSystem{name: "b"}
	c := &stubSystem{name: "c"}

	s.Add("a", a)
	s.Add("b", b)
	s.Add("c", c)

	got := s.Get()
	if len(got) != 3 {
		t.Fatalf("expected 3 systems, got %d", len(got))
	}

	if got[0].(*stubSystem).name != "a" {
		t.Fatal("expected first system 'a'")
	}

	if got[1].(*stubSystem).name != "b" {
		t.Fatal("expected second system 'b'")
	}

	if got[2].(*stubSystem).name != "c" {
		t.Fatal("expected third system 'c'")
	}
}

func TestAddDuplicateKeepsOrder(t *testing.T) {
	s := NewSystems()
	a1 := &stubSystem{name: "a1"}
	a2 := &stubSystem{name: "a2"}
	b := &stubSystem{name: "b"}

	s.Add("a", a1)
	s.Add("b", b)
	s.Add("a", a2) // re-register "a" — order must not change

	got := s.Get()
	if len(got) != 2 {
		t.Fatalf("expected 2 systems, got %d", len(got))
	}

	if got[0].(*stubSystem).name != "a2" {
		t.Fatal("expected updated system at position 0")
	}

	if got[1].(*stubSystem).name != "b" {
		t.Fatal("expected 'b' at position 1")
	}
}

func TestRemove(t *testing.T) {
	s := NewSystems()
	s.Add("a", &stubSystem{name: "a"})
	s.Add("b", &stubSystem{name: "b"})
	s.Add("c", &stubSystem{name: "c"})

	s.Remove("b")

	got := s.Get()
	if len(got) != 2 {
		t.Fatalf("expected 2 systems after remove, got %d", len(got))
	}

	if got[0].(*stubSystem).name != "a" {
		t.Fatal("expected first system 'a'")
	}

	if got[1].(*stubSystem).name != "c" {
		t.Fatal("expected second system 'c'")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	s := NewSystems()
	s.Add("a", &stubSystem{name: "a"})

	s.Remove("nonexistent")

	if len(s.Get()) != 1 {
		t.Fatal("should not remove anything")
	}
}

func TestClean(t *testing.T) {
	s := NewSystems()
	s.Add("a", &stubSystem{name: "a"})
	s.Add("b", &stubSystem{name: "b"})

	s.Clean()

	if len(s.Get()) != 0 {
		t.Fatalf("expected 0 after clean, got %d", len(s.Get()))
	}
}
