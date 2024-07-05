package edgeapplication

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
)

func mockApplicationID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalid(msg string) (string, error) {
	return "invalid", nil
}

func mockParseError(msg string) (string, error) {
	return "invalid", fmt.Errorf("error parsing input")
}

func TestDeleteWithAskInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		applicationID  string
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
			name:           "delete application by id",
			applicationID:  "1234",
			method:         "DELETE",
			endpoint:       "edge_applications/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockApplicationID,
			mockError:      nil,
		},
		{
			name:           "delete application - not found",
			applicationID:  "1234",
			method:         "DELETE",
			endpoint:       "edge_applications/1234",
			statusCode:     404,
			responseBody:   "Not Found",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockApplicationID,
			mockError:      fmt.Errorf("Failed to parse your response. Check your response and try again. If the error persists, contact Azion support"),
		},
		{
			name:           "error in input",
			applicationID:  "1234",
			method:         "DELETE",
			endpoint:       "edge_applications/invalid",
			statusCode:     400,
			responseBody:   "Bad Request",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf("invalid argument \"\" for \"--application-id\" flag: strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name:           "ask for application id success",
			applicationID:  "",
			method:         "DELETE",
			endpoint:       "edge_applications/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.OutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockApplicationID,
			mockError:      nil,
		},
		{
			name:           "ask for application id conversion failure",
			applicationID:  "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockInvalid,
			mockError:      fmt.Errorf(msg.ErrorConvertId.Error()),
		},
		{
			name:           "error - parse answer",
			applicationID:  "",
			method:         "",
			endpoint:       "",
			statusCode:     0,
			responseBody:   "",
			expectedOutput: "",
			expectError:    true,
			mockInputs:     mockParseError,
			mockError:      fmt.Errorf("error parsing input"),
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
			cobraCmd := NewCobraCmd(deleteCmd)

			if tt.applicationID != "" {
				cobraCmd.SetArgs([]string{"--application-id", tt.applicationID})
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

func TestCascadeDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("cascade delete application", func(t *testing.T) {
		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{}

		dat, _ := os.ReadFile("./fixtures/azion.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/666"),
			httpmock.StatusStringResponse(204, ""),
		)
		mock.Register(
			httpmock.REST("DELETE", "edge_functions/123"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		del := NewDeleteCmd(f)
		del.GetAzion = func(confPath string) (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}
		del.UpdateJson = func(cmd *DeleteCmd) error {
			return nil
		}
		del.f = f
		del.Io = f.IOStreams

		cmd := NewCobraCmd(del)

		cmd.SetArgs([]string{"--cascade"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, msg.CascadeSuccess, stdout.String())
	})
}
