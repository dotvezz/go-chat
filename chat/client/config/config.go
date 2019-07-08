package config

// New returns a new Config with default values
func New() Config {
	return Config{
		Host:     "localhost",
		Port:     ":1026",
		UserName: "Someone",
	}
}

// Config holds the configuration values for a chat client
type Config struct {
	// The address to connect to
	Host string `validate:"required"`
	// The port to connect on, must begin with ":"
	Port string `validate:"required,startswith=:"`
	// The username to display as the sender of messages
	UserName string `validate:"required"`
}
