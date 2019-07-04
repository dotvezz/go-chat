package chat

type Tracker interface {
	Connect(conn Connection)
	Start()
}

type Connection interface {
	Send(message Message)
	Receive() (Message, error)
}

type Message struct {
	From string
	To   string
	Body string
}
