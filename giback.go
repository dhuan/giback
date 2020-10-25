package main

import (
	"log"
	"os"

	"github.com/dhuan/giback/pkg/cmd"
	"github.com/dhuan/giback/pkg/shell"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "giback",
		Usage:  "Easily backup any files to git repositories.",
		Action: cmd.Default,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "c",
				Value: "",
				Usage: "Path to a configuration file.",
			},
			&cli.StringFlag{
				Name:  "w",
				Value: "",
				Usage: "Path to workspace.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "all",
				Usage:  "Run all units",
				Action: cmd.RunAll,
			},
		},
	}

	checkDependencies()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkDependencies() {
	_, err := shell.Run("", "which git", nil)

	if err != nil {
		log.Fatal("Giback requires git. Please make sure you have it installed before trying again.")
	}
}
