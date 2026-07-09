package csr

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// CSR read proxies the standard certificate endpoint, so the response is a
// Certificate whose csr field holds the generated signing request.
var successResponse string = `
{
  "state": "executed",
  "data": {
    "id": 1337,
    "name": "My CSR",
    "issuer": "string",
    "subject_name": ["example.com"],
    "validity": "string",
    "type": "certificate",
    "managed": false,
    "status": "pending",
    "status_detail": "",
    "csr": "-----BEGIN CERTIFICATE REQUEST-----MIIBdummycsr==-----END CERTIFICATE REQUEST-----",
    "challenge": "dns",
    "authority": "lets_encrypt",
    "key_algorithm": "rsa_2048",
    "active": true,
    "product_version": "1.0",
    "last_editor": "user@example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "last_modified": "2019-08-24T14:15:22Z",
    "renewed_at": "2019-08-24T14:15:22Z"
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
			name:      "describe a csr by id",
			request:   httpmock.REST("GET", "workspace/tls/certificates/1337"),
			response:  httpmock.JSONFromString(successResponse),
			args:      []string{"--csr-id", "1337"},
			expectErr: false,
		},
		{
			name:      "describe a csr - ask for id",
			request:   httpmock.REST("GET", "workspace/tls/certificates/1337"),
			response:  httpmock.JSONFromString(successResponse),
			expectErr: false,
			mockInput: func(s string) (string, error) {
				return "1337", nil
			},
		},
		{
			name:      "not found",
			request:   httpmock.REST("GET", "workspace/tls/certificates/9999"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--csr-id", "9999"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, _, _ := testutils.NewFactory(mock)
			descCmd := NewDescribeCmd(f)
			if tt.mockInput != nil {
				descCmd.AskInput = tt.mockInput
			}
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
