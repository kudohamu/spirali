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
		cli.Command{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "apply migrations",
			Action:  Up,
		},
		cli.Command{
			Name:    "down",
			Aliases: []string{"d"},
			Usage:   "rollback the latest migration",
			Action:  Down,
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

func openConfig(dir string) (*spirali.Config, error) {
	p := filepath.Join(dir, spirali.ConfigFileName)

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c, err := spirali.ReadConfig(file)
	if err != nil {
		return nil, err
	}
	return c, nil
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
