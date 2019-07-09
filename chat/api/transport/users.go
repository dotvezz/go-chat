package transport

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat/domain/message"
	"github.com/dotvezz/gochat/chat/domain/user"
	"github.com/gorilla/mux"
	"net/http"
)

// GetUser returns an http.HandlerFunc which searches for a specific user by user name
// If the user has never sent a message, then the user will not be found.
// An implementation of the user.Fetch usecase is injected as a dependency
func GetUser(fetch user.Fetch) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userName, ok := mux.Vars(request)["userName"]
		if !ok {
			http.NotFound(writer, request)
			return
		}
		u, err := fetch(userName)
		if err == user.NotFound {
			http.NotFound(writer, request)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonu, err := json.Marshal(buildUserResource(u))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// GetAllUsers returns an http.HandlerFunc which returns all users who have ever sent a message
// Pagination is not supported
// An implementation of the user.FetchAll usecase is injected as a dependency
func GetAllUsers(fetch user.FetchAll) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		us, err := fetch()
		if err == user.NotFound {
			http.NotFound(writer, request)
			return
		}
		urs := user.Resources{}

		for _, u := range us {
			urs.Data = append(urs.Data, buildUserResource(u))
		}

		jsonu, err := json.Marshal(urs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(jsonu)
	}
}

// DeleteUser returns an http.HandlerFunc which kicks the specified user
// An implementation of user.Kick is injected as a dependency
func DeleteUser(kick user.Kick) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userName, ok := mux.Vars(request)["userName"]
		if !ok {
			http.NotFound(writer, request)
			return
		}
		ue, err := kick(userName)
		if err == user.NotFound {
			http.NotFound(writer, request)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonu, err := json.Marshal(ue)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(jsonu)
	}
}

// buildUserResource takes a userName and builds the Rest Resource, including hypermedia
func buildUserResource(u user.User) user.Resource {
	ur := user.Resource{}
	ur.Data.Name = u.Name
	ur.Meta.Online = u.Online
	ur.Hypermedia.Self = user.GetPath(u.Name)
	ur.Hypermedia.Messages = message.GetOfUserPath(u.Name, 0)
	return ur
}
