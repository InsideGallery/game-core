package communications

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/InsideGallery/core/memory/registry"
	"github.com/InsideGallery/core/testutils"
)

// mockOutgoingMessage implements OutgoingMessage
type mockOutgoingMessage struct {
	msgType uint8
	data    []byte
}

func (m *mockOutgoingMessage) GetMessageType() uint8 {
	return m.msgType
}

func (m *mockOutgoingMessage) Encode() []byte {
	return m.data
}

// mockCommand implements Command
type mockCommand struct {
	msgType    uint8
	data       []byte
	executeErr error
	executed   atomic.Bool
}

func (m *mockCommand) GetMsgType() uint8 {
	return m.msgType
}

func (m *mockCommand) Decode(msg []byte) {
	m.data = msg
}

func (m *mockCommand) Encode() []byte {
	return m.data
}

func (m *mockCommand) Execute(_ context.Context) error {
	m.executed.Store(true)
	return m.executeErr
}

// mockCommandParser implements CommandParser
type mockCommandParser struct {
	cmd Command
	err error
}

func (m *mockCommandParser) Parse(_ []byte) (Command, error) {
	return m.cmd, m.err
}

// switchingMockParser fails on the first call and succeeds on subsequent calls
type switchingMockParser struct {
	callCount *atomic.Int32
	failErr   error
	cmd       *mockCommand
}

func (s *switchingMockParser) Parse(_ []byte) (Command, error) {
	n := s.callCount.Add(1)
	if n == 1 {
		return nil, s.failErr
	}

	return s.cmd, nil
}

func TestNewCommunicateComponent(t *testing.T) {
	c := NewCommunicateComponent(nil)
	if c == nil {
		t.Fatal("NewCommunicateComponent returned nil")
	}
	if c.incoming == nil {
		t.Fatal("incoming channel is nil")
	}
	if c.outgoing == nil {
		t.Fatal("outgoing channel is nil")
	}
}

func TestInitChannels(t *testing.T) {
	c := &CommunicateComponent{}
	c.InitChannels()

	if c.incoming == nil {
		t.Fatal("incoming channel is nil after InitChannels")
	}
	if c.outgoing == nil {
		t.Fatal("outgoing channel is nil after InitChannels")
	}
}

func TestClose(t *testing.T) {
	c := NewCommunicateComponent(nil)
	err := c.Close()
	testutils.Equal(t, err, nil)
}

func TestCloseDoubleCloseReturnsError(t *testing.T) {
	c := NewCommunicateComponent(nil)
	err := c.Close()
	testutils.Equal(t, err, nil)

	err = c.Close()
	testutils.Equal(t, err, ErrChannelsAlreadyClosed)
}

func TestWaitAndIsWaiting(t *testing.T) {
	c := NewCommunicateComponent(nil)

	testutils.Equal(t, c.IsWaiting(), false)

	c.Wait(true)
	testutils.Equal(t, c.IsWaiting(), true)

	c.Wait(false)
	testutils.Equal(t, c.IsWaiting(), false)
}

func TestGetOutgoing(t *testing.T) {
	c := NewCommunicateComponent(nil)
	ch := c.GetOutgoing()
	if ch == nil {
		t.Fatal("GetOutgoing returned nil channel")
	}
}

func TestGetIncoming(t *testing.T) {
	c := NewCommunicateComponent(nil)
	ch := c.GetIncoming()
	if ch == nil {
		t.Fatal("GetIncoming returned nil channel")
	}
}

func TestWriteSendsToOutgoing(t *testing.T) {
	c := NewCommunicateComponent(nil)
	data := []byte("hello")

	c.Write(data)

	select {
	case received := <-c.GetOutgoing():
		testutils.Equal(t, string(received), string(data))
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for data on outgoing channel")
	}
}

