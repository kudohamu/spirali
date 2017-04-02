package spirali

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("ReadConfig", func(t *testing.T) {
		t.Run("when valid content", func(t *testing.T) {
			content := `
        [dev]
        dsn = "foo"
        driver = "mysql"

        [prod]
        dsn = "bar"
        driver = "mysql"
      `

			c, err := ReadConfig(bytes.NewReader([]byte(content)))
			assert.NoError(t, err)
			assert.Equal(t, "foo", c.specificConfigs["dev"].Dsn)
			assert.Equal(t, "mysql", c.specificConfigs["dev"].Driver)
			assert.Equal(t, "bar", c.specificConfigs["prod"].Dsn)
			assert.Equal(t, "mysql", c.specificConfigs["prod"].Driver)
		})

		t.Run("when invalid content", func(t *testing.T) {
			content := `
dev:
  dsn: foo
  driver: mysql
prod:
  dsn: bar
  driver: mysql
`
			c, err := ReadConfig(bytes.NewReader([]byte(content)))
			assert.Error(t, err)
			assert.Nil(t, c)
		})
	})

	t.Run("WithEnv", func(t *testing.T) {
		cases := []struct {
			env         string
			expectEnv   string
			expectError error
		}{
			{
				env:         "dev",
				expectEnv:   "dev",
				expectError: nil,
			},
			{
				env:         "prod",
				expectEnv:   "",
				expectError: ErrEnvNotFound,
			},
		}

		for _, c := range cases {
			config := &Config{
				specificConfigs: map[string]*specificConfig{
					"dev": &specificConfig{},
				},
			}
			err := config.WithEnv(c.env)
			assert.Equal(t, c.expectEnv, config.Env)
			assert.Equal(t, c.expectError, err)
		}
	})

	t.Run("WithPath", func(t *testing.T) {
		cases := []struct {
			path string
		}{
			{path: "foo"},
		}

		for _, c := range cases {
			config := &Config{}
			config.WithPath(c.path)
			assert.Equal(t, c.path, config.Path)
		}
	})

	t.Run("Driver", func(t *testing.T) {
		config := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev":  &specificConfig{Driver: "foo"},
				"prod": &specificConfig{Driver: "bar"},
			},
		}

		cases := []struct {
			env    string
			expect string
		}{
			{
				env:    "dev",
				expect: "foo",
			},
			{
				env:    "prod",
				expect: "bar",
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			assert.Equal(t, c.expect, config.Driver())
		}
	})

	t.Run("Dsn", func(t *testing.T) {
		config := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev":  &specificConfig{Dsn: "hoge"},
				"prod": &specificConfig{Dsn: "huga"},
			},
		}

		cases := []struct {
			env    string
			expect string
		}{
			{
				env:    "dev",
				expect: "hoge",
			},
			{
				env:    "prod",
				expect: "huga",
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			assert.Equal(t, c.expect, config.Dsn())
		}
	})

	t.Run("Dir", func(t *testing.T) {
		config := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev":  &specificConfig{Dir: "hoge"},
				"prod": &specificConfig{Dir: "huga"},
			},
		}

		cases := []struct {
			env    string
			expect string
		}{
			{
				env:    "dev",
				expect: "hoge",
			},
			{
				env:    "prod",
				expect: "huga",
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			wd, _ := os.Getwd()
			assert.Equal(t, filepath.Join(wd, c.expect), config.Dir())
		}
	})
}
