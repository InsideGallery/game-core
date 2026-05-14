# Communications

Import path: `github.com/InsideGallery/game-core/engine/communications`

Package `communications` provides message and command contracts plus a communication component backed by incoming
and outgoing byte channels. It also includes a system that flushes queued outgoing messages for communication
entities stored in FrogoAI memory registry groups.

Key exports:

- `Command` describes decodable, encodable commands that execute with a context.
- `CommandParser` parses raw message bytes into commands.
- `OutgoingMessage` describes encodable outbound messages.
- `Communication` is the behavior required by the communication system.
- `CommunicateComponent` owns incoming/outgoing channels, an optional parser, an optional `net.Conn`, and an
  outgoing message queue.
- `NewCommunicateComponent` creates a component and initializes its channels.
- `Wait` and `IsWaiting` control whether sends and writes are skipped.
- `AddMessageToQueue`, `GetQueue`, and `ProcessOutgoingQueue` manage queued outbound messages.
- `ProcessIncomingMessages` parses and executes inbound command bytes when a parser is configured.
- `StartReadingMessages` starts a goroutine that processes messages from the incoming channel.
- `CommunicationSystem.Update` iterates configured registry groups and processes outgoing queues.
- `ErrChannelsAlreadyClosed` is returned when component channels are closed more than once.
