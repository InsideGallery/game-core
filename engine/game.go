package engine

import (
	"context"

	"github.com/InsideGallery/game-core/engine/communications"

	"github.com/InsideGallery/core/ecs"
)

// Game describe game object
type Game interface {
	Initialize() error
	Tick(ctx context.Context)
}

// Player describe player entity
type Player interface {
	ecs.Entity
	communications.Communication
}
