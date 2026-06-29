package scene

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"
)

// Scene is a self-contained game screen with its own ECS.
type Scene interface {
	Name() string
	Load(ctx context.Context)
	Unload(ctx context.Context)
	Update(ctx context.Context) error
	Draw(screen *ebiten.Image)
}
