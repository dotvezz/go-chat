package message

import "github.com/dotvezz/gochat/chat"

// Function signatures for message usecases
type Fetch func(line int) (Message, error)
type FetchN func(start, length int) ([]Message, error)
type FetchNByUsername func(username string, start, length int) ([]Message, error)
type Post func(message chat.Message)
