package application

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	msg "github.com/aziontech/azion-cli/messages/delete/application"
)

func mockApplicationID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func mockParseError(msg string) (string, error) {
	return "", fmt.Errorf("error parsing input")
}

func TestDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		args           []string
		statusCode     int
		responseBody   string
		expectedOutput string
		mockInputs     func(string) (string, error)
		err            error
	}{
		{
			name:           "delete application by id",
			args:           []string{"--application-id", "1234"},
			statusCode:     http.StatusOK,
			responseBody:   `{"state": "executed"}`,
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
		},
		{
			name:         "delete application - not found",
			args:         []string{"--application-id", "1234"},
			statusCode:   http.StatusNotFound,
			responseBody: "Not Found",
			err:          utils.ErrorNotFound404,
		},
		{
			name:           "ask for application id success",
			statusCode:     http.StatusOK,
			responseBody:   `{"state": "executed"}`,
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			mockInputs:     mockApplicationID,
		},
		{
			name:       "ask for application id conversion failure",
			mockInputs: mockInvalid,
			err:        msg.ErrorConvertId,
		},
		{
			name:       "error - parse answer",
			mockInputs: mockParseError,
			err:        fmt.Errorf("error parsing input"),
		},
		{
			name: "cascade is not offered on v6",
			args: []string{"--cascade"},
			err:  fmt.Errorf("unknown flag: --cascade"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(
				httpmock.REST("DELETE", "workspace/applications/1234"),
				httpmock.WithHeader(httpmock.StatusStringResponse(tt.statusCode, tt.responseBody), "Content-Type", "application/json"),
			)

			f, stdout, _ := testutils.NewFactory(mock)
			deleteCmd := NewDeleteCmd(f)
			if tt.mockInputs != nil {
				deleteCmd.AskInput = tt.mockInputs
			}
			cobraCmd := NewCobraCmd(deleteCmd)
			cobraCmd.SetArgs(tt.args)

			_, err := cobraCmd.ExecuteC()

			if tt.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, stdout.String())
		})
	}
}
