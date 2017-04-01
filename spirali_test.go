package spirali

import (
	"testing"

	"github.com/kudohamu/spirali/internal/driver"
	"github.com/stretchr/testify/assert"
)

func TestUp(t *testing.T) {
	metadata := &MetaData{
		Migrations: Migrations{
			&Migration{Version: 20170201, Name: "foo"},
			&Migration{Version: 20181111, Name: "bar"},
			&Migration{Version: 20170301, Name: "baz"},
			&Migration{Version: 20180101, Name: "hoge"},
			&Migration{Version: 20170101, Name: "huga"},
		},
	}
	readable := &nopReadable{}
	config := &Config{
		specificConfigs: map[string]*specificConfig{
			"dev": &specificConfig{Dsn: "aaa", Driver: "bbb"},
		},
		Env: "dev",
	}

	t.Run("when schema table is not created", func(t *testing.T) {
		driver := &driver.TDriver{Created: false}

		err := Up(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.True(t, driver.JustCreated)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions)
		assert.Equal(t, 5, driver.CountOfExec)
	})

	t.Run("when remain migrations not applied", func(t *testing.T) {
		driver := &driver.TDriver{
			Created:  true,
			Versions: []uint64{20170101, 20170201, 20170301},
		}

		err := Up(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.False(t, driver.JustCreated)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions)
		assert.Equal(t, 2, driver.CountOfExec)
	})

	t.Run("when all migrations are already applied", func(t *testing.T) {
		driver := &driver.TDriver{
			Created:  true,
			Versions: []uint64{20170101, 20170201, 20170301, 20180101, 20181111},
		}

		err := Up(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.False(t, driver.JustCreated)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions)
		assert.Equal(t, 0, driver.CountOfExec)
	})
}

func TestDown(t *testing.T) {
	metadata := &MetaData{
		Migrations: Migrations{
			&Migration{Version: 20170201, Name: "foo"},
			&Migration{Version: 20181111, Name: "bar"},
			&Migration{Version: 20170301, Name: "baz"},
			&Migration{Version: 20180101, Name: "hoge"},
			&Migration{Version: 20170101, Name: "huga"},
		},
	}
	readable := &nopReadable{}
	config := &Config{
		specificConfigs: map[string]*specificConfig{
			"dev": &specificConfig{Dsn: "aaa", Driver: "bbb"},
		},
		Env: "dev",
	}

	t.Run("when schema table is not created", func(t *testing.T) {
		driver := &driver.TDriver{Created: false}

		err := Down(metadata, config, driver, readable)
		assert.Error(t, err)
		assert.Equal(t, ErrSchemaVersionIsZero, err)
		assert.True(t, driver.Created)
		assert.True(t, driver.JustCreated)
	})

	t.Run("when valid state", func(t *testing.T) {
		driver := &driver.TDriver{
			Created:  true,
			Versions: []uint64{20170101, 20170201, 20170301, 20180101, 20181111},
		}
		err := Down(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101}, driver.Versions)
		assert.Equal(t, 1, driver.CountOfExec)
	})
}
