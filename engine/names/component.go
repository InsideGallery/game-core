package names

import "sync"

// NameComponent helper for naming entities
type NameComponent struct {
	name string

	mu sync.RWMutex
}

// NewNameComponent return named component
func NewNameComponent(name string) *NameComponent {
	return &NameComponent{
		name: name,
	}
}

// SetName set name
func (n *NameComponent) SetName(name string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.name = name
}

// GetName get name
func (n *NameComponent) GetName() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	name := n.name

	return name
}
