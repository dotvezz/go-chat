package chat

// Tracker is an abstraction for managing connections to the chat service.
// It must internally maintain all Connections and manage the sending and
// receiving of messages from and to the Connections
type Tracker interface {
	// Connect accepts a connection and begins tracking it, receiving and
	// sending messages to other connections on the channel
	Connect(conn Connection)
}

// Connection is an abstraction for managing connections between client and
// server.
type Connection interface {
	Send(message Message) error
	Receive() (Message, error)
	Close()
}

type Message struct {
	From      string
	To        string
	Body      string
	TimeStamp int64
}
