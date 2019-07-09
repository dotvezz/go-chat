package user

import "errors"

// ErrNotFound is returned by user-related functions when a user could not be found
var ErrNotFound = errors.New("couldn't find the user specified")
