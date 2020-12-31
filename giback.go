package main

import (
	"log"
	"os"

	"github.com/dhuan/giback/pkg/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "giback",
		Version: getVersion(),
		Usage:   "Easily backup any files to git repositories.",
		Action:  cmd.Default,
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
			&cli.BoolFlag{
				Name:  "v",
				Usage: "Verbose, print all shell operations.",
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

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	version := os.Getenv("GIBACK_VERSION")

	if version == "" {
		return "DEVELOPMENT"
	}

	return version
}
