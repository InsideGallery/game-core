package scene

import (
	"context"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// stubScene is a minimal Scene implementation for testing the Manager.
type stubScene struct {
	name     string
	loaded   bool
	unloaded bool
}

func (s *stubScene) Name() string                   { return s.name }
func (s *stubScene) Load(_ context.Context)         { s.loaded = true }
func (s *stubScene) Unload(_ context.Context)       { s.unloaded = true }
func (s *stubScene) Update(_ context.Context) error { return nil }
func (s *stubScene) Draw(_ *ebiten.Image)           {}

func TestNewManager(t *testing.T) {
	m := NewManager(context.Background())
	if m == nil {
		t.Fatal("expected non-nil manager")
	}

	if m.Scene() != nil {
		t.Fatal("expected nil active scene on new manager")
	}
}

func TestAdd(t *testing.T) {
	m := NewManager(context.Background())
	sc := &stubScene{name: "game"}
	m.Add("game", sc)

	if err := m.SwitchSceneTo("game"); err != nil {
		t.Fatalf("switch: %v", err)
	}

	if m.Scene() != sc {
		t.Fatal("expected active scene to be 'game'")
	}
}

func TestSwitchSceneToLoads(t *testing.T) {
	m := NewManager(context.Background())
	sc := &stubScene{name: "menu"}
	m.Add("menu", sc)

	if err := m.SwitchSceneTo("menu"); err != nil {
		t.Fatalf("switch: %v", err)
	}

	if !sc.loaded {
		t.Fatal("expected Load to be called on incoming scene")
	}
}

func TestSwitchSceneToUnloadsPrevious(t *testing.T) {
	m := NewManager(context.Background())
	sc1 := &stubScene{name: "menu"}
	sc2 := &stubScene{name: "game"}
	m.Add("menu", sc1)
	m.Add("game", sc2)

	if err := m.SwitchSceneTo("menu"); err != nil {
		t.Fatal(err)
	}

	if err := m.SwitchSceneTo("game"); err != nil {
		t.Fatal(err)
	}

	if !sc1.unloaded {
		t.Fatal("expected previous scene to be unloaded")
	}

	if m.Scene() != sc2 {
		t.Fatal("expected active scene to be 'game'")
	}
}

func TestSwitchSceneToNotFound(t *testing.T) {
	m := NewManager(context.Background())
	err := m.SwitchSceneTo("missing")

	if err != ErrSceneNotFound {
		t.Fatalf("expected ErrSceneNotFound, got %v", err)
	}
}

func TestSwitchMultipleTimes(t *testing.T) {
	m := NewManager(context.Background())
	for _, name := range []string{"a", "b", "c"} {
		m.Add(name, &stubScene{name: name})
	}

	for _, name := range []string{"a", "b", "c", "a"} {
		if err := m.SwitchSceneTo(name); err != nil {
			t.Fatalf("switch to %s: %v", name, err)
		}

		if m.Scene().Name() != name {
			t.Fatalf("expected scene %s, got %s", name, m.Scene().Name())
		}
	}
}
