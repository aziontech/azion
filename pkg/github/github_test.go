package github

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/go-git/go-git/v5"
	"go.uber.org/zap/zapcore"
)

var cloneOptions git.CloneOptions

func TestGetVersionGitHub(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name       string
		statusCode int
		latest     string
		releases   string
		wantTag    string
		wantErr    bool
	}{
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			latest:     `{"message": "Not Found"}`,
			wantTag:    "",
		},
		{
			name:       "latest is a production release",
			statusCode: http.StatusOK,
			latest:     `{"tag_name": "4.24.0", "published_at": "2026-10-16T18:12:39Z"}`,
			wantTag:    "4.24.0",
		},
		{
			name:       "latest is a stage tag published as a release",
			statusCode: http.StatusOK,
			latest:     `{"tag_name": "4.24.0-dev.20", "published_at": "2026-10-16T18:12:39Z"}`,
			releases:   `[{"tag_name": "4.24.0-dev.20"}, {"tag_name": "4.23.1", "prerelease": true}, {"tag_name": "4.23.0", "published_at": "2026-08-26T18:12:07Z"}]`,
			wantTag:    "4.23.0",
		},
		{
			name:       "latest flagged as pre-release and no production release",
			statusCode: http.StatusOK,
			latest:     `{"tag_name": "4.24.0", "prerelease": true}`,
			releases:   `[{"tag_name": "4.24.0", "prerelease": true}, {"tag_name": "4.23.0", "draft": true}]`,
			wantTag:    "",
		},
		{
			name:       "invalid body",
			statusCode: http.StatusOK,
			latest:     `{`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body := tt.latest
				if r.URL.Path == "/repos/aziontech/azion/releases" {
					body = tt.releases
				}
				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(body)); err != nil {
					panic(err)
				}
			}))
			defer server.Close()

			previous := githubAPI
			githubAPI = server.URL
			defer func() { githubAPI = previous }()

			gotTag, _, err := NewGithub().GetVersionGitHub("azion")
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

			if err := gh.Clone(&cloneOptions, tt.url, path); (err != nil) != tt.wantErr {
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
				err := os.WriteFile(file, []byte(".edge/\n.vulcan\n.open-next\n"), 0644)
				if err != nil {
					panic(err)
				}
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
			wantFile: "#Paths added by Azion CLI\n.edge/\n.vulcan\n.open-next\n",
			cleanup: func(path string) {
				os.RemoveAll(path)
			},
		},
		{
			name: "error opening file",
			setup: func() string {
				return "does/not/exist"
			},
			wantErr: true,
			cleanup: func(path string) {
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
