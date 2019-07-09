package message

import "errors"

// ErrNotFound is returned by message-related functions when a message could not be found
var ErrNotFound = errors.New("couldn't find the message specified")
