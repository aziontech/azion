package streams

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/streams"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockStreamID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidStreamID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorStreamID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		streamID       string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
	}{
		{
			name:           "delete stream by id",
			streamID:       "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/streams/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockStreamID,
		},
		{
			name:           "delete stream - not found",
			streamID:       "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/streams/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockStreamID,
		},
		{
			name:           "error in input",
			streamID:       "1234",
			method:         "DELETE",
			endpoint:       "workspace/stream/streams/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalidStreamID,
		},
		{
			name:           "ask for stream id success",
			streamID:       "",
			method:         "DELETE",
			endpoint:       "workspace/stream/streams/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockStreamID,
		},
		{
			name:           "error - parse answer",
			streamID:       "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseErrorStreamID,
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

			if tt.streamID != "" {
				cobraCmd.SetArgs([]string{"--stream-id", tt.streamID})
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
