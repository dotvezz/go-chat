package log

// Uses the log to Implement business logic related to users

import (
	"bufio"
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/domain/message"
	"os"
)

// FetchMessage builds and returns an implementation of the message.Fetch usecase
// The message IDs correspond directly to their line number in the log file.
// The implementation will error when it is unable to parse the specified line.
// The implementation will also error for log file access problems.
func FetchMessage(logFilePath string) message.Fetch {
	return func(targetLine int) (message.Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return message.Message{}, err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		cm := chat.Message{}

		// Let's be real here, a log file isn't remotely a structure that we should use for this in real life
		// But it's the data source we've got so let's just go with it.
		for curLine := 1; sc.Scan(); curLine++ {
			if targetLine != curLine {
				continue // Skip any line that isn't the one we want
			}
			err := json.Unmarshal(sc.Bytes(), &cm)
			if err != nil {
				return message.Message{}, err
			}

			m := message.Message{
				ID:        curLine,
				From:      cm.From,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			}

			return m, nil
		}

		return message.Message{}, message.NotFound
	}
}

// FetchMessagesOfSender builds and returns an implementation of the message.FetchN usecase
// The implementation will pull messages directly from the log
// The implementation will return an error only when it fails to open the log
func FetchNMessages(logFilePath string) message.FetchN {
	return func(start, length int) ([]message.Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return make([]message.Message, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		// Unmarshal the log into chat.Message because that's the structure which the log is built from
		cm := chat.Message{}
		ms := make([]message.Message, 0)

		// Let's be real here, a log file isn't remotely a structure that we should use for this in real life
		// But it's the data source we've got so let's just go with it.
		for curLine := 1; sc.Scan(); curLine++ {
			if curLine < start {
				continue // Skip lines before our starting point
			}
			if curLine >= start+length {
				break // Break out of the loop after we've got the length we want
			}
			err := json.Unmarshal(sc.Bytes(), &cm)
			if err != nil {
				continue
			}
			ms = append(ms, message.Message{
				ID:        curLine,
				From:      cm.From,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			})
		}
		return ms, nil
	}
}

// FetchMessagesOfSender builds and returns an implementation of the message.FetchNByUsername usecase
// The implementation function searches the log and returns messages where the given userName is the sender.
// The implementation will return an error only when it fails to open the log
func FetchMessagesOfSender(logFilePath string) message.FetchNByUsername {
	return func(username string, start, length int) ([]message.Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return make([]message.Message, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		// Unmarshal the log into chat.Message because that's the structure which the log is built from
		cm := chat.Message{}
		ms := make([]message.Message, 0)

		counter := 0
		// Let's be real here, a log file isn't remotely a structure that we should use for this in real life
		// But it's the data source we've got so let's just go with it.
		for curLine := 1; sc.Scan(); curLine++ {
			err := json.Unmarshal(sc.Bytes(), &cm)
			if cm.From != username {
				continue // Skip if it's not the user we want
			}
			counter++ // We're only incrementing this for lines that match the user
			if counter < start {
				continue // Skip lines before our starting point
			}
			if counter >= start+length {
				break // Break out of the loop after we've got the length we want
			}
			if err != nil {
				continue
			}
			ms = append(ms, message.Message{
				ID:        curLine,
				From:      cm.From,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			})
		}
		return ms, nil
	}
}
