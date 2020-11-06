package cmd

import (
	"log"

	"github.com/urfave/cli/v2"
)

func RunAll(c *cli.Context) error {
	config, workspacePath, shellRunOptions, err := runUnitPrepare(c, "")

	if err != nil {
		log.Fatal(err)
	}

	for _, unit := range config.Units {
		runErr := runUnit(unit, workspacePath, shellRunOptions)

		if runErr != nil {
			log.Fatal(runErr)
		}
	}

	return nil
}
