package users

import (
	"bufio"
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"os"
)

// The basic business entity for a User
type User struct {
	Name string
}

// Function signatures for user usecases
type FetchUser func(userName string) (User, error)
type FetchAllUsers func() ([]User, error)

// Builds and returns a usecase-type function that searches the log for a User.
// If any user has sent a message while using that name, it will return a User,
// Otherwise it will return an error
// The returned function will also error for log file access problems.
func FetchOne(logFilePath string) FetchUser {
	return func(userName string) (User, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDONLY, 0)
			if err != nil {
				return User{}, err
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
				return User{Name: userName}, nil
			}
		}

		return User{}, NotFound
	}
}

// Builds and returns a usecase-type function that gets all users from the log.
// The usecase function will return an error only when  it fails to open the log
func FetchAll(logFilePath string) FetchAllUsers {
	return func() ([]User, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDONLY, 0)
			if err != nil {
				return make([]User, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		m := chat.Message{}
		userNameMap := make(map[string]bool)
		us := make([]User, 0)
		for sc.Scan() {
			err := json.Unmarshal(sc.Bytes(), &m)
			if err != nil {
				continue
			}
			if _, ok := userNameMap[m.From]; !ok {
				userNameMap[m.From] = true
				us = append(us, User{Name: m.From})
			}
		}
		return us, nil
	}
}
