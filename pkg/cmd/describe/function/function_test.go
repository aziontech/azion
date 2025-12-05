package function

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
    "name": "string",
    "active": true,
    "runtime": "azion_js",
    "reference_count": 0,
    "last_modified": "2019-08-24T14:15:22Z",
    "product_version": "v1",
    "version": "1.0.0",
    "vendor": "azion",
    "execution_environment": "firewall",
    "last_editor": "string",
    "default_args": {
      "arg_01": "value_01"
    },
    "code": "console.log('hello');"
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
			name:      "describe a function",
			request:   httpmock.REST("GET", "workspace/functions/1337"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--function-id", "1337"},
			expectErr: false,
		},
		{
			name:      "describe a function - no function id",
			request:   httpmock.REST("GET", "workspace/functions/1337"),
			response:  httpmock.JSONFromString(successResponse),
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1337", nil
			},
		},
		{
			name:      "with code",
			request:   httpmock.REST("GET", "workspace/functions/1337"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--function-id", "1337", "--with-code"},
			expectErr: false,
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/functions/1234"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--function-id", "1234", "--with-code"},
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
