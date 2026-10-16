package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestRemoteConfigFile(t *testing.T) {
	workDir := t.TempDir()

	tests := []struct {
		name       string
		configFile string
		want       string
		wantErr    error
	}{
		{name: "not sent uses the default lookup", configFile: "", want: ""},
		{name: "absolute path inside the project", configFile: filepath.Join(workDir, "azion.staging.config.ts"), want: "azion.staging.config.ts"},
		{name: "nested path inside the project", configFile: filepath.Join(workDir, "configs", "azion.config.ts"), want: "configs/azion.config.ts"},
		{name: "file whose name starts with dots", configFile: filepath.Join(workDir, "..azion.config.ts"), want: "..azion.config.ts"},
		{name: "outside the project is rejected", configFile: filepath.Join(filepath.Dir(workDir), "azion.config.ts"), wantErr: msg.ErrorConfigFileOutsideProject},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := remoteConfigFile(workDir, tt.configFile)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("remoteConfigFile() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("remoteConfigFile() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("remoteConfigFile() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCallScriptConfigFile(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name       string
		configFile string
	}{
		{name: "sends the config file", configFile: "configs/azion.staging.config.ts"},
		{name: "sends empty by default", configFile: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload []map[string]string
			mock := &httpmock.Registry{}
			mock.Register(httpmock.MatchAny, func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				if err := json.Unmarshal(body, &payload); err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Request:    req,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(bytes.NewBufferString(`{"uuid": "12345"}`)),
				}, nil
			})

			f, _, _ := testutils.NewFactory(mock)
			cmd := NewDeployCmd(f)
			cmd.ConfigFile = tt.configFile
			if _, err := cmd.CallScript("token", "id", "secret", "prefix", "project", "azion", cmd); err != nil {
				t.Fatalf("CallScript() error = %v", err)
			}

			for _, env := range payload {
				if env["field"] == "AZCLI_CONFIG_FILE" {
					if env["instantiation_data_path"] != "envs.[6].value" {
						t.Errorf("instantiation_data_path = %q, want envs.[6].value", env["instantiation_data_path"])
					}
					if env["value"] != tt.configFile {
						t.Errorf("value = %q, want %q", env["value"], tt.configFile)
					}
					return
				}
			}
			t.Fatalf("AZCLI_CONFIG_FILE not sent in payload %v", payload)
		})
	}
}
