package testutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dhuan/giback/pkg/shell"
)

func defaultBackupFiles() []string {
	return []string{
		"some_file.txt This file shall be backed up.",
		"another_file.txt This file shall be backed up.",
		"unimportant_file.txt This file shall not be commited.",
	}
}

func ResetTestEnvironment() {
	workingDir, _ := os.Getwd()

	gibackRootPath := fmt.Sprintf("%s/..", workingDir)

	resetTestRepository(gibackRootPath)

	emptyWorkspace(gibackRootPath)

	emptyBackupFiles(gibackRootPath)

	AddFileToBackupFilesFolder(defaultBackupFiles())
}

func RunGiback(command string, options RunGibackOptions) ([]byte, error) {
	workspacePath := "./test/tmp/workspace"
	configPath := "./test/config/default.yml"
	workingDir, _ := os.Getwd()
	gibackRootPath := fmt.Sprintf("%s/..", workingDir)

	if options.ConfigFile != "" {
		configPath = fmt.Sprintf("./test/config/%s.yml", options.ConfigFile)
	}

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

func AssertHasNoError(t *testing.T, err error) {
	if err != nil {
		t.Error("Expected not to fail.")
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

func AddFileToWorkspace(repositoryId string, fileDefinitions []string) {
	if len(fileDefinitions) == 0 {
		return
	}

	workingDir, _ := os.Getwd()
	gibackRootPath := fmt.Sprintf("%s/..", workingDir)
	repositoryDir := fmt.Sprintf("%s/test/tmp/workspace/%s", gibackRootPath, repositoryId)

	fileName, fileContent := parseFileDefinition(fileDefinitions[0])

	fileNameFull := fmt.Sprintf("%s/%s", repositoryDir, fileName)

	err := ioutil.WriteFile(fileNameFull, []byte(fileContent), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	AddFileToWorkspace(repositoryId, fileDefinitions[1:])
}

func AddFileToBackupFilesFolder(fileDefinitions []string) {
	if len(fileDefinitions) == 0 {
		return
	}

	workingDir, _ := os.Getwd()
	gibackRootPath := fmt.Sprintf("%s/..", workingDir)
	backupDir := fmt.Sprintf("%s/test/backup_files", gibackRootPath)

	fileName, fileContent := parseFileDefinition(fileDefinitions[0])

	fileNameFull := fmt.Sprintf("%s/%s", backupDir, fileName)

	err := ioutil.WriteFile(fileNameFull, []byte(fileContent), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	AddFileToBackupFilesFolder(fileDefinitions[1:])
}

type RunGibackOptions struct {
	ConfigFile string
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
		containsDate := regexMatches("^[0-9]{4}\\/", lines[i])

		if !containsDate {
			linesTransformed[i] = lines[i]

			continue
		}

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

func emptyBackupFiles(gibackRootPath string) {
	basePath := fmt.Sprintf("%s/test/backup_files", gibackRootPath)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if strings.Index(path, "gitkeep") > -1 {
			return nil
		}

		if strings.Index(path, ".txt") == -1 && strings.Index(path, ".html") == -1 {
			return nil
		}

		err = os.Remove(path)

		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed while removing file to empty backup files: %s", path))
		}

		return nil
	})
}

func regexMatches(pattern string, subject string) bool {
	matched, err := regexp.MatchString(pattern, subject)

	if err != nil {
		return false
	}

	return matched
}

func parseFileDefinition(fileDefinition string) (string, string) {
	splitResult := strings.Split(fileDefinition, " ")

	if len(splitResult) < 2 {
		log.Fatalln(fmt.Sprintf("File definition could not be parsed: %s", fileDefinition))
	}

	fileName := splitResult[0]
	fileContent := splitResult[1:]

	return fileName, strings.Join(fileContent, " ")
}
