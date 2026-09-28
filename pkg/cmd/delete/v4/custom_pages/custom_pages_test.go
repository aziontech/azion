package custompages

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/custom_pages"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockCustomPageID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidCustomPageID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorCustomPageID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		customPageID   string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
	}{
		{
			name:           "delete custom page by id",
			customPageID:   "1234",
			method:         "DELETE",
			endpoint:       "workspace/custom_pages/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockCustomPageID,
		},
		{
			name:           "delete custom page - not found",
			customPageID:   "1234",
			method:         "DELETE",
			endpoint:       "workspace/custom_pages/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockCustomPageID,
		},
		{
			name:           "error in input",
			customPageID:   "1234",
			method:         "DELETE",
			endpoint:       "workspace/custom_pages/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidCustomPageID,
		},
		{
			name:           "ask for custom page id success",
			customPageID:   "",
			method:         "DELETE",
			endpoint:       "workspace/custom_pages/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockCustomPageID,
		},
		{
			name:           "error - parse answer",
			customPageID:   "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorCustomPageID,
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

			if tt.customPageID != "" {
				cobraCmd.SetArgs([]string{"--custom-page-id", tt.customPageID})
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
