package spirali

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ext = "sql"

// Migration is ...
type Migration struct {
	Version uint64 `json:"version"`
	Name    string `json:"name"`
}

// Migrations is array of migration.
type Migrations []*Migration

// NewMigration initialize migration struct.
func NewMigration(t time.Time, name string) (*Migration, error) {
	v, err := strconv.ParseUint(t.Format("20060102150405"), 10, 64)
	if err != nil {
		return nil, err
	}
	return &Migration{
		Version: v,
		Name:    name,
	}, nil
}

// GetUpFileName generate file name for applied migration.
func (m *Migration) GetUpFileName() string {
	return m.getFileName("up")
}

// GetDownFileName generate file name for rolled back migration.
func (m *Migration) GetDownFileName() string {
	return m.getFileName("down")
}

func (m *Migration) getFileName(suffix string) string {
	return strings.Join([]string{
		strconv.FormatUint(m.Version, 10),
		m.Name,
		strings.Join([]string{suffix, ext}, "."),
	}, "_")
}

// Up applies migrations not applied.
func (ms Migrations) Up(driver Driver, currentVersion uint64, readable Readable) error {
	ms.sort()

	for _, m := range ms {
		if currentVersion == 0 || currentVersion < m.Version {
			b, err := readable.Read(m.GetUpFileName())
			if err != nil {
				return err
			}
			if err := driver.Exec(string(b)); err != nil {
				return err
			}
			if err := driver.SetVersion(m.Version); err != nil {
				return err
			}
		}
	}
	return nil
}

// Down rolls back the latest migration.
func (ms Migrations) Down(driver Driver, readable Readable) error {
	ms.sort()

	if len(ms) == 0 {
		return ErrMigrationsNotExist
	}

	version, err := driver.GetCurrentVersion()
	if err != nil {
		return err
	}
	if version == 0 {
		return ErrSchemaVersionIsZero
	}

	for _, m := range ms {
		if version == m.Version {
			b, err := readable.Read(m.GetDownFileName())
			if err != nil {
				return err
			}
			if err := driver.Exec(string(b)); err != nil {
				return err
			}
			if err := driver.DeleteVersion(version); err != nil {
				return err
			}

			return nil
		}
	}

	return ErrMigrationFileNotFound
}

func (ms Migrations) sort() {
	sort.SliceStable(ms, func(i int, j int) bool {
		return ms[i].Version < ms[j].Version
	})
}

func writeUpTemplate(w io.Writer) {
	io.WriteString(w, `
-- write SQL for applying this migration.
`)
}

func writeDownTemplate(w io.Writer) {
	io.WriteString(w, `
-- write SQL for rolling back this migration.
`)
}
