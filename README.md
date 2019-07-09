# go-chat

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
(Playing it safe since I don't know if `os.Stdout.Fd()` works on Windows. I'd guess
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

The default is to connect to `localhost:1026`.

## Approach

### Code Structure

The API is built loosely following 
"[Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)" 
principles. Usecase types are defined in the `domain` package, and each type is 
implemented by a package which is named after the core implementation dependency.
For example, the `messages.Post` usecase is implemented by a function which uses
the chat tracker, and is therefore in a package named `tracker`.

The Chat Server is built primarily around two interfaces: `chat.Tracker` and `chat.Connection`.
The `chat.Tracker` is implemented with a single `chan` that all messages enter into. Every
Connection to the server can write messages to that channel, and the tracker distributes all
messages from the channel to every live connection.

The Chat Client is built primarily around the `chat.Connection` interface. It is the same 
`chat.Connection` that the Server uses, but has a slightly different implementation for
handling the `/nick name` command.

### Dependency Injection

The API uses dependency injection patterns for decoupling from the rest of the app. Usecases
are defined for business logic operations, and each usecase is implemented through either
the `chat.Tracker` or the log file. Implementations are injected where needed, including for
http handlers. If we wanted to remove the logfile requirement from some API operations, only
the implementation injected into the http handler needs to be updated.

## Third-Party Code

This project uses third party code through Go Modules:

- github.com/gorilla/mux
    - For setting up API routing
- github.com/marcusolsson/tui-go
    - For the client TUI
- gopkg.in/go-playground/validator.v9
    - For validation of configuration values
    
## Issues

- Connection-related Issues
  - As it stands, every user joins with an empty username. A handshake process when 
  initializing a connection, which includes identification, would be good to implement.
- Design Issues
  - Messages are serialized to json for transmission, and deserialized on the receiving end.
  This makes it hard to send messages using a simple telnet client. json can be manually
  written and sent, but it's unintuitive.
- API-Related Issues
  - The API uses pseudo-identifiers for `User` and `Message` resources. For users, it's
  the string username. For messages, it's the line which the message is on in the log file.
  - Related to the above: The "Post Message" API responds with Status Code `202 Accepted`,
  this is because of the difficulty guaranteeing that the ID (which should be in a `200 OK`
  response) would correctly match the line in the logfile the message lands on.
  - Also related to the above: Emptystring is a "valid" username, but it can't be used
  as a resource identifier.
- General Issues
  - There's no testing at all. I'll want to add testing, but for now see one of my other projects, [lime](http://github.com/dotvezz/lime/) for testing practice examples.
