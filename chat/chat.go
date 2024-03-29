package chat

import (
	"encoding/json"
	"flag"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
)

var configPath = flag.String("conf", "", "Path to the json config file")

// Tracker is an abstraction for managing connections to the chat service.
// It must internally maintain all Connections and manage the sending and
// receiving of messages from and to the Connections
type Tracker interface {
	// Connect accepts a connection and begins tracking it, receiving and
	// sending messages to other connections on the channel
	Connect(conn Connection)
	// Connections returns the currently active connections on the tracker
	Connections() []Connection
	// Post accepts a message to send to the tracker
	Post(message Message)
}

// Connection is an abstraction for managing connections between client and
// server.
type Connection interface {
	Send(message Message) error
	Receive() (Message, error)
	UserName() string
	Close()
}

// RestAPI is a Simple interface for serving HTTP requests
type RestAPI interface {
	ListenAndServe(addr string)
}

// Message is a structure that holds the body and metadata of a message sent
// to or from a client connection
type Message struct {
	From      string
	Body      string
	TimeStamp int64
}

// LoadConfig loads the config file, validates its contents, and hydrates it into the conf value parameter.
// Calls log.Fatal for any failures
func LoadConfig(conf interface{}) {
	// Parse flags used
	flag.Parse()

	// Just use the default Config if there's no path provided
	if configPath == nil || *configPath == "" {
		return
	}

	// Open the config file
	f, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the file's contents into the Config
	err = json.Unmarshal(f, conf)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the Config
	err = validator.New().Struct(conf)
	if err != nil {
		log.Fatal(err)
	}
}
