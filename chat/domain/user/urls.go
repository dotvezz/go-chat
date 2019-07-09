package user

import "fmt"

const Path = "/user/%s"
const ListPath = "/user/"

func GetPath(userName string) string {
	return fmt.Sprintf(Path, userName)
}
