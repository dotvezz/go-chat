package message

import "fmt"

// OfUserPath is the format pattern for a path to a paginated list of messages of a specific user
const OfUserPath = "/user/%s/message/?first=%d"

// ListPath is the format pattern for a path to a paginated list of messages
const ListPath = "/message/?first=%d"

// Path is the format pattern for a path to a specific message's Rest resource
const Path = "/message/%d"

// GetOfUserPath takes a username and a pagination parameter and returns a Rest path to a list of that user's
// message resources
func GetOfUserPath(userName string, first int) string {
	return fmt.Sprintf(OfUserPath, userName, first)
}

// GetPath takes an id parameter and returns a Rest path to that specific message resource
func GetPath(id int) string {
	return fmt.Sprintf(Path, id)
}

// GetListPath takes a pagination parameter and returns a Rest path to a list of message resources
func GetListPath(first int) string {
	return fmt.Sprintf(ListPath, first)
}
