#go-chat

[![Go Report Card](https://goreportcard.com/badge/github.com/dotvezz/go-chat)](https://goreportcard.com/report/github.com/dotvezz/go-chat)

go-chat is a small example of a chat server and chat client, written in Go. 
It has a RESTful API that can optionally be started to provide another interface
for interacting with messages and users. 

## Server

The server can be run in a few ways

- `go install github.com/dotvezz/go-chat/cmd/server`
  - This installs a binary called `server` to your `$GOPATH/bin` directory
- Clone this repository and `go run go-chat/cmd/server/main.go`
- Clone this repository and build `go-chat/cmd/server/main.go`

Every new user joins with an empty username. the `/nick` command changes usernames.

### Configuration 

The server has a simple configuration with three options:

```go
// Config holds the configuration values for a chat server
type Config struct {
	// The port that the server should listen on, must begin with ":"
	Port string `validate:"required,startswith=:"`
	// The path to the file to use for logging output.
	// Must be readable (Not Stdout) if the APIServerPort is set
	LogFile string `validate:"required"`
	// The port the REST API server should listen on. The API is disabled if this is left empty
	// The LogFile must be readable (Not Stdout) for the API server to start
	// Must begin with ":"
	APIServerPort string `validate:"startswith=:"`
}
```

The server has a default configuration with the API disabled and the chat server
listening on `:1026`. On UNIX-like operating systems, it logs to `/dev/stdout`
by default.

To load a new config, run the server with the `--conf={conf_path}` flag. The conf must be in
json format.

Note: For Windows, there is no default log output, so the log must be configured.
(Playing it safe since I don't know if `os.STDOUT.Fd()` works on Windows. I'd guess
it does, but not totally sure)

## API

The RESTful API is enabled when the server's `APIServerPort` config value is set to a valid
port, and the `LogFile` config value is set to a readable file (not `/dev/stdout`) (because
the log backs much of the API functionality).

The API has the following endpoints available, with a few example requests and responses:

- `GET: {host}/user/{userName}/message/`
    - Gets a paginated list of messages sent by a user
- `GET: {host}/user/{userName}`
    - Gets a specific user
    - Includes whether the user is online
```json
# GET: /user/ben
{
  "Data": {
    "Name": "ben"
  },
  "Meta": {
    "Online": true
  },
  "Hypermedia": {
    "Self": "/user/ben",
    "Messages": "/user/ben/message/?first=0"
  }
}
```
   
- `DELETE: {host}/user/{userName}`
    - Kicks a user from the chat server
- `GET: {host}/user/`
    - Gets a list of all users who have been on the server
- `GET: {host}/message/{messageID}`
    - Gets a specific message
```json
# GET: /message/52
{
  "ID": 52,
  "Data": {
    "From": "ben",
    "Body": "lol"
  },
  "Meta": {
    "TimeStamp": 1562638972
  },
  "Hypermedia": {
    "Self": "/message/52",
    "Sender": "/user/ben"
  }
}

```
- `GET: {host}/message/`
    - Gets a paginated list of messages
- `POST: {host}/message/`
    - Posts a message to the server

```json
# POST: /message/
# Request body:
{
  "Data": {
    "From": "ben",
    "Body": "lol"
  },
}

```

## Client

The client, like the server, can be run in a few ways:

- `go install github.com/dotvezz/go-chat/cmd/client`
  - This installs a binary called `client` to your `$GOPATH/bin` directory
- Clone this repository and `go run go-chat/cmd/client/main.go`
- Clone this repository and build `go-chat/cmd/client/main.go`

The client is a thin TUI with just a message box and input box. Ctrl+C or Esc will close
the TUI.

### Config

The client has a simple config:

```go
// Config holds the configuration values for a chat client
type Config struct {
	// The address to connect to
	Host string `validate:"required"`
	// The port to connect on, must begin with ":"
	Port string `validate:"required,startswith=:"`
}

```

## Third-Party Code

This project uses third party code through Go Modules:

- github.com/gorilla/mux
    - For setting up API routing
- github.com/marcusolsson/tui-go
    - For the client TUI
- gopkg.in/go-playground/validator.v9
    - For validation of configuration values
    
## Issues

- As it stands, every user joins with an empty username. This breaks the API
endpoints which use usernames, since emptystring can't be used as a resource identifier.
- The empty username is also just silly in general, a handshake process when initializing 
a connection, which includes identification, would be good to implement.