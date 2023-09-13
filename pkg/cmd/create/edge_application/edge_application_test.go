package edge_application

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/edge_application"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
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
			output: "Created edge application with ID 1694434702\n",
		},
		{
			name: "Creating the edge application with the --in flag",
			args: []string{"--in", "./fixtures/body_request.json"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.JSONFromFile("./fixtures/response.json"),
				)
				return &mock
			},
			output: "Created edge application with ID 1694434702\n",
		},
		{
			name: "Not found",
			args: []string{"--name", "sakura"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "edge_applications"),
					httpmock.StatusStringResponse(404, "Not Found"),
				)
				return &mock
			},
			err: fmt.Errorf(msg.ErrorCreate.Error(), "The given web page URL or API's endpoint doesn't exist or isn't available. Check that the identifying information is correct. If the error persists, contact Azion's support"),
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
