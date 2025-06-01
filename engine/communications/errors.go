package communications

import "errors"

// All kind of errors for game
var (
	ErrChannelsAlreadyClosed = errors.New("error channels already closed")
)
