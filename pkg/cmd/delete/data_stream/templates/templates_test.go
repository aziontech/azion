package templates

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/templates"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockTemplateID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidTemplateID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorTemplateID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		templateID     string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
	}{
		{
			name:           "delete template by id",
			templateID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/templates/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockTemplateID,
		},
		{
			name:           "delete template - not found",
			templateID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/templates/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockTemplateID,
		},
		{
			name:           "error in input",
			templateID:     "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/templates/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidTemplateID,
		},
		{
			name:           "ask for template id success",
			templateID:     "",
			method:         "DELETE",
			endpoint:       "workspace/stream/templates/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockTemplateID,
		},
		{
			name:           "error - parse answer",
			templateID:     "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorTemplateID,
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

			if tt.templateID != "" {
				cobraCmd.SetArgs([]string{"--template-id", tt.templateID})
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
