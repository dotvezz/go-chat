package connection

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"net"
)

const maxMessageLength = 1 << 8

// New builds and returns a Connection around a net.Conn
func New(c net.Conn) chat.Connection {
	conn := new(connection)
	conn.socket = c
	return conn
}

type connection struct {
	userName string
	socket   net.Conn
}

func (c *connection) Send(m chat.Message) error {
	mjson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.socket.Write(mjson)
	return err
}

func (c *connection) Receive() (m chat.Message, err error) {
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

	return
}

func (c *connection) Close() {
	_ = c.socket.Close()
}
