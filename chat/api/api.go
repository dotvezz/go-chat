package api

import (
	"github.com/dotvezz/gochat/chat"
	"github.com/gorilla/mux"
	"net/http"
)

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

func (api *restAPI) ListenAndServe(addr string) {
	go http.ListenAndServe(addr, api.router)
}
