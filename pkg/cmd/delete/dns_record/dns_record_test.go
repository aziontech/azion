package dnsrecord

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/dns_record"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		args           []string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
	}{
		{
			name:           "delete DNS record by id",
			args:           []string{"--zone-id", "1337", "--record-id", "1234"},
			method:         "DELETE",
			endpoint:       "workspace/dns/zones/1337/records/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DNSRecordDeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockID,
		},
		{
			name:         "delete DNS record - not found",
			args:         []string{"--zone-id", "1337", "--record-id", "1234"},
			method:       "DELETE",
			endpoint:     "workspace/dns/zones/1337/records/1234",
			statusCode:   404,
			responseBody: "Not Found",
			expectError:  true,
			mockInputs:   mockID,
		},
		{
			name:           "ask for ids success",
			args:           nil,
			method:         "DELETE",
			endpoint:       "workspace/dns/zones/1234/records/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DNSRecordDeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockID,
		},
		{
			name:        "ask for id conversion failure",
			args:        nil,
			method:      "",
			endpoint:    "",
			statusCode:  0,
			expectError: true,
			mockInputs:  mockInvalidID,
		},
		{
			name:        "error - parse answer",
			args:        nil,
			method:      "",
			endpoint:    "",
			statusCode:  0,
			expectError: true,
			mockInputs:  mockParseErrorID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(
				httpmock.REST(tt.method, tt.endpoint),
				httpmock.StatusStringResponse(tt.statusCode, tt.responseBody),
			)

			f, stdout, _ := testutils.NewFactory(mock)

			deleteCmd := NewDeleteCmd(f)
			deleteCmd.AskInput = tt.mockInputs
			cobraCmd := NewCobraCmd(deleteCmd, f)

			if tt.args != nil {
				cobraCmd.SetArgs(tt.args)
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
