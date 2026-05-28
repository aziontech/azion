package csr

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/delete/csr"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var deleteResponse string = `{"state": "executed"}`

func mockCSRID(string) (string, error) {
	return "1234", nil
}

func mockParseErrorCSRID(string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

// CSR delete proxies the standard certificate endpoint.
func TestDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		csrID          string
		request        httpmock.Matcher
		response       httpmock.Responder
		expectedOutput string
		expectError    bool
		mockInput      func(string) (string, error)
	}{
		{
			name:           "delete csr by id",
			csrID:          "1234",
			request:        httpmock.REST("DELETE", "workspace/tls/certificates/1234"),
			response:       httpmock.JSONFromString(deleteResponse),
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
		},
		{
			name:           "ask for csr id success",
			csrID:          "",
			request:        httpmock.REST("DELETE", "workspace/tls/certificates/1234"),
			response:       httpmock.JSONFromString(deleteResponse),
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInput:      mockCSRID,
		},
		{
			name:        "delete csr - not found",
			csrID:       "1234",
			request:     httpmock.REST("DELETE", "workspace/tls/certificates/1234"),
			response:    httpmock.StatusStringResponse(404, "Not Found"),
			expectError: true,
		},
		{
			name:        "error - parse answer",
			csrID:       "",
			request:     httpmock.REST("DELETE", "workspace/tls/certificates/1234"),
			response:    httpmock.JSONFromString(deleteResponse),
			expectError: true,
			mockInput:   mockParseErrorCSRID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, stdout, _ := testutils.NewFactory(mock)

			deleteCmd := NewDeleteCmd(f)
			if tt.mockInput != nil {
				deleteCmd.AskInput = tt.mockInput
			}
			cobraCmd := NewCobraCmd(deleteCmd, f)

			if tt.csrID != "" {
				cobraCmd.SetArgs([]string{"--csr-id", tt.csrID})
			}

			_, err := cobraCmd.ExecuteC()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, stdout.String())
			}
		})
	}
}
