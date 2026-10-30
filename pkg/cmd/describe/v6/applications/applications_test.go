package applications

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/describe/applications"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		request  httpmock.Matcher
		response httpmock.Responder
		args     []string
		format   string
		askInput func(string) (string, error)
		contains []string
		err      error
	}{
		{
			name:     "describe an application",
			request:  httpmock.REST("GET", "workspace/applications/1232132135"),
			response: httpmock.JSONFromFile("./fixtures/response.json"),
			args:     []string{"--application-id", "1232132135"},
			contains: []string{"1232132135", "my-app"},
		},
		{
			name:     "describe asking for the id",
			request:  httpmock.REST("GET", "workspace/applications/1232132135"),
			response: httpmock.JSONFromFile("./fixtures/response.json"),
			askInput: func(string) (string, error) { return "1232132135", nil },
			contains: []string{"1232132135", "my-app"},
		},
		{
			name:     "describe with json format",
			request:  httpmock.REST("GET", "workspace/applications/1232132135"),
			response: httpmock.JSONFromFile("./fixtures/response.json"),
			args:     []string{"--application-id", "1232132135"},
			format:   "json",
			contains: []string{`"version_state": "ready"`, `"product_version": "1.0"`},
		},
		{
			name:     "not found",
			request:  httpmock.REST("GET", "workspace/applications/1234"),
			response: httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:     []string{"--application-id", "1234"},
			err:      utils.ErrorNotFound404,
		},
		{
			name:     "invalid id answered",
			request:  httpmock.REST("GET", "workspace/applications/1234"),
			response: httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			askInput: func(string) (string, error) { return "abc", nil },
			err:      msg.ErrorConvertIdApplication,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, out, _ := testutils.NewFactory(mock)
			f.Format = tt.format
			describe := NewDescribeCmd(f)
			if tt.askInput != nil {
				describe.AskInput = tt.askInput
			}
			cmd := NewCobraCmd(describe, f)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			for _, s := range tt.contains {
				assert.Contains(t, out.String(), s)
			}
		})
	}
}
