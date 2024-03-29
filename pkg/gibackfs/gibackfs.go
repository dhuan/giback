package gibackfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/dhuan/giback/pkg/app"
	"github.com/dhuan/giback/pkg/utils"
	"gopkg.in/yaml.v2"
)

func GetUserConfig(appContext app.Context) (app.Config, string, error) {
	config := app.Config{}

	userConfigFilePath := appContext.ConfigFilePath

	workspacePath := appContext.WorkspacePath

	buffer, err := ioutil.ReadFile(userConfigFilePath)

	if err != nil {
		return config, workspacePath, err
	}

	fileContent := string(buffer)

	yamlErr := yaml.Unmarshal([]byte(fileContent), &config)

	for i := range config.Units {
		config.Units[i].RepositoryPath = fmt.Sprintf("%s/%s", workspacePath, config.Units[i].Id)
	}

	if yamlErr != nil {
		return config, workspacePath, yamlErr
	}

	return config, workspacePath, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func FolderExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
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

func ScanDirMany(pathList []string, excludePatterns []string) []string {
	var files []string

	if len(pathList) == 0 {
		return nil
	}

	files = ScanDir(pathList[0])

	if len(excludePatterns) > 0 {
		files = utils.FilterOut(files, excludePatterns)
	}

	files = append(files, ScanDirMany(pathList[1:], excludePatterns)...)

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
