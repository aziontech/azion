package firewallinstance

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/delete/firewall_instance"
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
		instanceID     string
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
			name:           "delete firewall instance by id",
			firewallID:     "1234",
			instanceID:     "5678",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234/functions/5678",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 5678),
			expectError:    false,
			mockInputs:     mockFirewallID,
			mockError:      nil,
		},
		{
			name:           "delete firewall instance - not found",
			firewallID:     "1234",
			instanceID:     "5678",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234/functions/5678",
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
			instanceID:     "5678",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/invalid/functions/5678",
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
			instanceID:     "5678",
			method:         "DELETE",
			endpoint:       "workspace/firewalls/1234/functions/5678",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 5678),
			expectError:    false,
			mockInputs: func(s string) (string, error) {
				if s == "Enter the ID of the Firewall:" {
					return "1234", nil
				}
				return "5678", nil
			},
			mockError: nil,
		},
		{
			name:           "ask for firewall id conversion failure",
			firewallID:     "",
			instanceID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidFirewallID,
			mockError:      msg.ErrorConvertFirewallId,
		},
		{
			name:           "error - parse answer",
			firewallID:     "",
			instanceID:     "",
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

			var args []string
			if tt.firewallID != "" {
				args = append(args, "--firewall-id", tt.firewallID)
			}
			if tt.instanceID != "" {
				args = append(args, "--instance-id", tt.instanceID)
			}
			cobraCmd.SetArgs(args)

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
