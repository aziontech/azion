package applications

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/application"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name    string
		args    []string
		mock    func() *httpmock.Registry
		payload map[string]interface{}
		output  string
		err     error
	}{
		{
			name: "create with flags",
			args: []string{"--name", "my-app", "--cache", "true", "--function", "false", "--active", "true", "--debug-enabled", "false"},
			payload: map[string]interface{}{
				"name":   "my-app",
				"active": true,
				"debug":  false,
				"modules": map[string]interface{}{
					"cache":     map[string]interface{}{"enabled": true},
					"functions": map[string]interface{}{"enabled": false},
				},
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1694434702),
		},
		{
			name: "create with file",
			args: []string{"--file", "./fixtures/create.json"},
			payload: map[string]interface{}{
				"name":   "from-file",
				"active": true,
				"modules": map[string]interface{}{
					"cache": map[string]interface{}{"enabled": true},
				},
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1694434702),
		},
		{
			name: "file does not exist",
			args: []string{"--file", "./fixtures/no_exist.json"},
			err:  utils.ErrorUnmarshalReader,
		},
		{
			name: "file without the required name",
			args: []string{"--file", "./fixtures/no_name.json"},
			err:  utils.ErrorUnmarshalReader,
		},
		{
			name: "invalid boolean flag",
			args: []string{"--name", "my-app", "--active", "maybe"},
			err:  utils.ErrorConvertingStringToBool,
		},
		{
			name: "api error",
			args: []string{"--name", "my-app"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "workspace/applications"),
					httpmock.StatusStringResponse(http.StatusConflict, `{"errors": [{"detail": "name already in use"}]}`),
				)
				return &mock
			},
			err: fmt.Errorf(msg.ErrorCreate.Error(), utils.ErrorNameInUse),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]interface{}
			mock := &httpmock.Registry{}
			if tt.mock != nil {
				mock = tt.mock()
			} else {
				mock.Register(
					httpmock.REST("POST", "workspace/applications"),
					httpmock.WithHeader(
						httpmock.RESTPayload(http.StatusOK, mustRead(t, "./fixtures/response.json"), func(p map[string]interface{}) { got = p }),
						"Content-Type", "application/json",
					),
				)
			}

			f, out, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if tt.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.output, out.String())
			assert.Equal(t, tt.payload, got)
		})
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}
