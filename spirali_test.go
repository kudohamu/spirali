package spirali

import (
	"bytes"
	"fmt"
	"testing"
	"time"

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
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions())
		assert.Equal(t, 5, driver.CountOfExec)
	})

	t.Run("when remain migrations not applied", func(t *testing.T) {
		driver := &driver.TDriver{
			Created: true,
			Rows: []*driver.Row{
				&driver.Row{Version: 20170101},
				&driver.Row{Version: 20170201},
				&driver.Row{Version: 20170301},
			},
		}

		err := Up(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.False(t, driver.JustCreated)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions())
		assert.Equal(t, 2, driver.CountOfExec)
	})

	t.Run("when all migrations are already applied", func(t *testing.T) {
		driver := &driver.TDriver{
			Created: true,
			Rows: []*driver.Row{
				&driver.Row{Version: 20170101},
				&driver.Row{Version: 20170201},
				&driver.Row{Version: 20170301},
				&driver.Row{Version: 20180101},
				&driver.Row{Version: 20181111},
			},
		}

		err := Up(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.False(t, driver.JustCreated)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101, 20181111}, driver.Versions())
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
			Created: true,
			Rows: []*driver.Row{
				&driver.Row{Version: 20170101},
				&driver.Row{Version: 20170201},
				&driver.Row{Version: 20170301},
				&driver.Row{Version: 20180101},
				&driver.Row{Version: 20181111},
			},
		}
		err := Down(metadata, config, driver, readable)
		assert.NoError(t, err)
		assert.Equal(t, []uint64{20170101, 20170201, 20170301, 20180101}, driver.Versions())
		assert.Equal(t, 1, driver.CountOfExec)
	})
}

func TestStatus(t *testing.T) {
	metadata := &MetaData{
		Migrations: Migrations{
			&Migration{Version: 20170201, Name: "foo"},
			&Migration{Version: 20181111, Name: "bar"},
			&Migration{Version: 20170301, Name: "baz"},
		},
	}
	config := &Config{
		specificConfigs: map[string]*specificConfig{
			"dev": &specificConfig{Dsn: "aaa", Driver: "bbb"},
		},
		Env: "dev",
	}
	header := "migration status of `dev` environment"

	t.Run("when schema table is not created", func(t *testing.T) {
		driver := &driver.TDriver{Created: false}
		var buff bytes.Buffer

		err := Status(metadata, config, driver, &buff)
		assert.NoError(t, err)
		assert.True(t, driver.Created)
		assert.True(t, driver.JustCreated)
		assert.Equal(t, fmt.Sprintf(`%s
%s
 not applied         | 20170201_foo
 not applied         | 20170301_baz
 not applied         | 20181111_bar
%s
`, header, separator, separator), buff.String())
	})

	t.Run("when remain migrations not applied", func(t *testing.T) {
		time1, _ := time.Parse("2006-01-02 15:04:05", "2018-12-01 11:22:33")
		time2, _ := time.Parse("2006-01-02 15:04:05", "2018-12-02 11:22:33")
		driver := &driver.TDriver{
			Created: true,
			Rows: []*driver.Row{
				&driver.Row{Version: 20170201, CreatedAt: time1},
				&driver.Row{Version: 20170301, CreatedAt: time2},
			},
		}
		var buff bytes.Buffer

		err := Status(metadata, config, driver, &buff)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`%s
%s
 2018-12-01 11:22:33 | 20170201_foo
 2018-12-02 11:22:33 | 20170301_baz
 not applied         | 20181111_bar
%s
`, header, separator, separator), buff.String())
	})

	t.Run("when all migrations are applied", func(t *testing.T) {
		time1, _ := time.Parse("2006-01-02 15:04:05", "2018-12-01 11:22:33")
		time2, _ := time.Parse("2006-01-02 15:04:05", "2018-12-02 11:22:33")
		time3, _ := time.Parse("2006-01-02 15:04:05", "2018-12-03 11:22:33")
		driver := &driver.TDriver{
			Created: true,
			Rows: []*driver.Row{
				&driver.Row{Version: 20170201, CreatedAt: time1},
				&driver.Row{Version: 20170301, CreatedAt: time2},
				&driver.Row{Version: 20181111, CreatedAt: time3},
			},
		}
		var buff bytes.Buffer

		err := Status(metadata, config, driver, &buff)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`%s
%s
 2018-12-01 11:22:33 | 20170201_foo
 2018-12-02 11:22:33 | 20170301_baz
 2018-12-03 11:22:33 | 20181111_bar
%s
`, header, separator, separator), buff.String())
	})
}
