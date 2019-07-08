package messages

import (
	"bufio"
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"os"
)

// The basic business entity for a Message. Essentially a copy of chat.Message
type Message struct {
	ID        int
	From      string
	To        string
	Body      string
	Timestamp int64
}

// Function signatures for message usecases
type FetchMessage func(line int) (Message, error)
type FetchNMessages func(start, length int) ([]Message, error)
type FetchNByUsername func(username string, start, length int) ([]Message, error)

// Builds and returns a usecase-type function that searches the log for a Message.
// The message IDs correspond directly to their line number in the log file.
// The returned function will error for parsing failures on the specified line.
// The returned function will also error for log file access problems.
func FetchOne(logFilePath string) FetchMessage {
	return func(targetLine int) (Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return Message{}, err
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
				return Message{}, err
			}

			m := Message{
				ID:        curLine,
				From:      cm.From,
				To:        cm.To,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			}

			return m, nil
		}

		return Message{}, NotFound
	}
}

// Builds and returns a usecase-type function that gets all messages from the log.
// The usecase function will return an error only when  it fails to open the log
func FetchN(logFilePath string) FetchNMessages {
	return func(start, length int) ([]Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return make([]Message, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		// Unmarshal the log into chat.Message because that's the structure which the log is built from
		cm := chat.Message{}
		ms := make([]Message, 0)

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
			ms = append(ms, Message{
				ID:        curLine,
				From:      cm.From,
				To:        cm.To,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			})
		}
		return ms, nil
	}
}

// Builds and returns a usecase-type function that gets all messages from the log.
// The usecase function will return an error only when  it fails to open the log
func FetchNBySender(logFilePath string) FetchNByUsername {
	return func(username string, start, length int) ([]Message, error) {
		var sc *bufio.Scanner
		{
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
			if err != nil {
				return make([]Message, 0), err
			}
			defer logFile.Close()
			sc = bufio.NewScanner(logFile)
		}

		// Unmarshal the log into chat.Message because that's the structure which the log is built from
		cm := chat.Message{}
		ms := make([]Message, 0)

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
			ms = append(ms, Message{
				ID:        curLine,
				From:      cm.From,
				To:        cm.To,
				Body:      cm.Body,
				Timestamp: cm.TimeStamp,
			})
		}
		return ms, nil
	}
}
