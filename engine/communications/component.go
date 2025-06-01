package communications

import (
	"context"
	"log"
	"net"
	"sync"
)

const bufferSize = 1000

// CommunicateComponent communication component
type CommunicateComponent struct {
	conn          net.Conn
	incoming      chan []byte
	outgoing      chan []byte
	outgoingQueue [][]byte
	parser        CommandParser
	wait          bool

	mu sync.RWMutex
}

// NewCommunicateComponent return new communication component
func NewCommunicateComponent(conn net.Conn) *CommunicateComponent {
	c := &CommunicateComponent{
		conn: conn,
	}
	c.InitChannels()

	return c
}

// InitChannels init channels
func (c *CommunicateComponent) InitChannels() {
	c.incoming = make(chan []byte, bufferSize)
	c.outgoing = make(chan []byte, bufferSize)
}

// Close close channels
func (c *CommunicateComponent) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrChannelsAlreadyClosed
		}
	}()

	c.close()

	if c.conn != nil {
		err = c.conn.Close()
	}

	return
}

func (c *CommunicateComponent) close() {
	close(c.incoming)
	close(c.outgoing)
}

// Wait mark connection waiting
func (c *CommunicateComponent) Wait(w bool) {
	c.mu.Lock()
	c.wait = w
	c.mu.Unlock()
}

// IsWaiting return true if still wait
func (c *CommunicateComponent) IsWaiting() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	w := c.wait

	return w
}

// GetOutgoing return outgoing channel
func (c *CommunicateComponent) GetOutgoing() chan []byte {
	return c.outgoing
}

// Write write data in chan
func (c *CommunicateComponent) Write(d []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in CommunicateComponent.Write", "panic", r)
		}
	}()

	if c.IsWaiting() {
		return
	}
	c.outgoing <- d
}

// GetIncoming return incoming channel
func (c *CommunicateComponent) GetIncoming() chan []byte {
	return c.incoming
}

// AddMessageToQueue add message to queue
func (c *CommunicateComponent) AddMessageToQueue(m OutgoingMessage) {
	c.mu.Lock()
	c.outgoingQueue = append(c.outgoingQueue, m.Encode())
	c.mu.Unlock()
}

// GetQueue retun copy of queue
func (c *CommunicateComponent) GetQueue() [][]byte {
	c.mu.Lock()
	defer c.mu.Unlock()

	d := make([][]byte, len(c.outgoingQueue))
	copy(d, c.outgoingQueue)
	c.outgoingQueue = [][]byte{}

	return d
}

// ProcessOutgoingQueue processing cache
func (c *CommunicateComponent) ProcessOutgoingQueue() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ProcessOutgoingQueue", "panic", r)
		}
	}()

	d := c.GetQueue()
	for _, e := range d {
		c.Write(e)
	}
}

// SetParser set command parser
func (c *CommunicateComponent) SetParser(parser CommandParser) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.parser = parser
}

// GetParser return command parser
func (c *CommunicateComponent) GetParser() CommandParser {
	c.mu.RLock()
	defer c.mu.RUnlock()

	p := c.parser

	return p
}

// ProcessIncomingMessages execute commands
func (c *CommunicateComponent) ProcessIncomingMessages(ctx context.Context, e []byte) error {
	parser := c.GetParser()
	if parser == nil {
		return nil
	}

	cmd, err := c.GetParser().Parse(e)
	if err != nil {
		return err
	}

	return cmd.Execute(ctx)
}

// StartReadingMessages starting processing incoming messages
func (c *CommunicateComponent) StartReadingMessages(ctx context.Context) {
	go func() {
		for e := range c.GetIncoming() {
			err := c.ProcessIncomingMessages(ctx, e)
			if err != nil {
				log.Println("Error processing", "err", err)
				continue
			}
		}
	}()
}

// Send send message for player
func (c *CommunicateComponent) Send(d OutgoingMessage) {
	if c.IsWaiting() {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in Send", "panic", r)
		}
	}()

	c.GetOutgoing() <- d.Encode()
}
