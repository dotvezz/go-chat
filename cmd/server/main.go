package main

import (
	"github.com/dotvezz/go-chat/chat"
	"github.com/dotvezz/go-chat/chat/api"
	"github.com/dotvezz/go-chat/chat/connection"
	"github.com/dotvezz/go-chat/chat/domain/server/config"
	"github.com/dotvezz/go-chat/chat/tracker"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Load the configuration
	conf := config.New()
	chat.LoadConfig(&conf)

	// Prepare the logger. Attempt to create a log file if the specified file doesn't exist
	var logger *log.Logger
	logFile, err := os.OpenFile(conf.LogFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger = log.New(logFile, "", 0)

	// Prepare the listener for connections on the configured port
	listener, err := net.Listen("tcp", conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Load the tracker that handles messages
	tr := tracker.New(timeStampProvider, logger)

	if conf.APIServerPort != "" {
		apiServer := api.New(conf.LogFile, tr)
		apiServer.ListenAndServe(conf.APIServerPort)
	}

	// Listen for connections and send them to the tracker
	for {
		c, _ := listener.Accept()
		tr.Connect(connection.NewServer(c))
	}
}

func timeStampProvider() int64 {
	return time.Now().Unix()
}
