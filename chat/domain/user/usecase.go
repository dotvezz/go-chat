package user

// Fetch defines a function which returns a user based on string input (typically userName)
type Fetch func(string) (User, error)

// FetchAll defines a function which takes no parameters and returns all users possible
type FetchAll func() ([]User, error)

// Kick defines a function which kicks a user based on a string input (typically userName)
// Technically the same signature as Fetch, but a little specificity in the name is for context
type Kick func(string) (User, error)

// GetBool defines a function which accepts a string and returns a bool.
// One implementation is transport.IsUserOnline, which accepts a username and returns true if the user is online
type GetBool func(string) bool
