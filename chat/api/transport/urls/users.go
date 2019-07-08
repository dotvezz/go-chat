package urls

import "fmt"

const User = "/user/%s"
const Users = "/user/"

func GetUser(userName string) string {
	return fmt.Sprintf(User, userName)
}