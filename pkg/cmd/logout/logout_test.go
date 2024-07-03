package logout

import (
	"context"
	"errors"
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/logout"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLogout(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name            string
		tokenSettings   token.Settings
		mockReadError   error
		mockDeleteError error
		mockWriteError  error
		expectedOutput  string
		expectedError   error
	}{
		{
			name: "logout - successful logout",
			tokenSettings: token.Settings{
				UUID: "1234-5678",
			},
			mockReadError:   nil,
			mockDeleteError: nil,
			mockWriteError:  nil,
			expectedOutput:  msg.SUCCESS,
			expectedError:   nil,
		},
		{
			name: "logout - no UUID",
			tokenSettings: token.Settings{
				UUID: "",
			},
			mockReadError:   nil,
			mockDeleteError: nil,
			mockWriteError:  nil,
			expectedOutput:  msg.SUCCESS,
			expectedError:   nil,
		},
		{
			name:            "logout - failed to read settings",
			tokenSettings:   token.Settings{},
			mockReadError:   errors.New("failed to get token dir"),
			mockDeleteError: nil,
			mockWriteError:  nil,
			expectedOutput:  "",
			expectedError:   errors.New("failed to get token dir"),
		},
		{
			name: "logout - failed to delete token",
			tokenSettings: token.Settings{
				UUID: "1234-5678",
			},
			mockReadError:   nil,
			mockDeleteError: errors.New("failed to delete token"),
			mockWriteError:  nil,
			expectedOutput:  "",
			expectedError:   fmt.Errorf(msg.ErrorLogout, "failed to delete token"),
		},
		{
			name: "logout - failed to write settings",
			tokenSettings: token.Settings{
				UUID: "1234-5678",
			},
			mockReadError:   nil,
			mockDeleteError: nil,
			mockWriteError:  errors.New("Failed to write settings.toml file"),
			expectedOutput:  "",
			expectedError:   errors.New("Failed to write settings.toml file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReadSettings := func() (token.Settings, error) {
				return tt.tokenSettings, tt.mockReadError
			}

			mockWriteSettings := func(settings token.Settings) error {
				return tt.mockWriteError
			}

			mockDeleteToken := func(ctx context.Context, uuid string) error {
				return tt.mockDeleteError
			}

			mock := &httpmock.Registry{}

			f, out, _ := testutils.NewFactory(mock)

			logoutCmd := &LogoutCmd{
				Io:            f.IOStreams,
				ReadSettings:  mockReadSettings,
				WriteSettings: mockWriteSettings,
				DeleteToken:   mockDeleteToken,
			}
			cmd := NewCobraCmd(logoutCmd, f)

			_, err := cmd.ExecuteC()
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, out.String())
			}
		})
	}
}
