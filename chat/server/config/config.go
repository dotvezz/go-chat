package config

// Config holds the configuration values for a chat server
type Config struct {
	// The port that the server should listen on, must begin with ":"
	Port string `validate:"required,startswith=:"`
	// The path to the file to use for logging output.
	LogFile string `validate:"required"`
}
