package main

import (
	mib_app "github.com/petara94/mib-go/internal/app"
	"github.com/petara94/mib-go/internal/db"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

var Version = "unknown"

func main() {
	app := cli.App{
		Usage:    "auth encrypted db",
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "db",
				Usage:   "db file",
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    "pass",
				Usage:   "db password",
				Aliases: []string{"p"},
			},
		},
		Version: Version,
		Action: func(c *cli.Context) error {
			repo := db.NewDB(c.String("db"), c.String("pass"))
			app := mib_app.NewApp(repo)

			if err := app.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
