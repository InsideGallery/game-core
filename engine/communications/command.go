package communications

import "context"

// Command interface for standard command
type Command interface {
	GetMsgType() uint8
	Decode(msg []byte)
	Encode() []byte
	Execute(ctx context.Context) error
}

// CommandParser command parser
type CommandParser interface {
	Parse(msg []byte) (cmd Command, err error)
}
