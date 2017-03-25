package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kudohamu/spirali"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "spirali"
	app.Usage = "golang based migration tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env, e",
			Usage: "config environment to use",
		},
		cli.StringFlag{
			Name:  "path, p",
			Usage: "migration dir containing config.yml",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create migration file for `up` and `down`",
			Action:  Create,
		},
	}

	app.Run(os.Args)
}

func initializeMetaDataFileIfNotExist(dir string) error {
	p := filepath.Join(dir, spirali.MetaDataFileName)

	if _, err := os.Stat(p); err == nil {
		return nil
	}

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()

	var m spirali.MetaData
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	if _, err := file.Write(b); err != nil {
		return err
	}
	return nil
}
