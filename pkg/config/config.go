package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/aziontech/azion-cli/messages/root"
)

const (
	DEFAULT_DIR      = ".azion"
	DEFAULT_SETTINGS = "settings.toml"
	DEFAULT_METRICS  = "metrics.json"
	DEFAULT_SCHEDULE = "schedule.json"
	DEFAULT_PROFILES = "profiles.json"
)

var (
	pathDir      string = DEFAULT_DIR
	pathSettings string = DEFAULT_SETTINGS
	pathProfiles string = DEFAULT_PROFILES
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
	pathProfiles = DEFAULT_PROFILES

	return nil
}

func GetPath() string {
	return pathDir
}

type DirPath struct {
	Dir      string
	Settings string
	Metrics  string
	Schedule string
	Profiles string
}

func Dir() DirPath {
	home := userHomeDir()

	if pathDir != DEFAULT_DIR {
		home = ""
	}

	dirPath := DirPath{
		Dir:      filepath.Join(home, pathDir),
		Settings: pathSettings,
		Profiles: pathProfiles,
		Metrics:  DEFAULT_METRICS,
		Schedule: DEFAULT_SCHEDULE,
	}
	return dirPath
}

func userHomeDir() string {
	env := "HOME"
	switch runtime.GOOS {
	case "windows":
		env = "USERPROFILE"
	case "plan9":
		env = "home"
	}
	if v := os.Getenv(env); v != "" {
		return v
	}
	// On some geese the home directory is not always defined.
	switch runtime.GOOS {
	case "android":
		return "/sdcard"
	case "ios":
		return "/"
	}
	return ""
}
