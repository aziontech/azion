package dev

import (
	"errors"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDev(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name           string
		mockVulcan     func() *vulcanPkg.VulcanPkg
		mockCommandRun func(f *cmdutil.Factory, comm string) error
		expectedError  error
	}{
		{
			name: "dev - successful execution",
			mockVulcan: func() *vulcanPkg.VulcanPkg {
				return &vulcanPkg.VulcanPkg{
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "echo 1"
					},
				}
			},
			mockCommandRun: func(f *cmdutil.Factory, comm string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "dev - failed command execution",
			mockVulcan: func() *vulcanPkg.VulcanPkg {
				return &vulcanPkg.VulcanPkg{
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "echo 1"
					},
				}
			},
			mockCommandRun: func(f *cmdutil.Factory, comm string) error {
				return errors.New("failed to run command")
			},
			expectedError: errors.New("Error executing Bundler: Failed to run dev command. Verify if the command is correct and check the output above for more details. Run the 'azion dev' command again or contact Azion's support"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			f, _, _ := testutils.NewFactory(mock)

			devCmd := NewDevCmd(f)
			devCmd.Vulcan = tt.mockVulcan
			devCmd.CommandRunInteractive = tt.mockCommandRun

			err := devCmd.Run(f)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDevArgsPath(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name      string
		configDir string
		wantArg   string
	}{
		{name: "config-dir sent", configDir: "my-config", wantArg: " --args-path 'my-config/args.json'"},
		{name: "config-dir not sent", configDir: "", wantArg: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			f, _, _ := testutils.NewFactory(mock)

			var gotCommand string
			devCmd := NewDevCmd(f)
			devCmd.ConfigDir = tt.configDir
			devCmd.Vulcan = func() *vulcanPkg.VulcanPkg {
				return &vulcanPkg.VulcanPkg{
					Command: func(flags, params string, f *cmdutil.Factory) string {
						return "bundler dev"
					},
				}
			}
			devCmd.CommandRunInteractive = func(f *cmdutil.Factory, comm string) error {
				gotCommand = comm
				return nil
			}

			require.NoError(t, devCmd.Run(f))
			assert.Equal(t, "bundler dev"+tt.wantArg, gotCommand)
		})
	}
}
