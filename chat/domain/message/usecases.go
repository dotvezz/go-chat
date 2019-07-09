package message

import "github.com/dotvezz/go-chat/chat"

// Fetch defines a function which takes an int and returns a message, typically the int is the ID of the message
type Fetch func(int) (Message, error)

// FetchN defines a function which takes two ints, canonically declared as `start` and `length` in this case, and
// returns a slice of messages
type FetchN func(start, length int) ([]Message, error)

// FetchNByString defines a function which takes a string and two ints. The two ints are still canonically declared
// as `start` and `length`, and the string is an additional param that might a body search string or sender's name
type FetchNByString func(str string, start, length int) ([]Message, error)

// Post defines a function which takes a message and has no response. The name is meant to suggest its typical intended
// use, posting a message to the tracker
type Post func(message chat.Message)
