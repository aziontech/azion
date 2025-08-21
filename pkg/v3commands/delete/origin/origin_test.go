package origin

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/origin"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mockAppID(msg string) (string, error) {
	return "1673635839", nil
}

func mockOriginKey(msg string) (string, error) {
	return "58755fef-e830-4ea4-b9e0-6481f1ef496d", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		applicationID  string
		originKey      string
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
			name:           "ask for application id success",
			applicationID:  "1673635839",
			originKey:      "58755fef-e830-4ea4-b9e0-6481f1ef496d",
			method:         "DELETE",
			endpoint:       "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "58755fef-e830-4ea4-b9e0-6481f1ef496d"),
			expectError:    false,
			mockInputs:     mockAppID,
			mockError:      nil,
		},
		{
			name:           "ask for origin key success",
			applicationID:  "1673635839",
			method:         "DELETE",
			endpoint:       "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "58755fef-e830-4ea4-b9e0-6481f1ef496d"),
			expectError:    false,
			mockInputs:     mockOriginKey,
			mockError:      nil,
		},
		{
			name:           "error in input",
			originKey:      "58755fef-e830-4ea4-b9e0-6481f1ef496d",
			method:         "DELETE",
			endpoint:       "edge_applications/invalid/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d",
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
			if tt.applicationID != "" && tt.originKey != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID, "--origin-key", tt.originKey})
			} else if tt.applicationID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID})
			} else {
				cobraCmd.SetArgs([]string{"--origin-key", tt.originKey})
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
