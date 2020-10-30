package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/git"
	"github.com/dhuan/giback/pkg/utils"
)

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
		log.Println(fmt.Sprintf("Repository has not been cloned yet. Will clone now: %s", unit.Repository))

		err = git.Clone(workspacePath, unit.Repository, unit.Id)

		if err != nil {
			return errors.New("Failed to clone repository.")
		}
	}

	if hasCloned {
		log.Println(fmt.Sprintf("Pulling git changes."))
		git.Pull(repositoryPath)
	}

	log.Println(fmt.Sprintf("Identifying files..."))

	pwd, _ := os.Getwd()

	vars := make(map[string]string)
	vars["PWD"] = pwd

	filePatterns := evaluateMany(vars, unit.Files)

	files := gibackfs.ScanDirMany(filePatterns)

	if len(unit.Exclude) > 0 {
		files = utils.FilterOut(files, unit.Exclude)
	}

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

	fileList := buildFileListFromStatusResult(statusResult)

	log.Println(fmt.Sprintf("Files affected in the repository: %s", fileList))

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

func buildFileListFromStatusResult(statusResults []git.GitStatusResult) string {
	var fileList []string

	for i := range statusResults {
		fileList = append(fileList, statusResults[i].File)
	}

	return strings.Join(fileList, ",")
}

func evaluateMany(vars map[string]string, strings []string) []string {
	stringsTransformed := make([]string, len(strings))

	for i := range strings {
		stringsTransformed[i] = evaluate(vars, strings[i])
	}

	return stringsTransformed
}

func evaluate(vars map[string]string, str string) string {
	result := str

	for i, v := range vars {
		search := fmt.Sprintf("{%s}", i)
		replace := v

		result = strings.ReplaceAll(str, search, replace)
	}

	return result
}
