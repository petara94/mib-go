package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

var version = "unknown"

func main() {
	app := cli.App{
		Usage:    "auth encrypted db",
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "db",
				Usage: "db file",
			},
		},
		Version: version,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
