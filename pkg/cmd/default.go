package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/gibackfs"

	"github.com/urfave/cli/v2"
)

func Default(c *cli.Context) error {
	config, workspacePath, err := gibackfs.GetUserConfig()

	if err != nil {
		log.Fatal(err)
	}

	unitId := c.Args().First()

	if unitId == "" {
		log.Println("Nothing to do.")

		return nil
	}

	unit, err := getUnitById(unitId, config.Units)

	if err != nil {
		log.Fatal(err)
	}

	runErr := runUnit(unit, workspacePath)

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

	return app.PushUnit{}, errors.New(fmt.Sprintf("Could not find unit with ID '%s'.", unitId))
}
