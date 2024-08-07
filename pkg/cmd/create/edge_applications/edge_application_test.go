package edge_applications

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/edge_application"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name   string
		args   []string
		mock   func() *httpmock.Registry
		output string
		err    error
	}{
		{
			name: "Create edge application with success",
			args: []string{"--name", "lulu"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.JSONFromFile("./fixtures/response.json"),
				)
				return &mock
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1694434702),
		},
		{
			name: "Create edge application with success all fields completed",
			args: []string{
				"--name", "lulu",
				"--http3", "true",
				"--debug-rules", "false",
				"--supported-ciphers", "AES256-GCM-SHA384",
				"--websocket", "true",
				"--http-port", "80",
				"--https-port", "443",
				"--delivery-protocol", "http",
				"--address", "https://example.com",
				"--origin-protocol-policy", "https-only",
				"--host-header", "example.com",
				"--browser-cache-settings", "max-age=3600",
				"--cdn-cache-settings", "max-age=3600",
			},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.JSONFromFile("./fixtures/response.json"),
				)
				return &mock
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1694434702),
		},
		{
			name: "Creating the edge application with the --file flag",
			args: []string{"--file", "./fixtures/body_request.json"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.JSONFromFile("./fixtures/response.json"),
				)
				return &mock
			},
			output: fmt.Sprintf(msg.OutputSuccess, 1694434702),
		},
		{
			name: "Error file json no exist",
			args: []string{"--file", "./fixtures/no_exist.json"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.JSONFromString("{}"),
				)
				return &mock
			},
			err: utils.ErrorUnmarshalReader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, out, _ := testutils.NewFactory(tt.mock())

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if err != nil && !(err.Error() == tt.err.Error()) {
				t.Errorf("Executec() err = %v, \nexpected %v", err, tt.args)
			}

			assert.Equal(t, tt.output, out.String())
		})
	}
}
