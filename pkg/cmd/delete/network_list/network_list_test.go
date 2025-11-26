package networklist

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/network_list"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockNetworkListID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidNetworkListID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorNetworkListID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		networkListID  string
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
			name:           "delete network list by id",
			networkListID:  "1234",
			method:         "DELETE",
			endpoint:       "workspace/network_lists/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "1234"),
			expectError:    false,
			mockInputs:     mockNetworkListID,
			mockError:      nil,
		},
		{
			name:           "delete network list - not found",
			networkListID:  "1234",
			method:         "DELETE",
			endpoint:       "workspace/network_lists/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockNetworkListID,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			networkListID:  "1234",
			method:         "DELETE",
			endpoint:       "workspace/network_lists/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidNetworkListID,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--network-list-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "ask for network list id success",
			networkListID:  "",
			method:         "DELETE",
			endpoint:       "workspace/network_lists/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, "1234"),
			expectError:    false,
			mockInputs:     mockNetworkListID,
			mockError:      nil,
		},
		{
			name:           "error - parse answer",
			networkListID:  "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorNetworkListID,
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

			if tt.networkListID != "" {
				cobraCmd.SetArgs([]string{"--network-list-id", tt.networkListID})
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
