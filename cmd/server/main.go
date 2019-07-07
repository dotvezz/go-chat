package main

import (
	"flag"
	"github.com/dotvezz/gochat/chat/connection"
	"github.com/dotvezz/gochat/chat/server/config"
	"github.com/dotvezz/gochat/chat/tracker"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Parse flags used
	flag.Parse()

	conf := config.Load()

	logFile, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", 0)
	l, err := net.Listen("tcp", conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	tr := tracker.New(timeStampProvider, logger)

	for {
		c, _ := l.Accept()
		go tr.Connect(connection.New(c))
	}
}

func timeStampProvider() int64 {
	return time.Now().Unix()
}
