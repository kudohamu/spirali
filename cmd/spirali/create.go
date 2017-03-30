package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/kudohamu/spirali"
	"github.com/urfave/cli"
)

// Create generate new migration files.
func Create(c *cli.Context) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.(error).Error())
		} else {
			fmt.Fprintf(os.Stdout, "generated migration files\n")
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
	metadata.Migrations = append(metadata.Migrations, m)
	if err := updateMetaData(metadata, dir); err != nil {
		panic(err)
	}
}
