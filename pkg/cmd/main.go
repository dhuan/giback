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

	log.Println(fmt.Sprintf("Running unit '%s'.", unit.Id))

	hasAccess := git.CheckAccess(workspacePath, unit.Repository)

	repositoryPath := getRepositoryPath(workspacePath, unit)

	if !hasAccess {
		return errors.New("You don't have access to this repository.")
	}

	hasCloned := hasClonedRepository(workspacePath, unit.Id)

	if !hasCloned {
		log.Println(fmt.Sprintf("Repository has not been cloned yet, yet. Will clone now: %s", unit.Repository))
		git.Clone(workspacePath, unit.Repository, unit.Id)
	}

	if hasCloned {
		log.Println(fmt.Sprintf("Pulling git changes."))
		git.Pull(repositoryPath)
	}

	log.Println(fmt.Sprintf("Identifying files..."))

	files := gibackfs.ScanDirMany(unit.Files)

	for i := range files {
		log.Println(fmt.Sprintf("%s", files[i]))
	}

	gibackfs.Copy(files, repositoryPath)

	log.Println(fmt.Sprintf("Files copied."))

	statusResult := git.Status(repositoryPath)

	if len(statusResult) == 0 {
		log.Println("Nothing has changed. No commit will be pushed to this repository.")

		return nil
	}

	err = git.AddAll(repositoryPath)

	if err != nil {
		log.Fatal("Failed to add.")
	}

	log.Println(fmt.Sprintf("Committing: %s", unit.Commit_Message))

	err = git.Commit(repositoryPath, unit.Commit_Message, unit.Author_Name, unit.Author_Email)

	if err != nil {
		log.Fatal("Failed to commit.", err)
	}

	log.Println(fmt.Sprintf("Pushing..."))

	err = git.Push(repositoryPath)

	if err != nil {
		log.Fatal("Failed to push.", err)
	}

	log.Println(fmt.Sprintf("Done!"))

	return nil
}

func hasClonedRepository(workspace string, id string) bool {
	repositoryPath := workspace + "/" + id

	return gibackfs.FolderExists(repositoryPath)
}

func getRepositoryPath(workspace string, unit app.PushUnit) string {
	return workspace + "/" + unit.Id
}
