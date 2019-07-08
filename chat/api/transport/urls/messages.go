package urls

import "fmt"

const MessagesOfUser = "/user/%s/message/?first=%d"
const Messages = "/message/?first=%d"
const Message = "/message/%d"

func GetMessagesOfUser(userName string, first int) string {
	return fmt.Sprintf(MessagesOfUser, userName, first)
}

func GetMessage(id int) string {
	return fmt.Sprintf(Message, id)
}

func GetMessages(first int) string {
	return fmt.Sprintf(Messages, first)
}
