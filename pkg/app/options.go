package app

import (
	"fmt"
	"log"
	"os/user"

	"github.com/urfave/cli/v2"
)

type Context struct {
	ConfigFilePath string
	WorkspacePath  string
	Verbose        bool
}

func BuildContext(c *cli.Context) (Context, error) {
	usr, err := user.Current()

	if err != nil {
		return Context{}, err
	}

	var workspacePath string
	var configFilePath string
	var verbose bool

	if c.IsSet("c") {
		configFilePath = c.String("c")
	} else {
		configFilePath = fmt.Sprintf("%s/.giback.yml", usr.HomeDir)
	}

	if c.IsSet("w") {
		workspacePath = c.String("w")
	} else {
		workspacePath = fmt.Sprintf("%s/.giback", usr.HomeDir)
	}

	if c.IsSet("v") {
		verbose = c.Bool("v")
	} else {
		verbose = false
	}

	if !fileExists(configFilePath) {
		log.Fatal(fmt.Sprintf("Config file does not exist:  %s", configFilePath))
	}

	if !folderExists(workspacePath) {
		log.Fatal(fmt.Sprintf("Workspace folder does not exist:  %s", workspacePath))
	}

	return Context{configFilePath, workspacePath, verbose}, nil
}
