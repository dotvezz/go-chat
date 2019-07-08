package config

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
