package git

import (
	"fmt"

	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/shell"
	"github.com/dhuan/giback/pkg/utils"
)

func CheckAccess(workspace string, repository string, shellRunOptions shell.RunOptions) bool {
	if !gibackfs.FolderExists(workspace) {
		return false
	}

	command := "git ls-remote " + repository

	_, runErr := shell.Run(workspace, command, nil, shellRunOptions)

	if runErr != nil {
		return false
	}

	return true
}

func Clone(workspace string, repositoryPath string, saveAs string, shellRunOptions shell.RunOptions) error {
	command := fmt.Sprintf("git clone %s %s", repositoryPath, saveAs)

	_, err := shell.Run(workspace, command, nil, shellRunOptions)

	return err
}

func Pull(repositoryPath string, shellRunOptions shell.RunOptions) error {
	commandFetch := "git fetch"

	_, errFetch := shell.Run(repositoryPath, commandFetch, nil, shellRunOptions)

	if errFetch != nil {
		return errFetch
	}

	commandRebase := "git rebase origin/master"

	_, errRebase := shell.Run(repositoryPath, commandRebase, nil, shellRunOptions)

	if errRebase != nil {
		return errRebase
	}

	return nil
}

func Status(repositoryPath string, shellRunOptions shell.RunOptions) []GitStatusResult {
	var result []GitStatusResult

	command := fmt.Sprintf("git status --short")

	statusOutput, _ := shell.Run(repositoryPath, command, nil, shellRunOptions)

	statusFiles := utils.SedReplaceGlobal(string(statusOutput), `^...`, "")

	statusFiles = utils.RemoveEmptyLines(statusFiles)

	for _, statusFile := range statusFiles {
		result = append(result, GitStatusResult{statusFile})
	}

	return result
}

func AddAll(repositoryPath string, shellRunOptions shell.RunOptions) error {
	_, err := shell.Run(repositoryPath, "git add .", nil, shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func Commit(repositoryPath string, message string, authorName string, authorEmail string, shellRunOptions shell.RunOptions) error {
	var env map[string]string

	env = make(map[string]string)

	env["GIT_COMMITTER_NAME"] = authorName
	env["GIT_COMMITTER_EMAIL"] = authorEmail
	env["GIT_AUTHOR_NAME"] = authorName
	env["GIT_AUTHOR_EMAIL"] = authorEmail

	command := fmt.Sprintf("git commit -m \"%s\"", message)

	_, err := shell.Run(repositoryPath, command, env, shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func Push(repositoryPath string, shellRunOptions shell.RunOptions) error {
	_, err := shell.Run(repositoryPath, "git push origin master", nil, shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

type GitStatusResult struct {
	File string
}
