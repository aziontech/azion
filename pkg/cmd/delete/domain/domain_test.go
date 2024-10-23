package domain

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/delete/domain"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockDomainID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func mockParseError(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		domainID       string
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
			name:           "delete domain by id",
			domainID:       "1234",
			method:         "DELETE",
			endpoint:       "domains/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockDomainID,
			mockError:      nil,
		},
		{
			name:           "delete domain - not found",
			domainID:       "1234",
			method:         "DELETE",
			endpoint:       "domains/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockDomainID,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			domainID:       "1234",
			method:         "DELETE",
			endpoint:       "domains/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--domain-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "ask for domain id success",
			domainID:       "",
			method:         "DELETE",
			endpoint:       "domains/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockDomainID,
			mockError:      nil,
		},
		{
			name:           "ask for domain id conversion failure",
			domainID:       "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      msg.ErrorConvertId,
		},
		{
			name:           "error - parse answer",
			domainID:       "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseError,
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

			if tt.domainID != "" {
				cobraCmd.SetArgs([]string{"--domain-id", tt.domainID})
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
