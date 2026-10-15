package app

import (
	"context"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/InsideGallery/game-core/client/event"
	"github.com/InsideGallery/game-core/client/scene"
)

// Config configures the App shell.
type Config struct {
	Width, Height int
	Title         string
	// Responsive: when true, Layout ignores Width and Height and returns the outside size times the device scale.
	Responsive bool
	// Setup creates scenes, registers them with the manager, and returns the initial scene name.
	Setup func(ctx context.Context, bus *event.Bus, m *scene.Manager, switchScene func(string)) string
}

// App is the thin Ebiten game-loop host. It owns the Bus and the scene Manager.
type App struct {
	cfg     Config
	manager *scene.Manager
	bus     *event.Bus
	ctx     context.Context
}

// New builds an App, calls Setup to populate the Manager, and loads the initial scene.
func New(cfg Config) *App {
	ctx := context.Background()
	bus := event.NewBus()
	m := scene.NewManager(ctx)

	a := &App{cfg: cfg, manager: m, bus: bus, ctx: ctx}

	if cfg.Setup != nil {
		switchScene := func(name string) {
			if err := m.SwitchSceneTo(name); err != nil {
				slog.Warn("switch scene", "name", name, "error", err)
			}
		}

		initial := cfg.Setup(ctx, bus, m, switchScene)
		if initial != "" {
			if err := m.SwitchSceneTo(initial); err != nil {
				slog.Warn("initial scene switch", "name", initial, "error", err)
			}
		}
	}

	return a
}

// Update delegates to the active scene.
func (a *App) Update() error {
	sc := a.manager.Scene()
	if sc == nil {
		return nil
	}

	return sc.Update(a.ctx)
}

// Draw delegates to the active scene.
func (a *App) Draw(screen *ebiten.Image) {
	sc := a.manager.Scene()
	if sc != nil {
		sc.Draw(screen)
	}
}

// viewPortSetter is an optional interface scenes may implement to receive the logical viewport size.
type viewPortSetter interface {
	SetViewPort(w, h int)
}

// DeviceScale returns the device scale factor of the monitor, or 1 when it is unknown.
func DeviceScale() float64 {
	m := ebiten.Monitor()
	if m == nil {
		return 1
	}

	if f := m.DeviceScaleFactor(); f > 0 {
		return f
	}

	return 1
}

// logicalSize converts an outside size to a logical size by the device scale, rounding up.
func logicalSize(ow, oh int, scale float64) (int, int) {
	return int(math.Ceil(float64(ow) * scale)), int(math.Ceil(float64(oh) * scale))
}

// Layout returns the logical window size and propagates it to the active scene if it implements
// viewPortSetter (e.g. a *BaseScene). With Config.Responsive it follows the outside size times the
// device scale; otherwise it returns Config.Width and Config.Height.
func (a *App) Layout(outsideW, outsideH int) (int, int) {
	w, h := a.cfg.Width, a.cfg.Height
	if a.cfg.Responsive {
		w, h = logicalSize(outsideW, outsideH, DeviceScale())
	}

	sc := a.manager.Scene()
	if sc != nil {
		if vps, ok := sc.(viewPortSetter); ok {
			vps.SetViewPort(w, h)
		}
	}

	return w, h
}

// Run configures the Ebiten window and starts the game loop.
func Run(a *App) error {
	ebiten.SetWindowSize(a.cfg.Width, a.cfg.Height)
	ebiten.SetWindowTitle(a.cfg.Title)

	if a.cfg.Responsive {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}

	return ebiten.RunGame(a)
}
