package message

import "errors"

// NotFound is returned by message-related functions when a message could not be found
var NotFound = errors.New("couldn't find the message specified")
