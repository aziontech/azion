package config

import (
	"os"
	"path/filepath"
)

const defaultPath = ".azion"

var pathDir string = defaultPath

type Config interface {
	GetString(key string) string
}

func SetPath(cp string) {
	pathDir = cp
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if pathDir != defaultPath {
		home = ""
	}

	return filepath.Join(home, pathDir), nil
}
