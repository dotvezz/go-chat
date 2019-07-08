package transport

// The basic Rest structure representing a Message resource
type Message struct {
	ID   int
	Data struct {
		From string
		To   string
		Body string
	}
	Hypermedia struct {
		Self      string
		Sender    string
		Recipient string
	}
}

// The basic Rest structure representing a User resource
type User struct {
	Data struct {
		Name string
	}
	Hypermedia struct {
		Self     string
		Messages string
	}
}

// The basic Rest response representing a collection of Users
type Users struct {
	Data []User
}

// The basic Rest response representing a collection of users
type Messages struct {
	Data       []Message
	Hypermedia struct {
		NextPage string
		PrevPage string
	}
}
