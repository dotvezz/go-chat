package user

import "errors"

// NotFound is returned by user-related functions when a user could not be found
var NotFound = errors.New("couldn't find the user specified")
