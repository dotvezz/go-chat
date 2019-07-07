package config

// New returns a Config with default values
func New() Config {
	return Config{
		Port: ":1026",
		// To be honest, I don't know if os.Stdout.Name works on Windows
		// So LogFile is not filled by default, to avoid cryptic validation errors
		// since "LogFile is required" is easier to understand out of context than
		// "Couldn't open /dev/stdout" is
	}
}
