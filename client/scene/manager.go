package scene

import (
	"context"
	"errors"
	"sync"
)

// ErrSceneNotFound is returned when the named scene is not registered.
var ErrSceneNotFound = errors.New("scene not found")

// Manager manages named scenes and transitions between them.
type Manager struct {
	mu     sync.RWMutex
	scenes map[string]Scene
	active string
	ctx    context.Context
}

// NewManager creates a Manager; ctx is passed to Load and Unload on scene transitions.
func NewManager(ctx context.Context) *Manager {
	return &Manager{
		scenes: make(map[string]Scene),
		ctx:    ctx,
	}
}

// Add registers a scene by name.
func (m *Manager) Add(name string, sc Scene) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.scenes[name] = sc
}

// SwitchSceneTo transitions to the named scene.
// It calls Load on the incoming scene and Unload on the outgoing one.
func (m *Manager) SwitchSceneTo(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sc, ok := m.scenes[name]
	if !ok {
		return ErrSceneNotFound
	}

	if m.active != "" {
		if prev, exists := m.scenes[m.active]; exists {
			prev.Unload(m.ctx)
		}
	}

	sc.Load(m.ctx)
	m.active = name

	return nil
}

// Scene returns the currently active scene, or nil if none is active.
func (m *Manager) Scene() Scene { //nolint:ireturn // manager returns interface by design
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.active == "" {
		return nil
	}

	return m.scenes[m.active]
}
