package spirali

import (
	"testing"

	"github.com/kudohamu/spirali/internal/driver"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	t.Run("NewMigration", func(t *testing.T) {

		cases := []struct {
			generator     VersionG
			name          string
			expectVersion uint64
			expectName    string
		}{
			{
				generator:     &IncrementalVersionG{CurrentVersion: 1},
				name:          "create_test_table",
				expectVersion: 2,
				expectName:    "create_test_table",
			},
		}
		for _, c := range cases {
			m, err := NewMigration(c.generator, c.name)

			assert.Nil(t, err)
			assert.Equal(t, c.expectVersion, m.Version)
			assert.Equal(t, c.expectName, m.Name)
		}
	})

	t.Run("GetUpFileName", func(t *testing.T) {
		cases := []struct {
			version uint64
			name    string
			expect  string
		}{
			{
				version: 20170101010101,
				name:    "create_test_table",
				expect:  "20170101010101_create_test_table_up.sql",
			},
		}
		for _, c := range cases {
			m := &Migration{
				Version: c.version,
				Name:    c.name,
			}

			assert.Equal(t, c.expect, m.GetUpFileName())
		}
	})

	t.Run("GetDownFileName", func(t *testing.T) {
		cases := []struct {
			version uint64
			name    string
			expect  string
		}{
			{
				version: 20170101010101,
				name:    "create_test_table",
				expect:  "20170101010101_create_test_table_down.sql",
			},
		}
		for _, c := range cases {
			m := &Migration{
				Version: c.version,
				Name:    c.name,
			}

			assert.Equal(t, c.expect, m.GetDownFileName())
		}
	})

	t.Run("Up", func(t *testing.T) {
		ms := Migrations{
			&Migration{Version: 2, Name: "foo"},
			&Migration{Version: 1, Name: "bar"},
			&Migration{Version: 3, Name: "baz"},
		}
		readable := &nopReadable{}

		t.Run("when current version is 0", func(t *testing.T) {
			driver := &driver.TDriver{
				Created:  true,
				Versions: []uint64{},
			}
			version := uint64(0)
			err := ms.Up(driver, version, readable)
			assert.NoError(t, err)
			assert.Equal(t, []uint64{1, 2, 3}, driver.Versions)
			assert.Equal(t, 3, driver.CountOfExec)
		})

		t.Run("when current version is not 0", func(t *testing.T) {
			driver := &driver.TDriver{
				Created:  true,
				Versions: []uint64{1, 2},
			}
			version := uint64(2)
			err := ms.Up(driver, version, readable)
			assert.NoError(t, err)
			assert.Equal(t, []uint64{1, 2, 3}, driver.Versions)
			assert.Equal(t, 1, driver.CountOfExec)
		})
	})

	t.Run("Down", func(t *testing.T) {
		ms := Migrations{
			&Migration{Version: 2, Name: "foo"},
			&Migration{Version: 1, Name: "bar"},
			&Migration{Version: 3, Name: "baz"},
		}
		readable := &nopReadable{}

		t.Run("when current version is 0", func(t *testing.T) {
			driver := &driver.TDriver{
				Created:  true,
				Versions: []uint64{},
			}
			err := ms.Down(driver, readable)
			assert.Error(t, err)
			assert.Equal(t, ErrSchemaVersionIsZero, err)
		})

		t.Run("when migration files do not exist", func(t *testing.T) {
			ms := Migrations{}

			driver := &driver.TDriver{
				Created:  true,
				Versions: []uint64{},
			}
			err := ms.Down(driver, readable)
			assert.Error(t, err)
			assert.Equal(t, ErrMigrationsNotExist, err)
		})

		t.Run("when valid state", func(t *testing.T) {
			driver := &driver.TDriver{
				Created:  true,
				Versions: []uint64{1, 2},
			}
			err := ms.Down(driver, readable)
			assert.NoError(t, err)
			assert.Equal(t, []uint64{1}, driver.Versions)
			assert.Equal(t, 1, driver.CountOfExec)
		})
	})
}

type nopReadable struct{}

func (nopReadable) Read(path string) ([]byte, error) {
	return []byte("foo"), nil
}
