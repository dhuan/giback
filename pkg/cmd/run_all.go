package cmd

import (
	"log"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/gibackfs"

	"github.com/urfave/cli/v2"
)

func RunAll(c *cli.Context) error {
	appContext, err := app.BuildContext(c)

	if err != nil {
		log.Fatal(err)
	}

	config, workspacePath, err := gibackfs.GetUserConfig(appContext)

	if err != nil {
		log.Fatal(err)
	}

	for _, unit := range config.Units {
		runErr := runUnit(unit, workspacePath)

		if runErr != nil {
			log.Fatal(runErr)
		}
	}

	return nil
}
