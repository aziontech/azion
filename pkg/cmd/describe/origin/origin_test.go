package origin

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
			name:      "describe an origin",
			request:   httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			response:  httpmock.JSONFromFile("./fixtures/origins.json"),
			args:      []string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"},
			expectErr: false,
		},
		{
			name:      "describe an origin - no app id",
			request:   httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			response:  httpmock.JSONFromFile("./fixtures/origins.json"),
			args:      []string{"--origin-key", "0000000-00000000-00a0a00s0as0-000000"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "123423424", nil
			},
		},
		{
			name:      "describe an origin - no origin key",
			request:   httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			response:  httpmock.JSONFromFile("./fixtures/origins.json"),
			args:      []string{"--application-id", "123423424"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "0000000-00000000-00a0a00s0as0-000000", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"},
			expectErr: true,
		},
		{
			name:      "no id sent",
			request:   httpmock.REST("GET", "edge_applications/123423424/origins/0000000-00000000-00a0a00s0as0-000000"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--application-id", "123423424", "--origin-key", "0000000-00000000-00a0a00s0as0-000000"},
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
