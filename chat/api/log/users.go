package log

// Uses the log to Implement business logic related to users

import (
	"bufio"
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/domain/user"
	"os"
)

// FetchUser builds and returns an implementation of the user.Fetch usecase
// The implementation scans the log for the given username
// If any user has sent a message while using that name, it will return a User,
// Otherwise it will return an error
// The implementation will also error for log file access problems.
// The implementation also takes a usecase dependency of type user.GetBool to determine
// if the user is online, to decouple this particular logic from the tracker
func FetchUser(logFilePath string, isOnline user.GetBool) user.Fetch {
	return func(userName string) (user.User, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return user.User{}, err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		m := chat.Message{}
		for sc.Scan() {
			err := json.Unmarshal(sc.Bytes(), &m)
			if err != nil {
				continue
			}
			if m.From == userName {
				return user.User{Name: userName, Online: isOnline(userName)}, nil
			}
		}

		return user.User{}, user.ErrNotFound
	}
}

// FetchAllUsers builds and returns an implementation of the user.FetchAll usecase
// The implementation scans the log directly and returns all users who have sent messages
// The implementation will return an error only when it fails to open the log
func FetchAllUsers(logFilePath string) user.FetchAll {
	return func() ([]user.User, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return make([]user.User, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		m := chat.Message{}
		userNameMap := make(map[string]bool)
		us := make([]user.User, 0)
		for sc.Scan() {
			err := json.Unmarshal(sc.Bytes(), &m)
			if err != nil {
				continue
			}
			if _, ok := userNameMap[m.From]; !ok {
				userNameMap[m.From] = true
				us = append(us, user.User{Name: m.From})
			}
		}
		return us, nil
	}
}
