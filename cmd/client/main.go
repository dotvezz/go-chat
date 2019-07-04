package main

import (
	"fmt"
	"github.com/dotvezz/gochat/chat"
	"github.com/dotvezz/gochat/chat/connection"
	"net"
	"os"
)

func main() {
	nc, err := net.Dial("tcp", "localhost:1026")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	conn := connection.New(nc)

	conn.Send(chat.Message{
		From: "me",
		To:   "everyone",
		Body: "Holy cow!",
	})
	for {
		m, err := conn.Receive()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			break
		}
		fmt.Println(m.Body)
	}
}
