package transport

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat/api/log/users"
	"github.com/dotvezz/gochat/chat/api/transport/urls"
	"github.com/gorilla/mux"
	"net/http"
)

// FetchUser returns an http.HandlerFunc which searches the log for a specific user by user name
// If the user has never sent a message, then the user will not be found.
// The necessary business logic is injected as a dependency
func FetchUser(fetch users.FetchUser) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userName, ok := mux.Vars(request)["userName"]
		if !ok {
			http.NotFound(writer, request)
			return
		}
		u, err := fetch(userName)
		if err == users.NotFound {
			http.NotFound(writer, request)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonu, err := json.Marshal(buildUserResource(u.Name))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// FetchUser returns an http.HandlerFunc which returns all users which have ever sent a message
// The necessary business logic is injected as a dependency
func FetchAllUsers(fetch users.FetchAllUsers) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		us, err := fetch()
		if err == users.NotFound {
			http.NotFound(writer, request)
			return
		}
		urs := Users{}

		for _, u := range us {
			urs.Data = append(urs.Data, buildUserResource(u.Name))
		}

		jsonu, err := json.Marshal(urs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(jsonu)
	}
}

// buildUserResource takes a userName and builds the Rest Resource, including hypermedia
func buildUserResource(userName string) User {
	ur := User{}
	ur.Data.Name = userName
	ur.Hypermedia.Self = urls.GetUser(userName)
	ur.Hypermedia.Messages = urls.GetMessagesOfUser(userName, 0)
	return ur
}