func TestWriteSkipsWhenWaiting(t *testing.T) {
	c := NewCommunicateComponent(nil)
	c.Wait(true)

	c.Write([]byte("should not send"))

	select {
	case <-c.GetOutgoing():
		t.Fatal("Write should not send data when waiting")
	case <-time.After(50 * time.Millisecond):
		// expected: nothing received
	}
}

func TestWriteOnClosedChannelDoesNotPanic(t *testing.T) {
	c := NewCommunicateComponent(nil)
	c.Close()

	// Should not panic due to defer recover
	c.Write([]byte("after close"))
}

func TestAddMessageToQueueAndGetQueue(t *testing.T) {
	c := NewCommunicateComponent(nil)
	msg1 := &mockOutgoingMessage{msgType: 1, data: []byte("msg1")}
	msg2 := &mockOutgoingMessage{msgType: 2, data: []byte("msg2")}

	c.AddMessageToQueue(msg1)
	c.AddMessageToQueue(msg2)

	queue := c.GetQueue()
	testutils.Equal(t, len(queue), 2)
	testutils.Equal(t, string(queue[0]), "msg1")
	testutils.Equal(t, string(queue[1]), "msg2")
}

func TestGetQueueClearsQueue(t *testing.T) {
	c := NewCommunicateComponent(nil)
	msg := &mockOutgoingMessage{msgType: 1, data: []byte("msg")}

	c.AddMessageToQueue(msg)
	queue := c.GetQueue()
	testutils.Equal(t, len(queue), 1)

	queue2 := c.GetQueue()
	testutils.Equal(t, len(queue2), 0)
}

func TestGetQueueEmptyInitially(t *testing.T) {
	c := NewCommunicateComponent(nil)
	queue := c.GetQueue()
	testutils.Equal(t, len(queue), 0)
}

func TestProcessOutgoingQueue(t *testing.T) {
	c := NewCommunicateComponent(nil)
	msg1 := &mockOutgoingMessage{msgType: 1, data: []byte("queued1")}
	msg2 := &mockOutgoingMessage{msgType: 2, data: []byte("queued2")}

	c.AddMessageToQueue(msg1)
	c.AddMessageToQueue(msg2)

	c.ProcessOutgoingQueue()

	received := make([]string, 0, 2)
	for i := 0; i < 2; i++ {
		select {
		case d := <-c.GetOutgoing():
			received = append(received, string(d))
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for queued message on outgoing channel")
		}
	}

	testutils.Equal(t, len(received), 2)
	testutils.Equal(t, received[0], "queued1")
	testutils.Equal(t, received[1], "queued2")
}

func TestProcessOutgoingQueueEmpty(t *testing.T) {
	c := NewCommunicateComponent(nil)
	// Should not panic or block
	c.ProcessOutgoingQueue()
}

func TestSetParserAndGetParser(t *testing.T) {
	c := NewCommunicateComponent(nil)

	testutils.Equal(t, c.GetParser(), nil)

	parser := &mockCommandParser{}
	c.SetParser(parser)

	got := c.GetParser()
	if got == nil {
		t.Fatal("GetParser returned nil after SetParser")
	}
}

func TestProcessIncomingMessagesNilParser(t *testing.T) {
	c := NewCommunicateComponent(nil)

	err := c.ProcessIncomingMessages(context.Background(), []byte("data"))
	testutils.Equal(t, err, nil)
}

func TestProcessIncomingMessagesWithParser(t *testing.T) {
	c := NewCommunicateComponent(nil)
	cmd := &mockCommand{msgType: 1}
	parser := &mockCommandParser{cmd: cmd}
	c.SetParser(parser)

	err := c.ProcessIncomingMessages(context.Background(), []byte("data"))
	testutils.Equal(t, err, nil)
	testutils.Equal(t, cmd.executed.Load(), true)
}

func TestProcessIncomingMessagesParserError(t *testing.T) {
	c := NewCommunicateComponent(nil)
	parseErr := errors.New("parse error")
	parser := &mockCommandParser{err: parseErr}
	c.SetParser(parser)

	err := c.ProcessIncomingMessages(context.Background(), []byte("bad data"))
	if err == nil {
		t.Fatal("expected error from ProcessIncomingMessages, got nil")
	}
	testutils.Equal(t, err.Error(), "parse error")
}

