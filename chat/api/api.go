package api

import (
	"github.com/dotvezz/go-chat/chat"
	"github.com/gorilla/mux"
	"net/http"
)

// New builds and returns a chat.RestAPI which serves the API for the chat app
func New(logFilePath string, tracker chat.Tracker) chat.RestAPI {
	router := mux.NewRouter()

	for _, endpoint := range initRoutes(logFilePath, tracker) {
		router.HandleFunc(endpoint.path, endpoint.handler).Methods(endpoint.method)
	}

	return &restAPI{
		router: router,
	}
}

type restAPI struct {
	router *mux.Router
}

// ListenAndServe serves the API for the chat app
func (api *restAPI) ListenAndServe(addr string) {
	go http.ListenAndServe(addr, api.router)
}
