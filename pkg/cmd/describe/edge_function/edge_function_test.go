package edgefunction

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
    "results":{
        "id":1337,
        "name":"SUUPA_FUNCTION",
        "language":"javascript",
        "code":"async function handleRequest(request) {return new Response(\"Hello World!\",{status:200})}",
        "json_args":{"a":1,"b":2},"function_to_run":"",
        "initiator_type":"edge_application",
        "active":true,
        "last_editor":"testando@azion.com",
        "modified":"2022-01-26T12:31:09.865515Z",
        "reference_count":0
    },
    "schema_version":3
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
	}{
		{
			name:      "describe a function",
			request:   httpmock.REST("GET", "edge_functions/123"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--function-id", "123"},
			expectErr: false,
		},
		{
			name:      "with code",
			request:   httpmock.REST("GET", "edge_functions/123"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--function-id", "123", "--with-code"},
			expectErr: false,
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "edge_functions/1234"),
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
