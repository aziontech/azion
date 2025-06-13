package variables

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/variables"
)

func mockVariableID(msg string) (string, error) {
	return "7a187044-4a00-4a4a-93ed-d230900421f3", nil
}

func mockParse(msg string) (string, error) {
	return "7a187044-4a00-4a4a-93ed-d230900421f3", utils.ErrorParseResponse
}
func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		variableID     string
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
			name:           "delete variable by id success",
			variableID:     "7a187044-4a00-4a4a-93ed-d230900421f3",
			method:         "DELETE",
			endpoint:       "variables/7a187044-4a00-4a4a-93ed-d230900421f3",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "7a187044-4a00-4a4a-93ed-d230900421f3"),
			expectError:    false,
			mockError:      nil,
		},
		{
			name:           "delete variable not found",
			variableID:     "7a187044-4a00-4a4a-93ed-d230900421f3",
			method:         "DELETE",
			endpoint:       "variables/7a187044-4a00-4a4a-93ed-d230900421f3",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockError:      nil,
		},
		{
			name:           "delete variable ask id",
			method:         "DELETE",
			endpoint:       "variables/7a187044-4a00-4a4a-93ed-d230900421f3",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockVariableID,
			mockError:      nil,
		},
		{
			name:           "delete variable ask id error parse",
			method:         "DELETE",
			endpoint:       "variables/7a187044-4a00-4a4a-93ed-d230900421f3",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParse,
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

			// Create your delete command instance (adjust as per your package structure)
			deleteCmd := NewDeleteCmd(f)
			deleteCmd.AskInput = tt.mockInputs

			// Create your Cobra command instance (adjust as per your package structure)
			cobraCmd := NewCobraCmd(deleteCmd, f)

			if tt.variableID != "" {
				// Set command line arguments based on test case
				cobraCmd.SetArgs([]string{"--variable-id", tt.variableID})
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
