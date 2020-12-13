package git

import (
	"fmt"
	"strings"

	"github.com/dhuan/giback/pkg/gibackfs"
	"github.com/dhuan/giback/pkg/shell"
	"github.com/dhuan/giback/pkg/utils"
)

func CheckAccess(workspace string, repository string, shellRunOptions shell.RunOptions) bool {
	if !gibackfs.FolderExists(workspace) {
		return false
	}

	command := "git ls-remote " + repository

	_, runErr := shell.Run(workspace, command, shellRunOptions)

	if runErr != nil {
		return false
	}

	return true
}

func Clone(workspace string, repositoryPath string, saveAs string, shellRunOptions shell.RunOptions) error {
	command := fmt.Sprintf("git clone %s %s", repositoryPath, saveAs)

	_, err := shell.Run(workspace, command, shellRunOptions)

	return err
}

func Pull(repositoryPath string, shellRunOptions shell.RunOptions) error {
	commandFetch := "git fetch"

	_, errFetch := shell.Run(repositoryPath, commandFetch, shellRunOptions)

	if errFetch != nil {
		return errFetch
	}

	commandRebase := "git rebase origin/master"

	_, errRebase := shell.Run(repositoryPath, commandRebase, shellRunOptions)

	if errRebase != nil {
		return errRebase
	}

	return nil
}

func Status(repositoryPath string, shellRunOptions shell.RunOptions) []GitStatusResult {
	var result []GitStatusResult

	command := fmt.Sprintf("git status --short")

	statusOutput, _ := shell.Run(repositoryPath, command, shellRunOptions)

	statusFiles := utils.SedReplaceGlobal(string(statusOutput), `^...`, "")

	statusFiles = utils.RemoveEmptyLines(statusFiles)

	for _, statusFile := range statusFiles {
		result = append(result, GitStatusResult{statusFile})
	}

	return result
}

func LsRemote(repository string, shellRunOptions shell.RunOptions) error {
	command := fmt.Sprintf("git ls-remote %s", repository)

	_, err := shell.Run("", command, shellRunOptions)

	return err
}

func Reset(repositoryPath string, shellRunOptions shell.RunOptions) error {
	_, err := shell.Run(repositoryPath, "git reset", shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func AddAll(repositoryPath string, shellRunOptions shell.RunOptions) error {
	_, err := shell.Run(repositoryPath, "git add .", shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func Add(repositoryPath string, files []string, shellRunOptions shell.RunOptions) error {
	filesJoined := strings.Join(files, " ")

	command := fmt.Sprintf("git add %s", filesJoined)

	_, err := shell.Run(repositoryPath, command, shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func Commit(repositoryPath string, message string, authorName string, authorEmail string, shellRunOptions shell.RunOptions) error {
	shellRunOptionsModifiedForCommit := shellRunOptions

	if len(shellRunOptionsModifiedForCommit.Env) == 0 {
		shellRunOptionsModifiedForCommit.Env = make(map[string]string)
	}

	shellRunOptionsModifiedForCommit.Env["GIT_COMMITTER_NAME"] = authorName
	shellRunOptionsModifiedForCommit.Env["GIT_COMMITTER_EMAIL"] = authorEmail
	shellRunOptionsModifiedForCommit.Env["GIT_AUTHOR_NAME"] = authorName
	shellRunOptionsModifiedForCommit.Env["GIT_AUTHOR_EMAIL"] = authorEmail

	command := fmt.Sprintf("git commit -m \"%s\"", message)

	_, err := shell.Run(repositoryPath, command, shellRunOptionsModifiedForCommit)

	if err != nil {
		return err
	}

	return nil
}

func Push(repositoryPath string, shellRunOptions shell.RunOptions) error {
	_, err := shell.Run(repositoryPath, "git push origin master", shellRunOptions)

	if err != nil {
		return err
	}

	return nil
}

func GetRepositoryMetadata(repositoryPath string, shellRunOptions shell.RunOptions) (RepositoryMetadata, error) {
	if !gibackfs.FolderExists(repositoryPath) {
		return RepositoryMetadata{}, nil
	}

	output, err := shell.Run(repositoryPath, "git config --get remote.origin.url", shellRunOptions)

	if err != nil {
		return RepositoryMetadata{}, err
	}

	address := strings.TrimSpace(string(output))

	return RepositoryMetadata{address}, nil
}

type GitStatusResult struct {
	File string
}

type RepositoryMetadata struct {
	Address string
}
