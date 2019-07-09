package message

// Message is the basic business entity for a Message. Essentially a copy of chat.Message
type Message struct {
	ID        int
	From      string
	Body      string
	Timestamp int64
}
