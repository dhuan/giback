package giback_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dhuan/giback/pkg/testutils"
)

func TestBackupSuccessfully(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, _ := testutils.RunGiback("my_backup")

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

	testutils.AssertGitLog(t, []string{
		"Super Man <superman@example.com> Backing up!",
	})
}

func TestBackupWithNoChanges(t *testing.T) {
	testutils.ResetTestEnvironment()

	testutils.RunGiback("my_backup")

	output, _ := testutils.RunGiback("my_backup")

	testutils.AssertGitLog(t, []string{
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

	testutils.AssertGitLog(t, []string{
		"Super Man <superman@example.com> Backing up!",
	})
}

func TestFailRunningUnexistingUnit(t *testing.T) {
	output, err := testutils.RunGiback("my_unexisting_backup")

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"Could not find unit with ID 'my_unexisting_backup'.",
	})
}

func withFullPath(path string) string {
	pwd, _ := os.Getwd()

	return fmt.Sprintf("%s/%s", pwd, path)
}
