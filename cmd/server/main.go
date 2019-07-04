package main

import (
	"github.com/dotvezz/gochat/chat/connection"
	"github.com/dotvezz/gochat/chat/tracker"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", ":1026")
	tr := tracker.New()
	go tr.Start()
	for {
		c, _ := l.Accept()
		go tr.Connect(connection.New(c))
	}
}
