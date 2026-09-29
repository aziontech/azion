package networklist

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "My Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`

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
			name:      "describe a network list",
			request:   httpmock.REST("GET", "workspace/network_lists/1337"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--network-list-id", "1337"},
			expectErr: false,
		},
		{
			name:      "describe a network list - no network list id",
			request:   httpmock.REST("GET", "workspace/network_lists/1337"),
			response:  httpmock.JSONFromString(successResponse),
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1337", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/network_lists/9999"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--network-list-id", "9999"},
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
