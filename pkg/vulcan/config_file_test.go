package vulcan

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigFileArg(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		want       string
	}{
		{name: "empty", configFile: "", want: ""},
		{name: "plain path", configFile: "azion.staging.config.ts", want: " --config-file 'azion.staging.config.ts'"},
		{name: "path with spaces", configFile: "my configs/azion.config.ts", want: " --config-file 'my configs/azion.config.ts'"},
		{name: "path with single quote", configFile: "it's/azion.config.ts", want: ` --config-file 'it'\''s/azion.config.ts'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConfigFileArg(tt.configFile); got != tt.want {
				t.Errorf("ConfigFileArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckConfigFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "azion.staging.config.ts")
	if err := os.WriteFile(file, []byte("export default {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		configFile string
		wantErr    bool
	}{
		{name: "empty is valid", configFile: ""},
		{name: "existing file", configFile: file},
		{name: "missing file", configFile: filepath.Join(dir, "missing.config.ts"), wantErr: true},
		{name: "directory", configFile: dir, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckConfigFile(tt.configFile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrConfigFileNotFound) {
				t.Errorf("CheckConfigFile() error = %v, want ErrConfigFileNotFound", err)
			}
		})
	}
}
