package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/kudohamu/spirali"
	"github.com/urfave/cli"
)

// Status shows migration status.
func Status(c *cli.Context) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.(error).Error())
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
	config.WithEnv(env)
	metadata, err := readMetaData(config)
	if err != nil {
		panic(err)
	}
	driver, err := spirali.NewDriver(config)
	if err != nil {
		panic(err)
	}

	// try to write migration status.
	if err := spirali.Status(metadata, config, driver, os.Stdout); err != nil {
		panic(err)
	}
}
