package app

import (
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

func folderExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}
