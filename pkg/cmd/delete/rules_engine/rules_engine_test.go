package rulesengine

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockAppID(msg string) (string, error) {
	return "4321", nil
}

func mockRuleID(msg string) (string, error) {
	return "1234", nil
}

func mockPhase(msg string) (string, error) {
	return "request", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		applicationID  string
		ruleID         string
		phase          string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(msg string) (string, error)
		mockError      error
	}{
		{
			name:           "ask for application id and rule id success",
			applicationID:  "4321",
			ruleID:         "1234",
			phase:          "request",
			method:         "DELETE",
			endpoint:       "edge_applications/4321/rules_engine/request/rules/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockAppID,
			mockError:      nil,
		},
		{
			name:           "ask for rule id success",
			phase:          "request",
			applicationID:  "4321",
			method:         "DELETE",
			endpoint:       "edge_applications/4321/rules_engine/request/rules/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockRuleID,
			mockError:      nil,
		},
		{
			name:           "ask for phase success",
			applicationID:  "4321",
			ruleID:         "1234",
			method:         "DELETE",
			endpoint:       "edge_applications/4321/rules_engine/request/rules/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockPhase,
			mockError:      nil,
		},
		{
			name:           "error in input",
			applicationID:  "invalid",
			ruleID:         "1234",
			phase:          "request",
			method:         "DELETE",
			endpoint:       "edge_applications/invalid/rules_engine/request/rules/1234",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--application-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "error in input",
			phase:          "request",
			ruleID:         "1234",
			method:         "DELETE",
			endpoint:       "edge_applications/invalid/rules_engine/request/rules/1234",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--application-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
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

			// Create your delete command instance (adjust as per your package structure)
			deleteCmd := NewDeleteCmd(f)
			deleteCmd.AskInput = tt.mockInputs

			// Create your Cobra command instance (adjust as per your package structure)
			cobraCmd := NewCobraCmd(deleteCmd, f)

			// Set command line arguments based on test case
			if tt.applicationID != "" && tt.ruleID != "" && tt.phase != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID, "--rule-id", tt.ruleID, "--phase", tt.phase})
			} else if tt.applicationID != "" && tt.ruleID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID, "--rule-id", tt.ruleID})
			} else if tt.applicationID != "" && tt.phase != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID, "--phase", tt.phase})
			} else if tt.ruleID != "" && tt.phase != "" {
				cobraCmd.SetArgs([]string{"--rule-id", tt.ruleID, "--phase", tt.phase})
			} else if tt.applicationID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID})
			} else if tt.ruleID != "" {
				cobraCmd.SetArgs([]string{"--rule-id", tt.ruleID})
			} else if tt.phase != "" {
				cobraCmd.SetArgs([]string{"--phase", tt.phase})
			}

			// Execute the command
			_, err := cobraCmd.ExecuteC()

			// Validate the results
			if tt.expectError {
				require.Error(t, err)
				logger.Debug("Expected error occurred", zap.Error(err))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, stdout.String())
				logger.Debug("Expected output", zap.String("output", stdout.String()))
			}
		})
	}
}
