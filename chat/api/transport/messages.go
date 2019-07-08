package transport

import (
	"encoding/json"
	"github.com/dotvezz/gochat/chat/api/log/messages"
	"github.com/dotvezz/gochat/chat/api/transport/urls"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"strconv"
)

// FetchUser returns an http.HandlerFunc which searches the log for a specific user by user name
// If the user has never sent a message, then the user will not be found.
// The necessary business logic is injected as a dependency
func FetchMessage(fetch messages.FetchMessage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		strID, ok := mux.Vars(request)["messageID"]
		if !ok {
			http.NotFound(writer, request)
			return
		}
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.NotFound(writer, request)
			return
		}
		m, err := fetch(id)
		if err == messages.NotFound {
			http.NotFound(writer, request)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonu, err := json.Marshal(buildMessageResource(m))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// FetchUser returns an http.HandlerFunc which searches the log for a specific user by user name
// If the user has never sent a message, then the user will not be found.
// The necessary business logic is injected as a dependency
func FetchMessages(fetch messages.FetchNMessages) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		first := request.FormValue("first")
		if first == "" {
			first = "0"
		}
		ifirst, err := strconv.Atoi(first)
		if err != nil {
			http.Error(writer, "'first' must be a positive integer", http.StatusBadRequest)
			return
		}
		ms, err := fetch(ifirst, 10)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		mrs := Messages{}
		for _, m := range ms {
			mrs.Data = append(mrs.Data, buildMessageResource(m))
		}
		mrs.Hypermedia.NextPage = urls.GetMessages(ifirst + 10)
		mrs.Hypermedia.PrevPage = urls.GetMessages(int(math.Max(0, float64(ifirst-10))))
		jsonu, err := json.Marshal(mrs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}


// FetchUser returns an http.HandlerFunc which searches the log for a specific user by user name
// If the user has never sent a message, then the user will not be found.
// The necessary business logic is injected as a dependency
func FetchMessagesOfUser(fetch messages.FetchNByUsername) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userName, ok := mux.Vars(request)["userName"]
		if !ok {
			http.NotFound(writer, request)
			return
		}
		first := request.FormValue("first")
		if first == "" {
			first = "0"
		}
		ifirst, err := strconv.Atoi(first)
		if err != nil {
			http.Error(writer, "'first' must be a positive integer", http.StatusBadRequest)
			return
		}
		ms, err := fetch(userName, ifirst, 10)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		mrs := Messages{}
		for _, m := range ms {
			mrs.Data = append(mrs.Data, buildMessageResource(m))
		}
		mrs.Hypermedia.NextPage = urls.GetMessagesOfUser(userName, ifirst + 10)
		mrs.Hypermedia.PrevPage = urls.GetMessagesOfUser(userName, int(math.Max(0, float64(ifirst-10))))
		jsonu, err := json.Marshal(mrs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// buildUserResource takes a message and builds the Rest Resource, including hypermedia
func buildMessageResource(m messages.Message) Message {
	mr := Message{}
	mr.ID = m.ID
	mr.Data.From = m.From
	mr.Data.To = m.To
	mr.Hypermedia.Self = urls.GetMessage(m.ID)
	mr.Hypermedia.Recipient = urls.GetUser(m.To)
	mr.Hypermedia.Sender = urls.GetUser(m.From)
	return mr
}
