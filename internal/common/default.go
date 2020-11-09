package common

import (
	"os"
	"path/filepath"
)

const (
	ServiceSpecFile = "service.yml"
)

func IsServiceDir(path string) bool {
	specFile, err := os.Stat(filepath.Join(path, ServiceSpecFile))
	if err != nil {
		return false
	}
	if specFile.IsDir() {
		return false
	}
	return true
}

func FileExists(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !file.IsDir()
}
