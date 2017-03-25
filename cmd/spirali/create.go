package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kudohamu/spirali"
	"github.com/urfave/cli"
)

// Create generate new migration files.
func Create(c *cli.Context) {
	defer func() {
		err := recover().(error)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else {
			fmt.Fprintf(os.Stdout, "generated migration files")
		}
	}()

	name := c.Args().First()
	if name == "" {
		panic(errors.New("base file name is missing"))
	}
	dir := c.GlobalString("path")
	if dir == "" {
		panic(errors.New("path is missing"))
	}
	dir = path.Clean(dir)

	if err := initializeMetaDataFileIfNotExist(dir); err != nil {
		panic(err)
	}

	m, err := spirali.Create(name, dir)
	if err != nil {
		panic(err)
	}
	metadata, err := readMetaData(dir)
	if err != nil {
		panic(err)
	}
	metadata.AddMigration(m.GetUpFileName(), m.GetDownFileName())
	if err := updateMetaData(metadata, dir); err != nil {
		panic(err)
	}
}

func readMetaData(dir string) (*spirali.MetaData, error) {
	p := filepath.Join(dir, spirali.MetaDataFileName)
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m, err := spirali.ReadMetaData(file)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func updateMetaData(m *spirali.MetaData, dir string) error {
	p := filepath.Join(dir, spirali.MetaDataFileName)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := m.Save(file); err != nil {
		return err
	}
	return nil
}
