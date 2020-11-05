package testutils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/dhuan/giback/pkg/shell"
)

func ResetTestEnvironment() {
	workingDir, _ := os.Getwd()

	gibackRootPath := fmt.Sprintf("%s/..", workingDir)

	resetTestRepository(gibackRootPath)

	emptyWorkspace(gibackRootPath)
}

func RunGiback(command string) ([]byte, error) {
	workspacePath := "./test/tmp/workspace"
	configPath := "./test/config/default.yml"
	workingDir, _ := os.Getwd()
	gibackRootPath := fmt.Sprintf("%s/..", workingDir)

	commandTranformed := fmt.Sprintf("./giback -w %s -c %s %s", workspacePath, configPath, command)

	output, err := shell.Run(gibackRootPath, commandTranformed, nil, shell.RunOptionsDefault())

	if err != nil {
		return output, err
	}

	return output, nil
}

func AssertHasError(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error to have been resulted.")
	}
}

func AssertOutput(t *testing.T, output []byte, expectedOutput []string) {
	outputLines := strings.Split(string(output), "\n")

	outputLines = cleanupLogMessages(outputLines)

	outputLines = cleanupEmptyLines(outputLines)

	AssertArraysEqual(t, expectedOutput, outputLines)
}

func AssertArraysEqual(t *testing.T, expected []string, result []string) {
	resultLinesJoined := strings.Join(result, "\n")

	if len(result) != len(expected) {
		t.Error(fmt.Sprintf("Output lines count does not match expected output! Output:\n\n%s", resultLinesJoined))

		return
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Error(fmt.Sprintf("Lines do not match!\nOutput:   %s\nExpected: %s\n\nWhole output:\n%s", result[i], expected[i], resultLinesJoined))
		}
	}
}

func AssertGitLog(t *testing.T, repositoryFolder string, expectedLog []string) {
	gitLog := getGitLog(repositoryFolder)

	AssertArraysEqual(t, expectedLog, gitLog)
}

func getGitLog(repositoryFolder string) []string {
	workingDir, _ := os.Getwd()
	gibackRootPath := fmt.Sprintf("%s/..", workingDir)

	command := "ssh -i ./test/tmp/id_rsa git@localhost -p 2222 \"git -C /srv/git/" + repositoryFolder + ".git log --pretty='format:%cn <%ce> %s'\""

	output, err := shell.Run(gibackRootPath, command, nil, shell.RunOptionsDefault())

	if err != nil {
		return []string{}
	}

	outputLines := outputToLines(output)

	return outputLines
}

func outputToLines(output []byte) []string {
	return strings.Split(string(output), "\n")
}

func cleanupLogMessages(lines []string) []string {
	linesTransformed := make([]string, len(lines))

	for i := range lines {
		split := strings.Split(lines[i], " ")

		if len(split) < 3 {
			continue
		}

		linesTransformed[i] = strings.Join(split[2:], " ")
	}

	return linesTransformed
}

func cleanupEmptyLines(lines []string) []string {
	var linesTransformed []string

	for i := range lines {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		linesTransformed = append(linesTransformed, lines[i])
	}

	return linesTransformed
}

func resetTestRepository(workingDir string) {
	var output []byte

	commands := []string{
		"ssh -i ./test/tmp/id_rsa git@localhost -p 2222 \"cd /srv/git/test.git && rm -rf ./* && git init --bare\"",
		"ssh -i ./test/tmp/id_rsa git@localhost -p 2222 \"cd /srv/git/test2.git && rm -rf ./* && git init --bare\"",
	}

	err := shell.RunMany(workingDir, commands, nil, &output, shell.RunOptionsDefault())

	if err != nil {
		log.Println(fmt.Sprintf("An error occurred while trying to reset the test repository:\n%s", output))

		log.Fatal(err)
	}
}

func emptyWorkspace(workingDir string) {
	var output []byte

	commands := []string{
		fmt.Sprintf("rm -rf %s/test/tmp/workspace/my_backup", workingDir),
		fmt.Sprintf("rm -rf %s/test/tmp/workspace/another_backup", workingDir),
	}

	err := shell.RunMany(workingDir, commands, nil, &output, shell.RunOptionsDefault())

	if err != nil {
		log.Println(fmt.Sprintf("An error occurred while trying to empty the test workspaces:\n%s", output))

		log.Fatal(err)
	}
}
