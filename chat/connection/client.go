package connection

import (
	"encoding/json"
	"github.com/dotvezz/go-chat/chat"
	"net"
)

// NewClient builds and returns a Connection around a net.Conn
func NewClient(c net.Conn) chat.Connection {
	clientConn := new(clientConnection)
	clientConn.socket = c
	return clientConn
}

type clientConnection struct {
	connection
}

// Send pushes a message over the socket. Returns an error whenever anything goes wrong at all.
// Also handles commands, which differentiates this from the connection and serverConnection implementations
// of chat.Connection
func (c *clientConnection) Send(m chat.Message) error {
	mjson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.socket.Write(mjson)
	if err == nil {
		c.parseCommand(m.Body)
	}
	return err
}
