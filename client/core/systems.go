package core

import "sync"

// Systems is an ordered, named collection of ECS systems.
type Systems struct {
	mu    sync.RWMutex
	list  map[string]System
	order []string
}

func NewSystems() *Systems {
	return &Systems{
		list: map[string]System{},
	}
}

// Add registers a system with the given name. Execution order follows registration order.
// If the name is already registered, only the stored system is updated (order unchanged).
func (s *Systems) Add(name string, system System) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.list[name]; !exists {
		s.order = append(s.order, name)
	}

	s.list[name] = system
}

// Remove unregisters a system by name.
func (s *Systems) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.order {
		if n == name {
			s.order = append(s.order[:i], s.order[i+1:]...)

			break
		}
	}

	delete(s.list, name)
}

// Get returns all systems in registration order.
func (s *Systems) Get() []System {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]System, len(s.order))
	for i, name := range s.order {
		result[i] = s.list[name]
	}

	return result
}

// Clean removes all systems.
func (s *Systems) Clean() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.order = nil
	s.list = map[string]System{}
}
