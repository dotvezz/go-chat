package message

import "fmt"

const OfUserPath = "/user/%s/message/?first=%d"
const ListPath = "/message/?first=%d"
const Path = "/message/%d"

func GetOfUserPath(userName string, first int) string {
	return fmt.Sprintf(OfUserPath, userName, first)
}

func GetPath(id int) string {
	return fmt.Sprintf(Path, id)
}

func GetListPath(first int) string {
	return fmt.Sprintf(ListPath, first)
}
