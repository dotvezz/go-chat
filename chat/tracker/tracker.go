package tracker

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"log"
)

// New builds and returns an implementation of the chat.Tracker interface
func New(timeStampProvider func() int64, logger *log.Logger) chat.Tracker {
	tr := &tracker{
		logger:      logger,
		timeStamp:   timeStampProvider,
		messages:    make(chan chat.Message),
		connections: make([]chat.Connection, 0),
	}
	tr.start()
	return tr
}

// tracker is the private implementation of the public chat.Tracker interface
type tracker struct {
	logger      *log.Logger
	timeStamp   func() int64
	messages    chan chat.Message
	connections []chat.Connection
}

// Subscribes the connection to the message channel
func (t *tracker) Connect(conn chat.Connection) {
	t.connections = append(t.connections, conn)
	go func() {
		for {
			m, err := conn.Receive()
			if err != nil {
				conn.Close()
				return
			}
			m.TimeStamp = t.timeStamp()
			t.messages <- m
		}
	}()
}

// start begins the goroutine that listens to the message channel and sends
// messages to connections
func (t *tracker) start() {
	go func() {
		for {
			message := <-t.messages
			jsonm, _ := json.Marshal(message)
			t.logger.Println(string(jsonm))
			deadConnections := make([]int, 0)
			for i, conn := range t.connections {
				if err := conn.Send(message); err != nil {
					// If a Send returns an error, just super-aggressively flag it to be removed
					deadConnections = append(deadConnections, i)
				}
			}

			// For flagged dead connections, kill without remorse
			for _, i := range deadConnections {
				t.connections = append(t.connections[:i], t.connections[i+1:]...)
			}
		}
	}()
}
