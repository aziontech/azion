package cachesetting

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name      string
		request   httpmock.Matcher
		response  httpmock.Responder
		args      []string
		expectErr bool
		mockInput func(string) (string, error)
	}{
		{
			name:      "describe a cache setting",
			request:   httpmock.REST("GET", "workspace/applications/1673635839/cache_settings/107313"),
			response:  httpmock.JSONFromFile("./fixtures/cache_settings.json"),
			args:      []string{"--application-id", "1673635839", "--cache-setting-id", "107313"},
			expectErr: false,
		},
		{
			name:      "describe a cache setting - no app id",
			request:   httpmock.REST("GET", "workspace/applications/1673635839/cache_settings/107313"),
			response:  httpmock.JSONFromFile("./fixtures/cache_settings.json"),
			args:      []string{"--cache-setting-id", "107313"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1673635839", nil
			},
		},
		{
			name:      "describe a cache setting - no cache id",
			request:   httpmock.REST("GET", "workspace/applications/1673635839/cache_settings/107313"),
			response:  httpmock.JSONFromFile("./fixtures/cache_settings.json"),
			args:      []string{"--application-id", "1673635839"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "107313", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/applications/1673635839/cache_settings/107313"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--application-id", "1673635839", "--cache-setting-id", "107313"},
			expectErr: true,
		},
		{
			name:      "no id sent",
			request:   httpmock.REST("GET", "workspace/applications/1673635839/cache_settings/0"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--application-id", "1673635839", "--cache-setting-id", "0"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, _, _ := testutils.NewFactory(mock)
			descCmd := NewDescribeCmd(f)
			descCmd.AskInput = tt.mockInput
			cmd := NewCobraCmd(descCmd, f)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
