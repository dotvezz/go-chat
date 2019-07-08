package api

import (
	"github.com/dotvezz/gochat/chat/api/log/messages"
	"github.com/dotvezz/gochat/chat/api/log/users"
	"github.com/dotvezz/gochat/chat/api/transport"
	"net/http"
)

type route struct {
	path    string
	handler http.HandlerFunc
}

func initRoutes(logFilePath string) []route {
	return []route{
		{
			path:    "/user/{userName}/message",
			handler: transport.FetchMessagesOfUser(messages.FetchNBySender(logFilePath)),
		},
		{
			path:    "/user/{userName}",
			handler: transport.FetchUser(users.FetchOne(logFilePath)),
		},
		{
			path:    "/user/",
			handler: transport.FetchAllUsers(users.FetchAll(logFilePath)),
		},
		{
			path:    "/message/{messageID}",
			handler: transport.FetchMessage(messages.FetchOne(logFilePath)),
		},
		{
			path:    "/message/",
			handler: transport.FetchMessages(messages.FetchN(logFilePath)),
		},
	}
}
