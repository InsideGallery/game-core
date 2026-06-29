package event

import "sync"

// Event carries a topic and arbitrary data.
type Event struct {
	Topic string
	Data  any
}

// Handler is a function called synchronously when a matching event is published.
type Handler func(Event)

// Bus is a synchronous in-process pub/sub bus.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// NewBus creates an empty Bus.
func NewBus() *Bus {
	return &Bus{subs: make(map[string][]Handler)}
}

// Subscribe registers h to be called for every event with the given topic.
func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[topic] = append(b.subs[topic], h)
}

// Publish calls all handlers registered for e.Topic synchronously.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subs[e.Topic]))
	copy(handlers, b.subs[e.Topic])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}
