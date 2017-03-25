package spirali

import (
	"strings"
	"time"
)

const ext = "sql"

// Migration is ...
type Migration struct {
	Version string
	Name    string
}

// NewMigration initialize migration struct.
func NewMigration(t time.Time, name string) *Migration {
	return &Migration{
		Version: t.Format("20060102150405"),
		Name:    name,
	}
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
		m.Version,
		m.Name,
		strings.Join([]string{suffix, ext}, "."),
	}, "_")
}
