package communications

import (
	"context"

	"github.com/InsideGallery/core/memory/registry"
	"github.com/InsideGallery/core/multiproc/worker"
)

// OutgoingMessage describe outgoing message
type OutgoingMessage interface {
	GetMessageType() uint8
	Encode() []byte
}

// Communication describe communication
type Communication interface {
	ProcessOutgoingQueue()
	AddMessageToQueue(m OutgoingMessage)
	GetIncoming() chan []byte
	GetOutgoing() chan []byte
	ProcessIncomingMessages(ctx context.Context, e []byte) error
	Close() error
	Write(d []byte)
}

// CommunicationSystem contains moveable entities
type CommunicationSystem struct {
	keys         []interface{}
	workersCount int
	reg          *registry.Registry[any, any, any]
}

// NewCommunicationSystem return new CommunicationSystem
func NewCommunicationSystem(
	reg *registry.Registry[any, any, any],
	workersCount int,
	keys ...interface{},
) *CommunicationSystem {
	return &CommunicationSystem{
		workersCount: workersCount,
		keys:         keys,
		reg:          reg,
	}
}

// EntitiesKeys return entities keys
func (c *CommunicationSystem) EntitiesKeys() []interface{} {
	return c.keys
}

// Update update move
func (c *CommunicationSystem) Update(ctx context.Context) error {
	for _, key := range c.EntitiesKeys() {
		g := c.reg.GetGroup(key).Iterator()

		worker.RunSyncMultipleWorkers(ctx, c.workersCount, func(_ context.Context) {
			for e := range g {
				m, ok := e.(Communication)
				if !ok {
					continue
				}

				m.ProcessOutgoingQueue()
			}
		})
	}

	return nil
}
