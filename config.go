package spirali

import (
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the spirali configuration.
type Config struct {
	specificConfigs map[string]*specificConfig
	Env             string
	Path            string
}

type specificConfig struct {
	Driver string `toml:"driver"`
	Dsn    string `toml:"dsn"`
	Dir    string `toml:"directory"`
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
func (c *Config) WithEnv(env string) *Config {
	c.Env = env
	return c
}

// WithPath sets path to config.
func (c *Config) WithPath(path string) *Config {
	c.Path = path
	return c
}

// Driver returns driver of current env.
func (c *Config) Driver() (string, error) {
	if sc, found := c.specificConfigs[c.Env]; found {
		return sc.Driver, nil
	}
	return "", ErrEnvNotFound
}

// Dsn returns dsn of current env.
func (c *Config) Dsn() (string, error) {
	if sc, found := c.specificConfigs[c.Env]; found {
		return sc.Dsn, nil
	}
	return "", ErrEnvNotFound
}

// Dir returns migration files directory.
func (c *Config) Dir() (string, error) {
	if sc, found := c.specificConfigs[c.Env]; found {
		currentDir, _ := os.Getwd()
		return filepath.Join(currentDir, sc.Dir), nil
	}
	return "", ErrEnvNotFound
}
