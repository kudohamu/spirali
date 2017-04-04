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
			env    string
			expect string
		}{
			{
				env:    "dev",
				expect: "dev",
			},
			{
				env:    "prod",
				expect: "prod",
			},
		}

		for _, c := range cases {
			config := &Config{
				specificConfigs: map[string]*specificConfig{
					"dev": &specificConfig{},
				},
			}
			config.WithEnv(c.env)
			assert.Equal(t, c.expect, config.Env)
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
				"dev": &specificConfig{Driver: "foo"},
			},
		}

		cases := []struct {
			env          string
			expectDriver string
			expectError  error
		}{
			{
				env:          "dev",
				expectDriver: "foo",
				expectError:  nil,
			},
			{
				env:          "prod",
				expectDriver: "",
				expectError:  ErrEnvNotFound,
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			d, err := config.Driver()
			assert.Equal(t, c.expectDriver, d)
			assert.Equal(t, c.expectError, err)
		}
	})

	t.Run("Dsn", func(t *testing.T) {
		config := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev": &specificConfig{Dsn: "hoge"},
			},
		}

		cases := []struct {
			env         string
			expectDsn   string
			expectError error
		}{
			{
				env:         "dev",
				expectDsn:   "hoge",
				expectError: nil,
			},
			{
				env:         "prod",
				expectDsn:   "",
				expectError: ErrEnvNotFound,
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			d, err := config.Dsn()
			assert.Equal(t, c.expectDsn, d)
			assert.Equal(t, c.expectError, err)
		}
	})

	t.Run("Dir", func(t *testing.T) {
		config := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev": &specificConfig{Dir: "hoge"},
			},
		}
		wd, _ := os.Getwd()

		cases := []struct {
			env         string
			expectDir   string
			expectError error
		}{
			{
				env:         "dev",
				expectDir:   filepath.Join(wd, "hoge"),
				expectError: nil,
			},
			{
				env:         "prod",
				expectDir:   "",
				expectError: ErrEnvNotFound,
			},
		}

		for _, c := range cases {
			config.WithEnv(c.env)
			d, err := config.Dir()

			assert.Equal(t, c.expectDir, d)
			assert.Equal(t, c.expectError, err)
		}
	})
}
