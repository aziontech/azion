package config

import (
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/messages/root"
)

const (
	DEFAULT_DIR      = ".azion"
	DEFAULT_SETTINGS = "settings.toml"
	DEFAULT_METRICS  = "metrics.json"
)

var (
	pathDir      string = DEFAULT_DIR
	pathSettings string = DEFAULT_SETTINGS
)

type Config interface {
	GetString(key string) string
}

func SetPath(path string) error {
	if filepath.Ext(path) != ".toml" {
		return root.ErrorReadFileSettingsToml
	}

	pathDir = filepath.Dir(path)
	pathSettings = filepath.Base(path)

	return nil
}

func GetPath() string {
	return pathDir
}

type DirPath struct {
	Dir      string
	Settings string
	Metrics  string
}

func Dir() (DirPath, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return DirPath{}, err
	}

	if pathDir != DEFAULT_DIR {
		home = ""
	}

	dirPath := DirPath{
		Dir:      filepath.Join(home, pathDir),
		Settings: pathSettings,
		Metrics:  DEFAULT_METRICS,
	}
	return dirPath, nil
}
