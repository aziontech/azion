package config

import (
	"os"
	"path/filepath"
)

type Config interface {
	GetString(key string) string
}

var pathDir string = ".azion"

func SetPath(cp string) {
	pathDir = cp
}

func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if len(pathDir) > 0 {
		home = ""
	}

	return filepath.Join(home, pathDir), nil
}
