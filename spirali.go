package spirali

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// Create generates new migration files.
func Create(vg VersionG, name string, config *Config, metadata *MetaData) (*MetaData, error) {
	m, err := NewMigration(vg, name)
	if err != nil {
		return nil, err
	}
	dir, err := config.Dir()
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
	io.WriteString(upWriter, "-- write SQL for applying this migration.")

	downfile, err := os.Create(filepath.Join(dir, m.GetDownFileName()))
	if err != nil {
		return nil, err
	}
	defer downfile.Close()
	downWriter := bufio.NewWriter(downfile)
	defer downWriter.Flush()
	io.WriteString(downWriter, "-- write SQL for rolling back this migration.")

	metadata.Migrations = append(metadata.Migrations, m)

	return metadata, nil
}

// Up applies migrations not applied.
func Up(metadata *MetaData, config *Config, driver Driver, readable Readable) error {
	dsn, err := config.Dsn()
	if err != nil {
		return err
	}
	if err := driver.Open(dsn); err != nil {
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

// Down rolls back the latest migration.
func Down(metadata *MetaData, config *Config, driver Driver, readable Readable) error {
	dsn, err := config.Dsn()
	if err != nil {
		return err
	}
	if err := driver.Open(dsn); err != nil {
		return err
	}
	defer driver.Close()

	if err := driver.Transaction(func() error {
		if err := driver.CreateVersionTableIfNotExists(); err != nil {
			return err
		}
		if err := metadata.Migrations.Down(driver, readable); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
