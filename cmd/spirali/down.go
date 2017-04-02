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
	configPath := c.GlobalString("path")
	if configPath == "" {
		panic(errors.New("path is missing"))
	}
	configPath = path.Clean(configPath)
	env := c.GlobalString("env")
	if env == "" {
		panic(errors.New("env is missing"))
	}

	// setup.
	config, err := openConfig(configPath)
	if err != nil {
		panic(err)
	}
	if e := config.WithEnv(env); e != nil {
		panic(e)
	}
	metadata, err := readMetaData(config)
	if err != nil {
		panic(err)
	}
	driver, err := spirali.NewDriver(config)
	if err != nil {
		panic(err)
	}
	readable := spirali.NewReadableFromDir(config.Dir())

	// try to roll back migrations.
	if err := spirali.Down(metadata, config, driver, readable); err != nil {
		panic(err)
	}
}
