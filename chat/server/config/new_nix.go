// +build linux darwin freebsd netbsd openbsd

package config

import "os"

// New returns a Config with default values
func New() Config {
	return Config{
		Port: ":1026",
		LogFile: os.Stdout.Name(),
	}
}
