package rollback

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func mockInvalidOriginKey(msg string) (string, error) {
	return "invalid", nil
}

func mockParseError(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestRollbackWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		originKey      string
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
			name:           "rollback with invalid origin key",
			originKey:      "invalid",
			method:         "UPDATE",
			endpoint:       "origins/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidOriginKey,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			originKey:      "invalid",
			method:         "UPDATE",
			endpoint:       "origins/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidOriginKey,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--origin-key\" flag: invalid syntax"),
		},
		{
			name:           "error - parse answer",
			originKey:      "",
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

			rollbackCmd := NewDeleteCmd(f)
			rollbackCmd.AskInput = tt.mockInputs
			rollbackCmd.GetAzionJsonContent = func(pathConf string) (*contracts.AzionApplicationOptions, error) {
				return &contracts.AzionApplicationOptions{
					Application: contracts.AzionJsonDataApplication{
						ID:   0001110001,
						Name: "namezin",
					},
					Bucket: "nomedobucket",
					Prefix: "001001001",
				}, nil
			}
			rollbackCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions, confPath string) error {
				return nil
			}
			cobraCmd := NewCobraCmd(rollbackCmd, f)

			if tt.originKey != "" {
				cobraCmd.SetArgs([]string{"--origin-key", tt.originKey})
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
