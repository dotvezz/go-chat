package user

import "fmt"

// Path is the format pattern for a specific user's Rest resource path
const Path = "/user/%s"

// GetPath takes a userName as a param and returns the Rest resource path for the given user
func GetPath(userName string) string {
	return fmt.Sprintf(Path, userName)
}
