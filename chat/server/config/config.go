package config

import (
	"encoding/json"
	"flag"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
)

var configPath = flag.String("conf", "", "Path to the json config file")

// Config holds the configuration values for a chat server
type Config struct {
	// The port that the server should listen on, must begin with ":"
	Port string `validate:"required,startswith=:"`
	// The path to the file to use for logging output.
	LogFile string `validate:"required"`
}

// Load loads the config file, validates its contents, and returns a hydrated config
// and calls log.Fatal for any failures
func Load() Config {
	// Create a Config struct with default values
	conf := New()

	// Just use the default Config if there's no path provided
	if configPath == nil || *configPath == "" {
		return conf
	}

	// Open the config file
	f, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the file's contents into the Config
	err = json.Unmarshal(f, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the Config
	err = validator.New().Struct(conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
