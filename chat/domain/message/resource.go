package message

// Resource is the basic Rest structure representing a Message resource
type Resource struct {
	ID   int
	Data struct {
		From string
		Body string
	}
	Meta struct {
		TimeStamp int64
	}
	Hypermedia struct {
		Self   string
		Sender string
	}
}

// Resources is the basic Rest response representing a collection of users
type Resources struct {
	Data       []Resource
	Hypermedia struct {
		NextPage string
		PrevPage string
	}
}
