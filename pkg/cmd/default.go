package cmd

import (
	"fmt"
	"log"

	"github.com/dhuan/giback/pkg/app"

	"github.com/urfave/cli/v2"
)

func Default(c *cli.Context) error {
	unitId := c.Args().First()

	if unitId == "" {
		log.Println(cli.ShowAppHelp(c))

		return nil
	}

	config, workspacePath, shellRunOptions, err := runUnitPrepare(c, unitId)

	if err != nil {
		log.Fatal(err)
	}

	unit, err := getUnitById(unitId, config.Units)

	if err != nil {
		log.Fatal(err)
	}

	runErr := runUnit(unit, workspacePath, shellRunOptions)

	if runErr != nil {
		log.Fatal(runErr)
	}

	return nil
}

func getUnitById(unitId string, units []app.PushUnit) (app.PushUnit, error) {
	for i := range units {
		if units[i].Id == unitId {
			return units[i], nil
		}
	}

	return app.PushUnit{}, fmt.Errorf("Could not find unit with ID '%s'.", unitId)
}
