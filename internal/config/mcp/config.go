package mcp

import (
	"github.com/yvv4git/go-juggler-mcp/internal/config"
)

// Config is the root application configuration.
type Config struct {
	Log   config.Log   `toml:"log"`
	Serve config.Serve `toml:"serve"`
}

// Load reads the configuration from the TOML file at path.
func Load(path string, cfg *Config) error {
	return config.Load(path, cfg)
}