func TestProcessIncomingMessagesCommandExecuteError(t *testing.T) {
	c := NewCommunicateComponent(nil)
	execErr := errors.New("execute error")
	cmd := &mockCommand{msgType: 1, executeErr: execErr}
	parser := &mockCommandParser{cmd: cmd}
	c.SetParser(parser)

	err := c.ProcessIncomingMessages(context.Background(), []byte("data"))
	if err == nil {
		t.Fatal("expected error from command Execute, got nil")
	}
	testutils.Equal(t, err.Error(), "execute error")
}

func TestSendEncodesAndSendsToOutgoing(t *testing.T) {
	c := NewCommunicateComponent(nil)
	msg := &mockOutgoingMessage{msgType: 1, data: []byte("encoded")}

	c.Send(msg)

	select {
	case received := <-c.GetOutgoing():
		testutils.Equal(t, string(received), "encoded")
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for sent message on outgoing channel")
	}
}

func TestSendSkipsWhenWaiting(t *testing.T) {
	c := NewCommunicateComponent(nil)
	c.Wait(true)
	msg := &mockOutgoingMessage{msgType: 1, data: []byte("should not send")}

	c.Send(msg)

	select {
	case <-c.GetOutgoing():
		t.Fatal("Send should not send data when waiting")
	case <-time.After(50 * time.Millisecond):
		// expected: nothing received
	}
}

func TestSendOnClosedChannelDoesNotPanic(t *testing.T) {
	c := NewCommunicateComponent(nil)
	c.Close()
	msg := &mockOutgoingMessage{msgType: 1, data: []byte("after close")}

	// Should not panic due to defer recover
	c.Send(msg)
}

func TestStartReadingMessages(t *testing.T) {
	c := NewCommunicateComponent(nil)
	cmd := &mockCommand{msgType: 1}
	parser := &mockCommandParser{cmd: cmd}
	c.SetParser(parser)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.StartReadingMessages(ctx)

	c.GetIncoming() <- []byte("test message")

	// Poll with timeout for the atomic flag
	deadline := time.After(time.Second)
	for {
		if cmd.executed.Load() {
			break
		}

		select {
		case <-deadline:
			t.Fatal("timeout waiting for command execution")
		default:
			time.Sleep(5 * time.Millisecond) //nolint:mnd
		}
	}
}

func TestStartReadingMessagesStopsOnChannelClose(t *testing.T) {
	c := NewCommunicateComponent(nil)
	parser := &mockCommandParser{cmd: &mockCommand{}}
	c.SetParser(parser)

	ctx := context.Background()
	c.StartReadingMessages(ctx)

	// Close channels - goroutine should exit when incoming is closed
	c.Close()

	// Give the goroutine time to finish
	time.Sleep(50 * time.Millisecond)
}

func TestProcessOutgoingQueueOnClosedChannelDoesNotPanic(t *testing.T) {
	c := NewCommunicateComponent(nil)
	msg := &mockOutgoingMessage{msgType: 1, data: []byte("queued")}
	c.AddMessageToQueue(msg)
	c.Close()

	// Should not panic due to defer recover in ProcessOutgoingQueue -> Write
	c.ProcessOutgoingQueue()
}

func TestCloseWithConnection(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	c := NewCommunicateComponent(client)
	err := c.Close()
	testutils.Equal(t, err, nil)
}

