package tracker

// Uses chat.Tracker to Implement business logic related to users

import (
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/domain/user"
)

// KickUser returns an implementation of the user.Kick usecase.
// The implementation uses a chat.Tracker, injected as a dependency, to find
// open connections and close the one belonging to the given user
func KickUser(tracker chat.Tracker) user.Kick {
	return func(userName string) (user.User, error) {
		for _, c := range tracker.Connections() {
			if c.UserName() == userName {
				c.Close()
				return user.User{Name: userName, Online: false}, nil
			}
		}
		return user.User{}, user.NotFound
	}
}

// IsUserOnline returns an implementation of user.GetBool which returns true of the user is online
// The implementation uses a chat.Tracker, injected as a dependency, to find
func IsUserOnline(tracker chat.Tracker) user.GetBool {
	return func(userName string) bool {
		for _, conn := range tracker.Connections() {
			if conn.UserName() == userName {
				return true
			}
		}

		return false
	}
}
