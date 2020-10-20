package git

import (
	"fmt"

	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/shell"
	"github.com/dhuan/giback/pkg/utils"
)

func CheckAccess(workspace string, repository string) bool {
	if !gibackfs.FolderExists(workspace) {
		return false
	}

	command := "git ls-remote " + repository

	_, runErr := shell.Run(workspace, command, nil)

	if runErr != nil {
		return false
	}

	return true
}

func Clone(workspace string, repositoryPath string, saveAs string) error {
	command := fmt.Sprintf("git clone %s %s", repositoryPath, saveAs)

	_, err := shell.Run(workspace, command, nil)

	return err
}

func Pull(repositoryPath string) error {
	commandFetch := "git fetch"

	_, errFetch := shell.Run(repositoryPath, commandFetch, nil)

	if errFetch != nil {
		return errFetch
	}

	commandRebase := "git rebase origin/master"

	_, errRebase := shell.Run(repositoryPath, commandRebase, nil)

	if errRebase != nil {
		return errRebase
	}

	return nil
}

func Status(repositoryPath string) []GitStatusResult {
	var result []GitStatusResult

	command := fmt.Sprintf("git status --short")

	statusOutput, _ := shell.Run(repositoryPath, command, nil)

	statusFiles := utils.SedReplaceGlobal(string(statusOutput), `^...`, "")

	statusFiles = utils.RemoveEmptyLines(statusFiles)

	for _, statusFile := range statusFiles {
		result = append(result, GitStatusResult{statusFile})
	}

	return result
}

func AddAll(repositoryPath string) error {
	_, err := shell.Run(repositoryPath, "git add .", nil)

	if err != nil {
		return err
	}

	return nil
}

func Commit(repositoryPath string, message string, authorName string, authorEmail string) error {
	var env map[string]string

	env = make(map[string]string)

	env["GIT_COMMITTER_NAME"] = authorName
	env["GIT_COMMITTER_EMAIL"] = authorEmail
	env["GIT_AUTHOR_NAME"] = authorName
	env["GIT_AUTHOR_EMAIL"] = authorEmail

	command := fmt.Sprintf("git commit -m \"%s\"", message)

	_, err := shell.Run(repositoryPath, command, env)

	if err != nil {
		return err
	}

	return nil
}

func Push(repositoryPath string) error {
	_, err := shell.Run(repositoryPath, "git push origin master", nil)

	if err != nil {
		return err
	}

	return nil
}

type GitStatusResult struct {
	file string
}
