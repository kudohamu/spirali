package spirali

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaData(t *testing.T) {
	t.Run("ReadMetaData", func(t *testing.T) {
		t.Run("when valid content", func(t *testing.T) {
			content := `
      {
        "migrations": [
           { "version": 20170102000000, "name": "create_foo_table" },
           { "version": 20170101000000, "name": "create_bar_table" },
           { "version": 20180101000000, "name": "create_baz_table" }
        ]
      }
      `

			m, err := ReadMetaData(bytes.NewReader([]byte(content)))
			assert.NoError(t, err)
			assert.Equal(t, uint64(20170101000000), m.Migrations[0].Version)
			assert.Equal(t, "create_bar_table", m.Migrations[0].Name)
			assert.Equal(t, uint64(20170102000000), m.Migrations[1].Version)
			assert.Equal(t, "create_foo_table", m.Migrations[1].Name)
			assert.Equal(t, uint64(20180101000000), m.Migrations[2].Version)
			assert.Equal(t, "create_baz_table", m.Migrations[2].Name)
		})

		t.Run("when invalid content", func(t *testing.T) {
			content := `migrations = []`

			m, err := ReadMetaData(bytes.NewReader([]byte(content)))
			assert.Error(t, err)
			assert.Nil(t, m)
		})
	})

	t.Run("Save", func(t *testing.T) {
		m := &MetaData{
			Migrations: Migrations{
				&Migration{
					Version: 1,
					Name:    "foo",
				},
				&Migration{
					Version: 2,
					Name:    "bar",
				},
				&Migration{
					Version: 3,
					Name:    "baz",
				},
			},
		}

		expect := "{\n  \"migrations\": [\n    {\n      \"version\": 1,\n      \"name\": \"foo\"\n    },\n    {\n      \"version\": 2,\n      \"name\": \"bar\"\n    },\n    {\n      \"version\": 3,\n      \"name\": \"baz\"\n    }\n  ]\n}"

		var buf bytes.Buffer
		err := m.Save(&buf)
		assert.NoError(t, err)
		assert.Equal(t, expect, buf.String())
	})
}
