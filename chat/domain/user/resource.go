package user

// Resource is the basic Rest structure representing a User resource
type Resource struct {
	Data struct {
		Name string
	}
	Meta struct {
		Online bool
	}
	Hypermedia struct {
		Self     string
		Messages string
	}
}

// Resources is the basic Rest response representing a collection of Users
type Resources struct {
	Data []Resource
}