func TestStartReadingMessagesWithError(t *testing.T) {
	c := NewCommunicateComponent(nil)
	parseErr := errors.New("parse failure")

	// Use a parser that fails on first call, succeeds on second
	callCount := atomic.Int32{}
	cmd := &mockCommand{msgType: 1}
	switchingParser := &switchingMockParser{
		callCount: &callCount,
		failErr:   parseErr,
		cmd:       cmd,
	}
	c.SetParser(switchingParser)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.StartReadingMessages(ctx)

	// First message triggers parse error - goroutine should continue
	c.GetIncoming() <- []byte("bad message")

	// Wait until first call is processed
	deadline := time.After(time.Second)
	for callCount.Load() < 1 {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for first parse call")
		default:
			time.Sleep(5 * time.Millisecond) //nolint:mnd
		}
	}

	// Second message succeeds
	c.GetIncoming() <- []byte("good message")

	deadline = time.After(time.Second)
	for !cmd.executed.Load() {
		select {
		case <-deadline:
			t.Fatal("timeout waiting for command execution")
		default:
			time.Sleep(5 * time.Millisecond) //nolint:mnd
		}
	}
}

// Tests for system.go

func TestNewCommunicationSystem(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	sys := NewCommunicationSystem(reg, 1, "key1", "key2")

	if sys == nil {
		t.Fatal("NewCommunicationSystem returned nil")
	}
}

func TestEntitiesKeys(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	sys := NewCommunicationSystem(reg, 1, "key1", "key2")

	keys := sys.EntitiesKeys()
	testutils.Equal(t, len(keys), 2)
	testutils.Equal(t, keys[0], "key1")
	testutils.Equal(t, keys[1], "key2")
}

func TestEntitiesKeysEmpty(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	sys := NewCommunicationSystem(reg, 1)

	keys := sys.EntitiesKeys()
	testutils.Equal(t, len(keys), 0)
}

// mockCommunication implements Communication for testing Update
type mockCommunication struct {
	processed atomic.Bool
}

func (m *mockCommunication) ProcessOutgoingQueue()                                       { m.processed.Store(true) }
func (m *mockCommunication) AddMessageToQueue(_ OutgoingMessage)                         {}
func (m *mockCommunication) GetIncoming() chan []byte                                    { return nil }
func (m *mockCommunication) GetOutgoing() chan []byte                                    { return nil }
func (m *mockCommunication) ProcessIncomingMessages(_ context.Context, _ []byte) error   { return nil }
func (m *mockCommunication) Close() error                                                { return nil }
func (m *mockCommunication) Write(_ []byte)                                              {}

func TestUpdateWithCommunicationEntity(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	groupKey := "players"
	sys := NewCommunicationSystem(reg, 1, groupKey)

	mc := &mockCommunication{}
	err := reg.Add(groupKey, "player1", mc)
	testutils.Equal(t, err, nil)

	err = sys.Update(context.Background())
	testutils.Equal(t, err, nil)
	testutils.Equal(t, mc.processed.Load(), true)
}

func TestUpdateWithNonCommunicationEntity(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	groupKey := "objects"
	sys := NewCommunicationSystem(reg, 1, groupKey)

	// Add a non-Communication entity (just a string)
	err := reg.Add(groupKey, "obj1", "not a communication")
	testutils.Equal(t, err, nil)

	// Should not panic, should just skip the entity
	err = sys.Update(context.Background())
	testutils.Equal(t, err, nil)
}

func TestUpdateEmptyRegistry(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	sys := NewCommunicationSystem(reg, 1, "empty_group")

	err := sys.Update(context.Background())
	testutils.Equal(t, err, nil)
}

func TestUpdateMultipleGroups(t *testing.T) {
	reg := registry.NewRegistry[any, any, any]()
	sys := NewCommunicationSystem(reg, 1, "g1", "g2")

	mc1 := &mockCommunication{}
	mc2 := &mockCommunication{}
	err := reg.Add("g1", "p1", mc1)
	testutils.Equal(t, err, nil)
	err = reg.Add("g2", "p2", mc2)
	testutils.Equal(t, err, nil)

	err = sys.Update(context.Background())
	testutils.Equal(t, err, nil)
	testutils.Equal(t, mc1.processed.Load(), true)
	testutils.Equal(t, mc2.processed.Load(), true)
}
