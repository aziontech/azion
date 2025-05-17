package cachesetting

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockAppId(msg string) (string, error) {
	return "1673635839", nil
}

func mockCacheId(msg string) (string, error) {
	return "107313", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		applicationID  string
		cacheSettingID string
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
			cacheSettingID: "107313",
			method:         "DELETE",
			endpoint:       "edge_application/applications/1673635839/cache_settings/107313",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 107313),
			expectError:    false,
			mockInputs:     mockAppId,
			mockError:      nil,
		},
		{
			name:           "ask for cache setting id success",
			applicationID:  "1673635839",
			method:         "DELETE",
			endpoint:       "edge_application/applications/1673635839/cache_settings/107313",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 107313),
			expectError:    false,
			mockInputs:     mockCacheId,
			mockError:      nil,
		},
		{
			name:           "error in input",
			cacheSettingID: "107313",
			method:         "DELETE",
			endpoint:       "edge_application/applications/1673635839/cache_settings/107313",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--cache-setting-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
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

			if tt.applicationID != "" && tt.cacheSettingID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID, "--cache-setting-id", tt.cacheSettingID})
			} else if tt.applicationID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID})
			} else {
				cobraCmd.SetArgs([]string{"--cache-setting-id", tt.cacheSettingID})
			}

			_, err := cobraCmd.ExecuteC()
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
