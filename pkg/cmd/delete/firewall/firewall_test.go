package firewall

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/delete/firewall"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockFirewallID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidFirewallID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorFirewallID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		firewallID     string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
		mockError      error
	}{
		{
			name:           "delete firewall by id",
			firewallID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockFirewallID,
			mockError:      nil,
		},
		{
			name:           "delete firewall - not found",
			firewallID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockFirewallID,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			firewallID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidFirewallID,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--firewall-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "ask for firewall id success",
			firewallID:     "",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockFirewallID,
			mockError:      nil,
		},
		{
			name:           "ask for firewall id conversion failure",
			firewallID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidFirewallID,
			mockError:      msg.ErrorConvertId,
		},
		{
			name:           "error - parse answer",
			firewallID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorFirewallID,
			mockError:      utils.ErrorParseResponse,
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

			if tt.firewallID != "" {
				cobraCmd.SetArgs([]string{"--firewall-id", tt.firewallID})
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
