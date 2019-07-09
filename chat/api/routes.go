package api

import (
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/api/log"
	"github.com/dotvezz/gochat/chat/api/tracker"
	"github.com/dotvezz/gochat/chat/api/transport"
	"net/http"
)

type route struct {
	path    string
	handler http.HandlerFunc
	method  string
}

// InitRoutes builds the routes used to define the APIs, as well as building http.HandlerFuncs for each route.
// It injects the dependencies as needed for each handlerFunc and usecase
func initRoutes(logFilePath string, tr chat.Tracker) []route {
	return []route{
		{
			path:    "/user/{userName}/message/",
			handler: transport.GetMessagesOfUser(log.FetchMessagesOfSender(logFilePath)),
			method:  http.MethodGet,
		},
		{
			path:    "/user/{userName}",
			handler: transport.GetUser(log.FetchUser(logFilePath, tracker.IsUserOnline(tr))),
			method:  http.MethodGet,
		},
		{
			path:    "/user/{userName}",
			handler: transport.DeleteUser(tracker.KickUser(tr)),
			method:  http.MethodDelete,
		},
		{
			path:    "/user/",
			handler: transport.GetAllUsers(log.FetchAllUsers(logFilePath)),
			method:  http.MethodGet,
		},
		{
			path:    "/message/{messageID}",
			handler: transport.GetMessage(log.FetchMessage(logFilePath)),
			method:  http.MethodGet,
		},
		{
			path:    "/message/",
			handler: transport.GetNMessages(log.FetchNMessages(logFilePath)),
			method:  http.MethodGet,
		},
		{
			path:    "/message/",
			handler: transport.PostMessage(tracker.PostMessage(tr)),
			method:  http.MethodPost,
		},
	}
}
