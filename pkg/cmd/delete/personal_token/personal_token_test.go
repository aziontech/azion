package personaltoken

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/delete/personal_token"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func mockTokenID(msg string) (string, error) {
	return "5c9c1854-45dd-11ee-be56-0242ac120002", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteCmd(t *testing.T) {
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
		mockInputs     func(msg string) (string, error)
		mockError      error
	}{
		{
			name:           "Delete personal token by ID",
			args:           []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			method:         "DELETE",
			endpoint:       "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, "5c9c1854-45dd-11ee-be56-0242ac120002"),
			expectError:    false,
			mockInputs:     mockTokenID,
			mockError:      nil,
		},
		{
			name:         "Delete personal token not found",
			args:         []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			method:       "DELETE",
			endpoint:     "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002",
			statusCode:   404,
			responseBody: "Not Found",
			expectError:  true,
			mockInputs:   mockTokenID,
			mockError:    fmt.Errorf(msg.ErrorFailToDelete, utils.ErrorNotFound404),
		},
		{
			name:           "error - parse answer",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
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

			cobraCmd.SetArgs(tt.args)

			// Execute the command and capture any error
			_, err := cobraCmd.ExecuteC()

			if tt.expectError {
				require.Error(t, err)
				assert.EqualError(t, err, tt.mockError.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, stdout.String())
			}
		})
	}
}
