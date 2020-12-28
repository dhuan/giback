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
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost/srv/git/test.git",
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
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost/srv/git/test.git",
		"Identifying files...",
		withFullPath("backup_files/another_file.txt"),
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Files affected in the repository: another_file.txt,some_file.txt",
		"Committing: Backing up!",
		"Pushing...",
		"Done!",
		"Running unit 'another_backup'.",
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost/srv/git/test2.git",
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

func TestBackupSuccessfullyWithIgnoredFiles(t *testing.T) {
	testutils.ResetTestEnvironment()

	_, err := testutils.RunGiback("my_backup", testutils.RunGibackOptions{})

	testutils.AssertHasNoError(t, err)

	testutils.AddFileToWorkspace("my_backup", []string{
		"this_file_shall_not_be_commited.html File content.",
	})

	testutils.AddFileToBackupFilesFolder([]string{
		"a_new_file.txt File content.",
	})

	output, err := testutils.RunGiback("my_backup", testutils.RunGibackOptions{})

	testutils.AssertOutput(t, output, []string{
		"Running unit 'my_backup'.",
		"Pulling git changes.",
		"Identifying files...",
		withFullPath("backup_files/a_new_file.txt"),
		withFullPath("backup_files/another_file.txt"),
		withFullPath("backup_files/some_file.txt"),
		"Files copied.",
		"Files affected in the repository: a_new_file.txt",
		"The following files located in the repository folder will not be commited as they are not included in the backup files: this_file_shall_not_be_commited.html",
		"Committing: Backing up!",
		"Pushing...",
		"Done!",
	})

	testutils.AssertGitLog(t, "test", []string{
		"Super Man <superman@example.com> Backing up!",
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

func TestFailRunningAllWithInvalidConfigMissingFields(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, err := testutils.RunGiback("all", testutils.RunGibackOptions{
		ConfigFile: "invalid_missing_fields",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"The following units are not configured properly:",
		"another_backup:",
		"Missing the following fields: repository,files",
		"Check the manual to find out how to properly configure Giback.",
	})
}

func TestFailRunningSingleWithInvalidConfigMissingFields(t *testing.T) {
	testutils.ResetTestEnvironment()

	_, err := testutils.RunGiback("my_backup", testutils.RunGibackOptions{
		ConfigFile: "invalid_missing_fields",
	})

	testutils.AssertHasNoError(t, err)

	output, err := testutils.RunGiback("another_backup", testutils.RunGibackOptions{
		ConfigFile: "invalid_missing_fields",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"'another_backup' is not configured properly:",
		"Missing the following fields: repository,files",
		"Check the manual to find out how to properly configure Giback.",
	})
}

func TestBackupSuccessfullyWithCustomKey(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, _ := testutils.RunGiback("my_backup", testutils.RunGibackOptions{
		ConfigFile: "with_keys",
	})

	testutils.AssertOutput(t, output, []string{
		"Running unit 'my_backup'.",
		"Repository has not been cloned yet. Will clone now: ssh://git@localhost/srv/git/test.git",
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

func TestFailSingleWithCustomKey(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, err := testutils.RunGiback("another_backup", testutils.RunGibackOptions{
		ConfigFile: "with_keys",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"The following repositories failed to be communicated with. Please verify that you indeed have access to these repositories.",
		"ssh://git@localhost/srv/git/test2.git",
	})
}

func TestFailAllWithCustomKey(t *testing.T) {
	testutils.ResetTestEnvironment()

	output, err := testutils.RunGiback("all", testutils.RunGibackOptions{
		ConfigFile: "with_keys",
	})

	testutils.AssertHasError(t, err)

	testutils.AssertOutput(t, output, []string{
		"The following repositories failed to be communicated with. Please verify that you indeed have access to these repositories.",
		"ssh://git@localhost/srv/git/test2.git",
	})
}

func withFullPath(path string) string {
	pwd, _ := os.Getwd()

	return fmt.Sprintf("%s/%s", pwd, path)
}
