package gibackfs

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/dhuan/giback/pkg/app"
	"gopkg.in/yaml.v2"
)

func GetUserConfig() (app.Config, string, error) {
	config := app.Config{}

	userConfigFilePath, workspacePath, err := getUserConfigFilePath()

	if err != nil {
		return config, workspacePath, err
	}

	buffer, err := ioutil.ReadFile(userConfigFilePath)

	if err != nil {
		return config, workspacePath, err
	}

	fileContent := string(buffer)

	yamlErr := yaml.Unmarshal([]byte(fileContent), &config)

	if yamlErr != nil {
		return config, workspacePath, yamlErr
	}

	return config, workspacePath, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

func FolderExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

func ScanDir(path string) []string {
	var files []string
	var filesFiltered []string

	basePath := getGlobBasePath(path)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)

		return nil
	})

	for _, file := range files {
		match, _ := filepath.Match(path, file)

		if match {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return filesFiltered
}

func ScanDirMany(pathList []string) []string {
	var files []string

	if len(pathList) == 0 {
		return nil
	}

	files = ScanDir(pathList[0])

	files = append(files, ScanDirMany(pathList[1:])...)

	return files
}

func getGlobBasePath(pathGlob string) string {
	split := strings.Split(pathGlob, "/")

	i := len(split) - 1

	return strings.Join(split[0:i], "/")
}

func Copy(files []string, copyTo string) error {
	if len(files) == 0 {
		return nil
	}

	file := files[0]

	input, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	fileName := getFileNameFromFullPathFile(file)

	destinationFile := copyTo + "/" + fileName

	err = ioutil.WriteFile(destinationFile, input, 0644)

	if err != nil {
		return err
	}

	return Copy(files[1:], copyTo)
}

func getFileNameFromFullPathFile(filePath string) string {
	split := strings.Split(filePath, "/")

	return split[len(split)-1]
}

func getUserConfigFilePath() (string, string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", "", err
	}

	configFilePath := fmt.Sprintf("%s/.giback.yml", usr.HomeDir)

	workspacePath := fmt.Sprintf("%s/.giback", usr.HomeDir)

	return configFilePath, workspacePath, nil
}
