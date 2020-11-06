package giback_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dhuan/giback/pkg/testutils"
)

func TestBackupSuccessfully(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, _ := testutils.RunGiback("my_backup", testutils.RunGibackOptions{})

	testutils.AssertOutput(t, output, []string{
		"Running unit 'my_backup'.",
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost:2222/srv/git/test.git",
		"Identifying files...",
		withFullPath("backup_files/another_file.txt"),
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Files affected in the repository: another_file.txt,some_file.txt",
		"Committing: Backing up!",
		"Pushing...",
		"Done!",
	})

	testutils.AssertGitLog(t, "test", []string{
		"Super Man <superman@example.com> Backing up!",
	})
}

func TestBackupSuccessfullyAll(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, _ := testutils.RunGiback("all", testutils.RunGibackOptions{})

	testutils.AssertOutput(t, output, []string{
		"Running unit 'my_backup'.",
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost:2222/srv/git/test.git",
		"Identifying files...",
		withFullPath("backup_files/another_file.txt"),
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Files affected in the repository: another_file.txt,some_file.txt",
		"Committing: Backing up!",
		"Pushing...",
		"Done!",
		"Running unit 'another_backup'.",
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost:2222/srv/git/test2.git",
		"Identifying files...",
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Files affected in the repository: some_file.txt",
		"Committing: Backup.",
		"Pushing...",
		"Done!",
	})

	testutils.AssertGitLog(t, "test", []string{
		"Super Man <superman@example.com> Backing up!",
	})

	testutils.AssertGitLog(t, "test2", []string{
		"Someone <someone@example.com> Backup.",
	})
}

func TestBackupWithNoChanges(t *testing.T) {
	testutils.ResetTestEnvironment()

	testutils.RunGiback("my_backup", testutils.RunGibackOptions{})

	output, _ := testutils.RunGiback("my_backup", testutils.RunGibackOptions{})

	testutils.AssertGitLog(t, "test", []string{
		"Super Man <superman@example.com> Backing up!",
	})

	testutils.AssertOutput(t, output, []string{
		"Running unit 'my_backup'.",
		"Pulling git changes.",
		"Identifying files...",
		withFullPath("backup_files/another_file.txt"),
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Nothing has changed. No commit will be pushed to this repository.",
	})

	testutils.AssertGitLog(t, "test", []string{
		"Super Man <superman@example.com> Backing up!",
	})
}

func TestFailRunningUnexistingUnit(t *testing.T) {
	output, err := testutils.RunGiback("my_unexisting_backup", testutils.RunGibackOptions{})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"Could not find unit with ID 'my_unexisting_backup'.",
	})
}

func TestFailRunningAllWithUnmatchingRepositories(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, err := testutils.RunGiback("all", testutils.RunGibackOptions{})

	testutils.AssertHasNoError(t, err)

	output, err = testutils.RunGiback("all", testutils.RunGibackOptions{
		ConfigFile: "invalid",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"The following repositories defined in your configuration don't match with the ones located in your workspace: another_backup.",
	})
}

func TestFailRunningSingleWithUnmatchingRepositories(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, err := testutils.RunGiback("all", testutils.RunGibackOptions{})

	testutils.AssertHasNoError(t, err)

	output, err = testutils.RunGiback("another_backup", testutils.RunGibackOptions{
		ConfigFile: "invalid",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"The repository configured for \"another_backup\" does not match with the one cloned at your workspace.",
	})
}

func withFullPath(path string) string {
	pwd, _ := os.Getwd()

	return fmt.Sprintf("%s/%s", pwd, path)
}
