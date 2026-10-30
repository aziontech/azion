package application

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/update/application"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name    string
		args    []string
		mock    func(t *testing.T, got *map[string]interface{}) *httpmock.Registry
		payload map[string]interface{}
		output  string
		err     error
	}{
		{
			name: "update with flags",
			args: []string{"--application-id", "1337", "--name", "renamed", "--cache", "false", "--image-processor", "true"},
			mock: patchOK("workspace/applications/1337"),
			payload: map[string]interface{}{
				"name": "renamed",
				"modules": map[string]interface{}{
					"cache":           map[string]interface{}{"enabled": false},
					"image_processor": map[string]interface{}{"enabled": true},
				},
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1337),
		},
		{
			name: "update with file uses the id from the flag",
			args: []string{"--application-id", "1337", "--file", "./fixtures/update.json"},
			mock: patchOK("workspace/applications/1337"),
			payload: map[string]interface{}{
				"name":   "renamed",
				"active": true,
				"debug":  false,
				"modules": map[string]interface{}{
					"cache": map[string]interface{}{"enabled": true},
				},
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1337),
		},
		{
			name: "file does not exist",
			args: []string{"--application-id", "1337", "--file", "./fixtures/no_exist.json"},
			err:  utils.ErrorUnmarshalReader,
		},
		{
			name: "invalid boolean flag",
			args: []string{"--application-id", "1337", "--active", "unactive"},
			err:  utils.ErrorConvertingStringToBool,
		},
		{
			name: "no field informed",
			args: []string{"--application-id", "1337"},
			err:  msg.ErrorNoFieldInformed,
		},
		{
			name: "not found",
			args: []string{"--application-id", "1337", "--name", "renamed"},
			mock: func(t *testing.T, _ *map[string]interface{}) *httpmock.Registry {
				mock := &httpmock.Registry{}
				mock.Register(
					httpmock.REST("PATCH", "workspace/applications/1337"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
				return mock
			},
			err: utils.ErrorNotFound404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]interface{}
			mock := &httpmock.Registry{}
			if tt.mock != nil {
				mock = tt.mock(t, &got)
			}

			f, out, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

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

func patchOK(path string) func(t *testing.T, got *map[string]interface{}) *httpmock.Registry {
	return func(t *testing.T, got *map[string]interface{}) *httpmock.Registry {
		body, err := os.ReadFile("./fixtures/response.json")
		require.NoError(t, err)

		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", path),
			httpmock.WithHeader(
				httpmock.RESTPayload(http.StatusOK, string(body), func(p map[string]interface{}) { *got = p }),
				"Content-Type", "application/json",
			),
		)
		return mock
	}
}
