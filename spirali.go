package spirali

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Create generate new migration files.
func Create(name string, dir string) (*Migration, error) {
	t := time.Now()
	m := NewMigration(t, name)
	upfile, err := os.Create(filepath.Join(dir, m.GetUpFileName()))
	if err != nil {
		return nil, err
	}
	defer upfile.Close()
	upWriter := bufio.NewWriter(upfile)
	defer upWriter.Flush()
	WriteUpTemplate(upWriter)

	downfile, err := os.Create(filepath.Join(dir, m.GetDownFileName()))
	if err != nil {
		return nil, err
	}
	defer downfile.Close()
	downWriter := bufio.NewWriter(downfile)
	defer downWriter.Flush()
	WriteDownTemplate(downWriter)

	return m, nil
}

// WriteUpTemplate is ...
func WriteUpTemplate(w io.Writer) {
	io.WriteString(w, `
-- write SQL for applying this migration.
`)
}

// WriteDownTemplate is ...
func WriteDownTemplate(w io.Writer) {
	io.WriteString(w, `
-- write SQL for rolling back this migration.
`)
}
