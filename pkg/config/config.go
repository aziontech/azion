package config

import (
	"os"
	"path/filepath"
)

const DEFAULT_PATH = ".azion"

var pathDir string = DEFAULT_PATH 

type Config interface {
	GetString(key string) string
}

func SetPath(cp string) {
	pathDir = cp
}

func GetPath() string {
	return pathDir
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if pathDir != DEFAULT_PATH {
		home = ""
	}

	return filepath.Join(home, pathDir), nil
}
