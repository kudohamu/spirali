package spirali

import (
	"io"

	"github.com/BurntSushi/toml"
)

// ConfigFileName ...
const ConfigFileName = "config.toml"

// Config represents the spirali configuration.
type Config struct {
	specificConfigs map[string]*specificConfig
	Env             string
	Dir             string
}

type specificConfig struct {
	Driver string `toml:"driver"`
	Dsn    string `toml:"dsn"`
}

// ReadConfig reads the spirali configuration from io.Reader.
func ReadConfig(r io.Reader) (*Config, error) {
	var c Config
	if _, err := toml.DecodeReader(r, &c.specificConfigs); err != nil {
		return nil, err
	}
	return &c, nil
}

// WithEnv sets env to config.
func (c *Config) WithEnv(env string) error {
	for key := range c.specificConfigs {
		if key == env {
			c.Env = env
			return nil
		}
	}
	return ErrEnvNotFound
}

// WithDir sets dir to config.
func (c *Config) WithDir(dir string) {
	c.Dir = dir
}

// Driver returns driver of current env.
func (c *Config) Driver() string {
	return c.specificConfigs[c.Env].Driver
}

// Dsn returns dsn of current env.
func (c *Config) Dsn() string {
	return c.specificConfigs[c.Env].Dsn
}
