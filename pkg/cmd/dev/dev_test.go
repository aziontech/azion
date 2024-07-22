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
		isFirewall     bool
		expectedError  error
	}{
		{
			name: "dev - successful execution without firewall",
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
			isFirewall:    false,
			expectedError: nil,
		},
		{
			name: "dev - successful execution with firewall",
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
			isFirewall:    true,
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
			isFirewall:    false,
			expectedError: errors.New("Error executing Vulcan: Failed to run dev command. Verify if the command is correct and check the output above for more details. Run the 'azion dev' command again or contact Azion's support"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			f, _, _ := testutils.NewFactory(mock)

			devCmd := NewDevCmd(f)
			devCmd.Vulcan = tt.mockVulcan
			devCmd.CommandRunInteractive = tt.mockCommandRun

			isFirewall = tt.isFirewall

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
