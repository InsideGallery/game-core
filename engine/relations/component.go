package relations

import (
	"sync"

	"github.com/InsideGallery/core/memory/registry"
)

var store = registry.NewRegistry[string, uint64, any]()

// RelationComponent contains parents and childs
type RelationComponent struct {
	parent        map[string]uint64
	child         map[string][]Child
	ignoreZeroIDs bool

	mu sync.RWMutex
}

// NewRelationComponent return new ralation component
func NewRelationComponent(parents map[string]uint64, ignoreZeroIDs bool) *RelationComponent {
	return &RelationComponent{
		parent:        parents,
		child:         make(map[string][]Child),
		ignoreZeroIDs: ignoreZeroIDs,
	}
}

// SetParent set parent by type
func (r *RelationComponent) SetParent(entityType string, entityID uint64) {
	r.mu.Lock()
	r.parent[entityType] = entityID
	r.mu.Unlock()
}

// RemParent remove parent by type
func (r *RelationComponent) RemParent(entityType string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.parent, entityType)
}

// GetParentID return parent id
func (r *RelationComponent) GetParentID(entityType string) uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id := r.parent[entityType]

	return id
}

// Parent return parent entity
func (r *RelationComponent) Parent(entityType string) (Parent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if id, e := r.parent[entityType]; e && id != 0 {
		entity, err := store.Get(entityType, id)
		if err != nil {
			return nil, err
		}

		return entity.(Parent), err
	}

	return nil, ErrNotFoundParent
}

// ConstructChild construct relation with parent
func (r *RelationComponent) ConstructChild(childType string, entity Child) error {
	for t := range r.parent {
		parent, err := r.Parent(t)
		if r.ignoreZeroIDs && err == ErrNotFoundParent {
			continue
		} else if err != nil {
			return err
		}

		parent.Attach(childType, entity)
	}

	return nil
}

// GetParents return parents
func (r *RelationComponent) GetParents() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p := map[string]uint64{}
	for k, v := range r.parent {
		p[k] = v
	}

	return p
}

// DestroyChild destroy relation with parent
func (r *RelationComponent) DestroyChild(childType string, entity Child) error {
	for t := range r.GetParents() {
		parent, err := r.Parent(t)
		if r.ignoreZeroIDs && err == ErrNotFoundParent {
			continue
		} else if err != nil {
			return err
		}

		parent.Detach(childType, entity)
	}

	return nil
}

// Attach attach child
func (r *RelationComponent) Attach(entityType string, c Child) {
	r.mu.Lock()
	r.child[entityType] = append(r.child[entityType], c)
	r.mu.Unlock()
}

// Detach detach child
func (r *RelationComponent) Detach(entityType string, c Child) {
	r.mu.Lock()
	for i, v := range r.child[entityType] {
		if v == c {
			r.child[entityType] = append(r.child[entityType][:i], r.child[entityType][i+1:]...)
			break
		}
	}
	r.mu.Unlock()
}

// GetChildren return children
func (r *RelationComponent) GetChildren(entityType string) []Child {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if c, e := r.child[entityType]; e {
		children := make([]Child, len(r.child[entityType]))
		copy(children, c)

		return children
	}

	return []Child{}
}

// GetAllChildren return all children
func (r *RelationComponent) GetAllChildren() map[string][]Child {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c := map[string][]Child{}
	for k, v := range r.child {
		c[k] = make([]Child, len(v))
		copy(c[k], v)
	}

	return c
}
