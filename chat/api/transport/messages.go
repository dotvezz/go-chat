package transport

import (
	"encoding/json"
	"github.com/dotvezz/go-chat/chat"
	"github.com/dotvezz/go-chat/chat/domain/message"
	"github.com/dotvezz/go-chat/chat/domain/user"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"strconv"
)

// GetMessage returns an http.HandlerFunc which searches for a specific message
// An implementation of the message.Fetch usecase is injected as a dependency
func GetMessage(fetch message.Fetch) http.HandlerFunc {
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
		if err == message.ErrNotFound {
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

// GetNMessages returns an http.HandlerFunc which returns a list of messages
// Pagination is supported with the 'first' query
// An implementation of the message.FetchN usecase is injected as a dependency
func GetNMessages(fetch message.FetchN) http.HandlerFunc {
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

		mrs := message.Resources{}
		for _, m := range ms {
			mrs.Data = append(mrs.Data, buildMessageResource(m))
		}
		mrs.Hypermedia.NextPage = message.GetListPath(ifirst + 10)
		mrs.Hypermedia.PrevPage = message.GetListPath(int(math.Max(0, float64(ifirst-10))))
		jsonu, err := json.Marshal(mrs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// PostMessage returns an http.HandlerFunc which accepts a message resource in json format and posts
// that message to the tracker
// An implementation of the message.Post usecase is injected as a dependency
func PostMessage(post message.Post) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		dec := json.NewDecoder(request.Body)
		mr := message.Resource{}
		err := dec.Decode(&mr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		me := chat.Message{
			Body: mr.Data.Body,
			From: mr.Data.From,
		}
		writer.WriteHeader(http.StatusAccepted)
		post(me)
	}
}

// GetMessagesOfUser returns an http.HandlerFunc which searches messages from a specific user by user name
// An implementation of the message.FetchNByString usecase is injected as a dependency
func GetMessagesOfUser(fetch message.FetchNByString) http.HandlerFunc {
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

		mrs := message.Resources{}
		for _, m := range ms {
			mrs.Data = append(mrs.Data, buildMessageResource(m))
		}
		mrs.Hypermedia.NextPage = message.GetOfUserPath(userName, ifirst+10)
		mrs.Hypermedia.PrevPage = message.GetOfUserPath(userName, int(math.Max(0, float64(ifirst-10))))
		jsonu, err := json.Marshal(mrs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = writer.Write(jsonu)
	}
}

// buildMessageResource takes a message and builds the Rest Resource, including hypermedia
func buildMessageResource(m message.Message) message.Resource {
	mr := message.Resource{}
	mr.ID = m.ID
	mr.Data.From = m.From
	mr.Data.Body = m.Body
	mr.Meta.TimeStamp = m.Timestamp
	mr.Hypermedia.Self = message.GetPath(m.ID)
	mr.Hypermedia.Sender = user.GetPath(m.From)
	return mr
}
