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

	// check options.
	name := c.Args().First()
	if name == "" {
		panic(errors.New("base file name is missing"))
	}
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
	if err := initializeMetaDataFileIfNotExist(config); err != nil {
		panic(err)
	}
	metadata, err := readMetaData(config)
	if err != nil {
		panic(err)
	}
	vg := &spirali.TimestampBasedVersionG{}

	// try to create migration files.
	newMetaData, err := spirali.Create(vg, name, config, metadata)
	if err != nil {
		panic(err)
	}
	if err := updateMetaData(newMetaData, config); err != nil {
		panic(err)
	}
}
