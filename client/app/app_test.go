package app

import (
	"context"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/InsideGallery/game-core/client/event"
	"github.com/InsideGallery/game-core/client/scene"
)

func TestLayoutReturnsLogicalSize(t *testing.T) {
	cfg := Config{Width: 800, Height: 600, Title: "test"}
	a := New(cfg)

	w, h := a.Layout(1920, 1080)
	if w != 800 || h != 600 {
		t.Fatalf("Layout: want 800x600, got %dx%d", w, h)
	}
}

func TestLayoutDifferentSizes(t *testing.T) {
	tests := []struct{ w, h int }{
		{320, 240},
		{1280, 720},
		{1920, 1080},
	}

	for _, tt := range tests {
		a := New(Config{Width: tt.w, Height: tt.h})
		gw, gh := a.Layout(0, 0)
		if gw != tt.w || gh != tt.h {
			t.Errorf("Layout(%d,%d): want %dx%d, got %dx%d", tt.w, tt.h, tt.w, tt.h, gw, gh)
		}
	}
}

func TestSetupCalledAndSceneRegistered(t *testing.T) {
	const sceneName = "main"
	setupCalled := false

	cfg := Config{
		Width:  100,
		Height: 100,
		Setup: func(_ context.Context, _ *event.Bus, m *scene.Manager, _ func(string)) string {
			setupCalled = true
			m.Add(sceneName, &stubScene{name: sceneName})
			return sceneName
		},
	}

	a := New(cfg)
	if !setupCalled {
		t.Fatal("Setup was not called")
	}

	sc := a.manager.Scene()
	if sc == nil {
		t.Fatal("expected active scene after Setup, got nil")
	}

	if sc.Name() != sceneName {
		t.Fatalf("active scene: want %q, got %q", sceneName, sc.Name())
	}
}

func TestNewNilSetup(t *testing.T) {
	a := New(Config{Width: 10, Height: 10})
	if a == nil {
		t.Fatal("expected non-nil App")
	}

	if a.manager.Scene() != nil {
		t.Fatal("expected nil active scene with no Setup")
	}
}

func TestUpdateNilScene(t *testing.T) {
	a := New(Config{Width: 10, Height: 10})
	if err := a.Update(); err != nil {
		t.Fatalf("Update with no scene: %v", err)
	}
}

func TestLayoutPropagatesViewPort(t *testing.T) {
	const w, h = 640, 480

	spy := &spyViewPortSetter{stubScene: stubScene{name: "spy"}}
	cfg := Config{
		Width:  w,
		Height: h,
		Setup: func(_ context.Context, _ *event.Bus, m *scene.Manager, _ func(string)) string {
			m.Add("spy", spy)
			return "spy"
		},
	}

	a := New(cfg)
	a.Layout(1024, 768)

	if !spy.setViewPortCalled {
		t.Fatal("expected SetViewPort to be called on active scene")
	}

	if spy.vpW != w || spy.vpH != h {
		t.Fatalf("SetViewPort: want %dx%d, got %dx%d", w, h, spy.vpW, spy.vpH)
	}
}

// stubScene is a minimal Scene for testing.
type stubScene struct {
	name string
}

func (s *stubScene) Name() string                    { return s.name }
func (s *stubScene) Load(_ context.Context)          {}
func (s *stubScene) Unload(_ context.Context)        {}
func (s *stubScene) Update(_ context.Context) error  { return nil }
func (s *stubScene) Draw(_ *ebiten.Image) {}

// spyViewPortSetter wraps stubScene and implements viewPortSetter.
type spyViewPortSetter struct {
	stubScene
	setViewPortCalled bool
	vpW, vpH          int
}

func (s *spyViewPortSetter) SetViewPort(w, h int) {
	s.setViewPortCalled = true
	s.vpW = w
	s.vpH = h
}
