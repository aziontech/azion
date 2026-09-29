package firewallinstance

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
			name:      "describe a firewall function instance",
			request:   httpmock.REST("GET", "workspace/firewalls/1234/functions/5678"),
			response:  httpmock.JSONFromFile("./fixtures/response.json"),
			args:      []string{"--firewall-id", "1234", "--instance-id", "5678"},
			expectErr: false,
		},
		{
			name:      "describe a firewall function instance - no firewall id",
			request:   httpmock.REST("GET", "workspace/firewalls/1234/functions/5678"),
			response:  httpmock.JSONFromFile("./fixtures/response.json"),
			args:      []string{"--instance-id", "5678"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1234", nil
			},
		},
		{
			name:      "describe a firewall function instance - no instance id",
			request:   httpmock.REST("GET", "workspace/firewalls/1234/functions/5678"),
			response:  httpmock.JSONFromFile("./fixtures/response.json"),
			args:      []string{"--firewall-id", "1234"},
			expectErr: false,
			mockInput: func(s string) (string, error) {
				if s == "Enter the Firewall's ID:" {
					return "1234", nil
				}
				return "5678", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/firewalls/1234/functions/9999"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--firewall-id", "1234", "--instance-id", "9999"},
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
