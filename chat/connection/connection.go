package connection

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"net"
)

const maxMessageLength = 1 << 8

func New(c net.Conn) chat.Connection {
	conn := new(connection)
	conn.conn = c
	return conn
}

type connection struct {
	userName string
	conn     net.Conn
}

func (c *connection) Send(m chat.Message) {
	mjson, _ := json.Marshal(m)
	_, _ = c.conn.Write(mjson)
}

func (c *connection) Receive() (m chat.Message, err error) {
	bs := make([]byte, maxMessageLength)
	length, err := c.conn.Read(bs)

	if err != nil {
		return
	}

	err = json.Unmarshal(bs[:length], &m)
	if err != nil {
		return
	}
	return
}
