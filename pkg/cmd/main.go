package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/git"

	"github.com/urfave/cli/v2"
)

func Main(c *cli.Context) error {
	config, workspacePath, err := gibackfs.GetUserConfig()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nresult: ")
	fmt.Printf("%+v\n", config)

	for _, unit := range config.Units {
		runErr := runUnit(unit, workspacePath)

		if runErr != nil {
			log.Fatal(runErr)
		}
	}

	return nil
}

func runUnit(unit app.PushUnit, workspacePath string) error {
	var err error

	hasAccess := git.CheckAccess(workspacePath, unit.Key, unit.Repository)

	repositoryPath := getRepositoryPath(workspacePath, unit)

	if !hasAccess {
		return errors.New("You don't have access to this repository.")
	}

	hasCloned := hasClonedRepository(workspacePath, unit.Id)

	if !hasCloned {
		git.Clone(workspacePath, unit.Repository, unit.Id)
	}

	if hasCloned {
		git.Pull(repositoryPath)
	}

	files := gibackfs.ScanDir(unit.Files[0])

	gibackfs.Copy(files, repositoryPath)

	statusResult := git.Status(repositoryPath)

	if len(statusResult) == 0 {
		log.Println("Nothing has changed. No commit will be pushed to this repository.")

		return nil
	}

	err = git.AddAll(repositoryPath)

	if err != nil {
		log.Fatal("Failed to add.")
	}

	err = git.Commit(repositoryPath, "some commit message", "fulano silva", "fulano.silva@example.com")

	if err != nil {
		log.Fatal("Failed to commit.", err)
	}

	err = git.Push(repositoryPath)

	if err != nil {
		log.Fatal("Failed to push.", err)
	}

	return nil
}

func hasClonedRepository(workspace string, id string) bool {
	repositoryPath := workspace + "/" + id

	return gibackfs.FolderExists(repositoryPath)
}

func getRepositoryPath(workspace string, unit app.PushUnit) string {
	return workspace + "/" + unit.Id
}
