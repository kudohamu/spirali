package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/kudohamu/spirali"
	"github.com/urfave/cli"
)

// Down rolls back the latest migration.
func Down(c *cli.Context) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.(error).Error())
		} else {
			fmt.Fprintf(os.Stdout, "rolled back the latest migration\n")
		}
	}()

	// check options.
	dir := c.GlobalString("path")
	if dir == "" {
		panic(errors.New("path is missing"))
	}
	dir = path.Clean(dir)
	env := c.GlobalString("env")
	if dir == "" {
		panic(errors.New("env is missing"))
	}

	// setup.
	config, err := openConfig(dir)
	if err != nil {
		panic(err)
	}
	metadata, err := readMetaData(dir)
	if err != nil {
		panic(err)
	}
	if e := config.WithEnv(env); e != nil {
		panic(e)
	}
	driver, err := spirali.NewDriver(config)
	if err != nil {
		panic(err)
	}
	readable := spirali.NewReadableFromDir(dir)

	// try to roll back migrations.
	if err := spirali.Down(metadata, config, driver, readable); err != nil {
		panic(err)
	}
}
