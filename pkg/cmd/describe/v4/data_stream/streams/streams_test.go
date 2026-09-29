package streams

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
    "name": "My Stream",
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "created": "2019-08-24T14:15:22Z",
    "active": true,
    "product_version": "1.0",
    "inputs": [],
    "transform": [],
    "outputs": []
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
			name:      "describe a stream",
			request:   httpmock.REST("GET", "workspace/stream/streams/1337"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--stream-id", "1337"},
			expectErr: false,
		},
		{
			name:      "describe a stream - no stream id",
			request:   httpmock.REST("GET", "workspace/stream/streams/1337"),
			response:  httpmock.JSONFromString(successResponse),
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1337", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/stream/streams/9999"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--stream-id", "9999"},
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
