package main

import (
	"log"
	"os"

	"github.com/dhuan/giback/pkg/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "giback",
		Usage:  "Easily backup any files to git repositories.",
		Action: cmd.Default,
		Commands: []*cli.Command{
			{
				Name:   "all",
				Usage:  "Run all units",
				Action: cmd.RunAll,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
