package spirali

import (
	"reflect"
	"testing"

	"github.com/kudohamu/spirali/internal/driver"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	t.Run("NewDriver", func(t *testing.T) {
		c := &Config{
			specificConfigs: map[string]*specificConfig{
				"dev":  &specificConfig{Driver: "mysql"},
				"prod": &specificConfig{Driver: "invalid"},
			},
		}

		t.Run("when driver name is valid", func(t *testing.T) {
			c.WithEnv("dev")
			d, err := NewDriver(c)
			assert.NoError(t, err)
			assert.Equal(t, reflect.TypeOf(&driver.Mysql{}), reflect.TypeOf(d))
		})

		t.Run("when driver name is invalid", func(t *testing.T) {
			c.WithEnv("prod")
			d, err := NewDriver(c)
			assert.Equal(t, ErrUnknownDriver, err)
			assert.Nil(t, d)
		})
	})
}
