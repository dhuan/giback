package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/git"
	"github.com/dhuan/giback/pkg/shell"
	"github.com/dhuan/giback/pkg/utils"
)

func runUnit(unit app.PushUnit, workspacePath string, shellRunOptions shell.RunOptions) error {
	var err error

	log.Println(fmt.Sprintf("Running unit '%s'.", unit.Id))

	hasAccess := git.CheckAccess(workspacePath, unit.Repository, shellRunOptions)

	repositoryPath := getRepositoryPath(workspacePath, unit)

	if !hasAccess {
		return errors.New("You don't have access to this repository.")
	}

	hasCloned := hasClonedRepository(workspacePath, unit.Id)

	if !hasCloned {
		log.Println(fmt.Sprintf("Repository has not been cloned yet. Will clone now: %s", unit.Repository))

		err = git.Clone(workspacePath, unit.Repository, unit.Id, shellRunOptions)

		if err != nil {
			return errors.New("Failed to clone repository.")
		}
	}

	if hasCloned {
		log.Println(fmt.Sprintf("Pulling git changes."))
		git.Pull(repositoryPath, shellRunOptions)
	}

	log.Println(fmt.Sprintf("Identifying files..."))

	pwd, _ := os.Getwd()

	vars := make(map[string]string)
	vars["PWD"] = pwd

	filePatterns := evaluateMany(vars, unit.Files)

	excludePatterns := evaluateMany(vars, unit.Exclude)

	files := gibackfs.ScanDirMany(filePatterns, excludePatterns)

	fileNames := utils.GetFileNameMany(files)

	statusBeforeCopyResult := git.Status(repositoryPath, shellRunOptions)

	statusFilesBeforeCopy := buildFileListFromStatusResult(statusBeforeCopyResult)

	for i := range files {
		log.Println(fmt.Sprintf("%s", files[i]))
	}

	gibackfs.Copy(files, repositoryPath)

	log.Println(fmt.Sprintf("Files copied."))

	statusResult := git.Status(repositoryPath, shellRunOptions)

	ignoredFiles := utils.StringListDiff(statusFilesBeforeCopy, fileNames)

	filesToBeAdded := utils.FilterOutMatching(
		buildFileListFromStatusResult(statusResult),
		ignoredFiles,
	)

	if len(filesToBeAdded) == 0 {
		log.Println("Nothing has changed. No commit will be pushed to this repository.")

		return nil
	}

	log.Println(fmt.Sprintf(
		"Files affected in the repository: %s",
		strings.Join(filesToBeAdded, ","),
	))

	if len(ignoredFiles) > 0 {
		log.Println(fmt.Sprintf(
			"The following files located in the repository folder will not be commited as they are not included in the backup files: %s",
			strings.Join(ignoredFiles, ","),
		))
	}

	if err = git.Reset(repositoryPath, shellRunOptions); err != nil {
		log.Fatal("Failed to reset.")
	}

	err = git.Add(repositoryPath, filesToBeAdded, shellRunOptions)

	if err != nil {
		log.Fatal("Failed to add.")
	}

	log.Println(fmt.Sprintf("Committing: %s", unit.Commit_Message))

	err = git.Commit(repositoryPath, unit.Commit_Message, unit.Author_Name, unit.Author_Email, shellRunOptions)

	if err != nil {
		log.Fatal("Failed to commit.", err)
	}

	log.Println(fmt.Sprintf("Pushing..."))

	err = git.Push(repositoryPath, shellRunOptions)

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

func buildFileListFromStatusResult(statusResults []git.GitStatusResult) []string {
	var fileList []string

	for i := range statusResults {
		fileList = append(fileList, statusResults[i].File)
	}

	return fileList
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

func buildShellRunOptions(context app.Context) shell.RunOptions {
	return shell.RunOptions{
		Debug: context.Verbose,
	}
}

func checkDependencies(shellRunOptions shell.RunOptions) {
	_, err := shell.Run("", "which git", nil, shellRunOptions)

	if err != nil {
		log.Fatal("Giback requires git. Please make sure you have it installed before trying again.")
	}
}

func checkRepos(config app.Config, workspacePath string, shellRunOptions shell.RunOptions) ([]string, error) {
	var invalidRepositories []string

	for i := range config.Units {
		unit := config.Units[i]
		repositoryPath := fmt.Sprintf("%s/%s", workspacePath, unit.Id)

		if !gibackfs.FolderExists(repositoryPath) {
			continue
		}

		repositoryMetadata, err := git.GetRepositoryMetadata(repositoryPath, shellRunOptions)

		if err != nil {
			return nil, err
		}

		if repositoryMetadata.Address != unit.Repository {
			invalidRepositories = append(invalidRepositories, unit.Id)
		}
	}

	return invalidRepositories, nil
}

func runUnitPrepare(c *cli.Context, unitId string) (app.Config, string, shell.RunOptions, error) {
	allMode := false

	if unitId == "" {
		allMode = true
	}

	appContext, err := app.BuildContext(c)
	if err != nil {
		return app.Config{}, "", shell.RunOptions{}, err
	}

	shellRunOptions := buildShellRunOptions(appContext)

	checkDependencies(shellRunOptions)

	config, workspacePath, err := gibackfs.GetUserConfig(appContext)
	if err != nil {
		return app.Config{}, "", shell.RunOptions{}, err
	}

	units := config.Units

	if !allMode {
		unit, err := getUnitById(unitId, config.Units)

		if err != nil {
			return app.Config{}, "", shell.RunOptions{}, err
		}

		units = []app.PushUnit{unit}
	}

	invalidUnits, invalidUnitsIds := app.ValidateUnits(units)

	if allMode && len(invalidUnits) > 0 {
		log.Fatal(fmt.Sprintf(
			"The following units are not configured properly:\n\n%s\n\nCheck the manual to find out how to properly configure Giback.",
			buildInvalidUnitsErrorMessage(invalidUnits, invalidUnitsIds),
		))
	}

	if !allMode && len(invalidUnits) > 0 {
		invalidUnitErrorMessage := invalidUnits[0]

		log.Fatal(fmt.Sprintf(
			"'%s' is not configured properly:\n\n%s\n\nCheck the manual to find out how to properly configure Giback.",
			unitId,
			invalidUnitErrorMessage,
		))
	}

	invalidRepos, err := checkRepos(config, workspacePath, shellRunOptions)
	if err != nil {
		return app.Config{}, "", shell.RunOptions{}, err
	}

	if !allMode && len(invalidRepos) > 0 && utils.IndexOfString(invalidRepos, unitId) > -1 {
		return app.Config{}, "", shell.RunOptions{}, errors.New(fmt.Sprintf(
			"The repository configured for \"%s\" does not match with the one cloned at your workspace.",
			unitId,
		))
	}

	if allMode && len(invalidRepos) > 0 {
		return app.Config{}, "", shell.RunOptions{}, errors.New(fmt.Sprintf(
			"The following repositories defined in your configuration don't match with the ones located in your workspace: %s.",
			strings.Join(invalidRepos, ","),
		))
	}

	return config, workspacePath, shellRunOptions, nil
}

func buildInvalidUnitsErrorMessage(invalidUnitErrorMessage []string, unitIds []string) string {
	var messages []string

	for i := range invalidUnitErrorMessage {
		messages = append(messages, fmt.Sprintf("%s:\n%s", unitIds[i], invalidUnitErrorMessage[i]))
	}

	return strings.Join(messages, "\n\n")
}
