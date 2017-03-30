package spirali

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	testTime, _ := time.Parse("2006-01-02 15:04:05", "2017-01-02 03:04:05")

	t.Run("NewMigration", func(t *testing.T) {

		cases := []struct {
			time          time.Time
			name          string
			expectVersion uint64
			expectName    string
		}{
			{
				time:          testTime,
				name:          "create_test_table",
				expectVersion: 20170102030405,
				expectName:    "create_test_table",
			},
		}
		for _, c := range cases {
			m, err := NewMigration(c.time, c.name)

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
}
