package crl

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/delete/crl"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var deleteResponse string = `{"state": "executed"}`

func mockCRLID(string) (string, error) {
	return "1234", nil
}

func mockParseErrorCRLID(string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		crlID          string
		request        httpmock.Matcher
		response       httpmock.Responder
		expectedOutput string
		expectError    bool
		mockInput      func(string) (string, error)
	}{
		{
			name:           "delete crl by id",
			crlID:          "1234",
			request:        httpmock.REST("DELETE", "workspace/tls/crls/1234"),
			response:       httpmock.JSONFromString(deleteResponse),
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
		},
		{
			name:           "ask for crl id success",
			crlID:          "",
			request:        httpmock.REST("DELETE", "workspace/tls/crls/1234"),
			response:       httpmock.JSONFromString(deleteResponse),
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInput:      mockCRLID,
		},
		{
			name:        "delete crl - not found",
			crlID:       "1234",
			request:     httpmock.REST("DELETE", "workspace/tls/crls/1234"),
			response:    httpmock.StatusStringResponse(404, "Not Found"),
			expectError: true,
		},
		{
			name:        "error - parse answer",
			crlID:       "",
			request:     httpmock.REST("DELETE", "workspace/tls/crls/1234"),
			response:    httpmock.JSONFromString(deleteResponse),
			expectError: true,
			mockInput:   mockParseErrorCRLID,
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

			if tt.crlID != "" {
				cobraCmd.SetArgs([]string{"--crl-id", tt.crlID})
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
