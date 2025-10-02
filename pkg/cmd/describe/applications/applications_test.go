package applications

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
	}{
		{
			name:      "describe an application",
			request:   httpmock.REST("GET", "edge_application/applications/1232132135"),
			response:  httpmock.JSONFromFile("./fixtures/response.json"),
			args:      []string{"--application-id", "1232132135"},
			expectErr: false,
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "edge_application/applications/1234"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--application-id", "1234"},
			expectErr: true,
		},
		{
			name:      "no id sent",
			request:   httpmock.REST("GET", "edge_application/applications/1234"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, _, _ := testutils.NewFactory(mock)
			cmd := NewCmd(f)
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
