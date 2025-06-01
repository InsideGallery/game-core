package relations

import (
	"github.com/InsideGallery/core/ecs"
)

// Child describe child entity
type Child interface {
	ecs.Entity
	Parent(entityType string) (Parent, error)
	SetParent(entityType string, entityID uint64)
	RemParent(entityType string)
	Construct() error
	Destroy() error
}

// Parent describe parent entity
type Parent interface {
	ecs.Entity
	Attach(entityType string, entity Child)
	Detach(entityType string, entity Child)
	GetChildren(entityType string) []Child
	GetAllChildren() map[string][]Child
}
