package tracker

// Uses chat.Tracker to Implement business logic related to users

import (
	"github.com/dotvezz/go-chat/chat"
	"github.com/dotvezz/go-chat/chat/domain/message"
)

// PostMessage builds an returns an implementation of the message.Post usecase
// The implementation posts the message to the chat.Tracker, an injected dependency
func PostMessage(tracker chat.Tracker) message.Post {
	return func(message chat.Message) {
		tracker.Post(message)
	}
}
