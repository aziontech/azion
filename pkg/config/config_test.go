package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/messages/root"
)

func TestSetPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr error
		wantDir string
		wantSet string
	}{
		{
			name:    "valid path",
			path:    "/home/user/.azion/settings.toml",
			wantErr: nil,
			wantDir: "/home/user/.azion",
			wantSet: "settings.toml",
		},
		{
			name:    "invalid path extension",
			path:    "/home/user/.azion/settings.json",
			wantErr: root.ErrorReadFileSettingsToml,
			wantDir: DEFAULT_DIR,
			wantSet: DEFAULT_SETTINGS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to defaults before each test
			pathDir = DEFAULT_DIR
			pathSettings = DEFAULT_SETTINGS

			err := SetPath(tt.path)
			if err != tt.wantErr {
				t.Errorf("SetPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if pathDir != tt.wantDir {
				t.Errorf("SetPath() pathDir = %v, want %v", pathDir, tt.wantDir)
			}
			if pathSettings != tt.wantSet {
				t.Errorf("SetPath() pathSettings = %v, want %v", pathSettings, tt.wantSet)
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	tests := []struct {
		name     string
		setPath  string
		wantPath string
	}{
		{
			name:     "default path",
			setPath:  "",
			wantPath: DEFAULT_DIR,
		},
		{
			name:     "custom path",
			setPath:  "/home/user/.azion/settings.toml",
			wantPath: "/home/user/.azion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to defaults before each test
			pathDir = DEFAULT_DIR
			pathSettings = DEFAULT_SETTINGS

			if tt.setPath != "" {
				SetPath(tt.setPath)
			}
			if gotPath := GetPath(); gotPath != tt.wantPath {
				t.Errorf("GetPath() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		name    string
		setPath string
		envHome string
		wantDir DirPath
		wantErr bool
	}{
		{
			name:    "default path",
			setPath: "",
			envHome: home,
			wantDir: DirPath{
				Dir:      filepath.Join(home, DEFAULT_DIR),
				Settings: DEFAULT_SETTINGS,
				Metrics:  DEFAULT_METRICS,
				Schedule: DEFAULT_SCHEDULE,
			},
			wantErr: false,
		},
		{
			name:    "custom path",
			setPath: "/home/user/.azion/settings.toml",
			envHome: home,
			wantDir: DirPath{
				Dir:      "/home/user/.azion",
				Settings: "settings.toml",
				Metrics:  DEFAULT_METRICS,
				Schedule: DEFAULT_SCHEDULE,
			},
			wantErr: false,
		},
		{
			name:    "error getting home dir",
			setPath: "",
			envHome: "",
			wantDir: DirPath{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to defaults before each test
			pathDir = DEFAULT_DIR
			pathSettings = DEFAULT_SETTINGS

			// Set environment variable for HOME
			if tt.envHome != "" {
				os.Setenv("HOME", tt.envHome)
			} else {
				os.Unsetenv("HOME")
			}

			if tt.setPath != "" {
				SetPath(tt.setPath)
			}

			gotDir, err := Dir()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDir != tt.wantDir {
				t.Errorf("Dir() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}
