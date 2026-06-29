package core

import (
	"context"

	"github.com/InsideGallery/core/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// System is an ECS system with both Update and Draw capabilities.
type System interface {
	ecs.System
	Draw(ctx context.Context, screen *ebiten.Image)
}

// SystemWindow is a system that draws UI overlays in screen space.
type SystemWindow interface {
	System
	ScreenDraw(ctx context.Context, screen *ebiten.Image)
}
