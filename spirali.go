package spirali

import (
	"bufio"
	"os"
	"path/filepath"
	"time"
)

// Create generates new migration files.
func Create(name string, dir string) (*Migration, error) {
	t := time.Now()
	m, err := NewMigration(t, name)
	if err != nil {
		return nil, err
	}

	upfile, err := os.Create(filepath.Join(dir, m.GetUpFileName()))
	if err != nil {
		return nil, err
	}
	defer upfile.Close()
	upWriter := bufio.NewWriter(upfile)
	defer upWriter.Flush()
	writeUpTemplate(upWriter)

	downfile, err := os.Create(filepath.Join(dir, m.GetDownFileName()))
	if err != nil {
		return nil, err
	}
	defer downfile.Close()
	downWriter := bufio.NewWriter(downfile)
	defer downWriter.Flush()
	writeDownTemplate(downWriter)

	return m, nil
}

// Up applies migrations.
func Up(metadata *MetaData, config *Config, driver Driver, readable Readable) error {

	if err := driver.Open(config.Dsn()); err != nil {
		return err
	}
	defer driver.Close()

	if err := driver.Transaction(func() error {
		if err := driver.CreateVersionTableIfNotExists(); err != nil {
			return err
		}

		version, err := driver.GetCurrentVersion()
		if err != nil {
			return err
		}
		if err := metadata.Migrations.Up(driver, version, readable); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
