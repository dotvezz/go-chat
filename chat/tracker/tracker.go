package tracker

import (
	"fmt"
	"github.com/dotvezz/gochat/chat"
)

func New() chat.Tracker {
	return &tracker{
		messages:    make(chan chat.Message),
		connections: make([]chat.Connection, 0),
	}
}

type tracker struct {
	messages    chan chat.Message
	connections []chat.Connection
}

func (t *tracker) Connect(conn chat.Connection) {
	fmt.Println("Client Connected")
	t.connections = append(t.connections, conn)
	for {
		m, err := conn.Receive()
		if err != nil {
			fmt.Println("Client Disconnected")
			return
		}
		t.messages <- m
	}
}

func (t *tracker) Start() {
	for {
		message := <-t.messages
		for _, conn := range t.connections {
			conn.Send(message)
		}
	}
}
