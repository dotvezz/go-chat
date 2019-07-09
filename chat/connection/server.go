package connection

import (
	"encoding/json"
	"github.com/dotvezz/go-chat/chat"
	"net"
)

// NewServer builds and returns a Connection around a net.Conn
func NewServer(c net.Conn) chat.Connection {
	clientConn := new(serverConnection)
	clientConn.socket = c
	return clientConn
}

type serverConnection struct {
	connection
}

// Receive returns a message when it comes in over the socket.
// Also handles commands, which differentiates this from the connection and clientConnection
// implementations of chat.Connection
// Returns an error if there were any communication or decoding problems.
func (c *serverConnection) Receive() (m chat.Message, err error) {
read:
	bs := make([]byte, maxMessageLength)
	length, err := c.socket.Read(bs)

	if err != nil {
		return
	}

	err = json.Unmarshal(bs[:length], &m)
	if err != nil {
		return
	}

	// Ignore messages with empty body
	if m.Body == "" {
		goto read
	}
	c.parseCommand(m.Body)
	return
}
