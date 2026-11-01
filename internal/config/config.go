package config

import (
	"os"
)

// Config holds application configuration settings.
type Config struct {
	Port string
}

// Load reads configuration from the environment with default values.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port: port,
	}
}
