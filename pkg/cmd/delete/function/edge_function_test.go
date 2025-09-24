package function

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/function"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockFunctionID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidFunctionID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorFunctionID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		functionID     string
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
			name:           "delete function by id",
			functionID:     "1234",
			method:         "DELETE",
			endpoint:       "edge_functions/functions/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "1234"),
			expectError:    false,
			mockInputs:     mockFunctionID,
			mockError:      nil,
		},
		{
			name:           "delete function - not found",
			functionID:     "1234",
			method:         "DELETE",
			endpoint:       "edge_functions/functions/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockFunctionID,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			functionID:     "1234",
			method:         "DELETE",
			endpoint:       "edge_functions/functions/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidFunctionID,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--function-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "ask for function id success",
			functionID:     "",
			method:         "DELETE",
			endpoint:       "edge_functions/functions/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "1234"),
			expectError:    false,
			mockInputs:     mockFunctionID,
			mockError:      nil,
		},
		{
			name:           "ask for function id conversion failure",
			functionID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidFunctionID,
			mockError:      msg.ErrorConvertIdFunction,
		},
		{
			name:           "error - parse answer",
			functionID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorFunctionID,
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

			if tt.functionID != "" {
				cobraCmd.SetArgs([]string{"--function-id", tt.functionID})
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
