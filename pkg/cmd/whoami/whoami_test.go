package whoami

import (
    "errors"
    "testing"

    msg "github.com/aziontech/azion-cli/messages/whoami"
    "github.com/aziontech/azion-cli/pkg/httpmock"
    "github.com/aziontech/azion-cli/pkg/logger"
    "github.com/aziontech/azion-cli/pkg/testutils"
    "github.com/aziontech/azion-cli/pkg/token"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/zap/zapcore"
)

func TestWhoami(t *testing.T) {
    logger.New(zapcore.DebugLevel)

    tests := []struct {
        name           string
        tokenSettings  token.Settings
        mockReadError  error
        expectedOutput string
        expectedError  error
    }{
        {
            name: "whoami - logged in",
            tokenSettings: token.Settings{
                Email:    "test@example.com",
                ClientId: "abcd-1234",
            },
            mockReadError:  nil,
            expectedOutput: " Client ID: abcd-1234\n Email: test@example.com\n",
            expectedError:  nil,
        },
        {
            name: "whoami - not logged in",
            tokenSettings: token.Settings{
                Email: "",
            },
            mockReadError:  nil,
            expectedOutput: "",
            expectedError:  msg.ErrorNotLoggedIn,
        },
        {
            name:           "whoami - failed to read settings",
            tokenSettings:  token.Settings{},
            mockReadError:  errors.New("failed to get token dir"),
            expectedOutput: "",
            expectedError:  errors.New("failed to get token dir"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockReadSettings := func() (token.Settings, error) {
                return tt.tokenSettings, tt.mockReadError
            }

            mock := &httpmock.Registry{}

            f, out, _ := testutils.NewFactory(mock)
            f.Flags.NoColor = true

            whoamiCmd := &WhoamiCmd{
                Io:           f.IOStreams,
                ReadSettings: mockReadSettings,
                F:            f,
            }
            cmd := NewCobraCmd(whoamiCmd, f)

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
