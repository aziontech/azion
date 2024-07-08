package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func TestGetVersionGitHub(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name       string
		repoName   string
		wantTag    string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "not found",
			repoName:   "test-repo",
			wantTag:    "",
			statusCode: http.StatusNotFound,
			response:   `{"message": "Not Found"}`,
			wantErr:    false,
		},
		{
			name:       "successful response",
			repoName:   "azion-cli",
			wantTag:    "1.30.0",
			statusCode: http.StatusOK,
			response:   `{"tag_name": "1.30.0"}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			oldURL := ApiURL
			ApiURL = server.URL + "/repos/aziontech/%s/releases/tag/1.30.0"
			defer func() { ApiURL = oldURL }()

			gh := NewGithub()
			gh.GetVersionGitHub = getVersionGitHub

			gotTag, err := gh.GetVersionGitHub(tt.repoName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersionGitHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTag != tt.wantTag {
				t.Errorf("GetVersionGitHub() = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}

func TestClone(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name    string
		url     string
		wantErr bool
		setup   func() (string, func())
	}{
		{
			name:    "invalid url",
			url:     "invalid-url",
			wantErr: true,
			setup: func() (string, func()) {
				return t.TempDir(), func() {}
			},
		},
		{
			name: "valid url",
			url:  "https://github.com/aziontech/azion-cli.git",
			setup: func() (string, func()) {
				// Create a temporary directory
				dir := t.TempDir()
				// Return the directory path and a cleanup function
				return dir, func() {
					// Cleanup steps if needed
					os.RemoveAll(dir)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setup()
			defer cleanup()

			gh := NewGithub()
			gh.Clone = clone

			if err := gh.Clone(tt.url, path); (err != nil) != tt.wantErr {
				t.Errorf("Clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNameRepo(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "with .git",
			url:  "https://github.com/aziontech/azion-cli.git",
			want: "azion-cli",
		},
		{
			name: "without .git",
			url:  "https://github.com/aziontech/azion-cli",
			want: "azion-cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh := NewGithub()
			gh.GetNameRepo = getNameRepo

			if got := gh.GetNameRepo(tt.url); got != tt.want {
				t.Errorf("GetNameRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckGitignore(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name    string
		setup   func() string
		want    bool
		wantErr bool
		cleanup func(string)
	}{
		{
			name: "exists",
			setup: func() string {
				path := t.TempDir()
				file := filepath.Join(path, ".gitignore")
				os.WriteFile(file, []byte(".edge/\n.vulcan\n"), 0644)
				return path
			},
			want:    true,
			wantErr: false,
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
		},
		{
			name: "does not exist",
			setup: func() string {
				return t.TempDir()
			},
			want:    false,
			wantErr: false,
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			gh := NewGithub()
			gh.CheckGitignore = checkGitignore

			got, err := gh.CheckGitignore(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckGitignore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckGitignore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteGitignore(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name     string
		setup    func() string
		wantErr  bool
		wantFile string
		cleanup  func(string)
	}{
		{
			name: "success",
			setup: func() string {
				return t.TempDir()
			},
			wantErr:  false,
			wantFile: "#Paths added by Azion CLI\n.edge/\n.vulcan\n",
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
		},
		{
			name: "error opening file",
			setup: func() string {
				path := t.TempDir()
				file := filepath.Join(path, ".gitignore")
				os.WriteFile(file, []byte{}, 0444) // Create a read-only file
				return path
			},
			wantErr: true,
			cleanup: func(path string) {
				os.Chmod(filepath.Join(path, ".gitignore"), 0644) // Restore permissions before cleanup
				os.RemoveAll(path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			defer tt.cleanup(path)

			gh := NewGithub()
			gh.WriteGitignore = writeGitignore

			err := gh.WriteGitignore(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteGitignore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				got, err := os.ReadFile(filepath.Join(path, ".gitignore"))
				if err != nil {
					t.Errorf("Error reading .gitignore file: %v", err)
					return
				}
				if string(got) != tt.wantFile {
					t.Errorf("WriteGitignore() = %v, want %v", string(got), tt.wantFile)
				}
			}
		})
	}
}
