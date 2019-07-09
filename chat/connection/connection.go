package connection

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat"
	"net"
	"strings"
)

// Max message length is 256 bytes
const maxMessageLength = 256

// connection is an implementation of the chat.Connection interface
type connection struct {
	userName string
	socket   net.Conn
}

// UserName returns the userName as defined by this connection's most recent /nick command
func (c *connection) UserName() string {
	return c.userName
}

// Close terminates this connection
func (c *connection) Close() {
	_ = c.socket.Close()
}

// Send pushes a message over the socket. Returns an error whenever anything goes wrong at all.
func (c *connection) Send(m chat.Message) error {
	mjson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = c.socket.Write(mjson)
	return err
}

// Receive returns a message when it comes in over the socket.
// Returns an error if there were any communication or decoding problems.
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

// parseCommand checks a message body for commands, such as changing name with /nick
func (c *connection) parseCommand(body string) {
	if len(body) > 0 && body[0] == '/' {
		ss := strings.Split(body[1:], " ")
		// Check for username change
		if len(ss) > 1 && ss[0] == "nick" {
			c.userName = ss[1]
		}
	}
}
