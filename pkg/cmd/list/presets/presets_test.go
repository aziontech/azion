package presets

import (
	"errors"
	"strings"
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

func newMockVulcan(captured *string) func() *vulcanPkg.VulcanPkg {
	return func() *vulcanPkg.VulcanPkg {
		return &vulcanPkg.VulcanPkg{
			CheckVulcanMajor: func(currentVersion string, f *cmdutil.Factory, vulcan *vulcanPkg.VulcanPkg) error {
				return nil
			},
			Command: func(flags, params string, f *cmdutil.Factory) string {
				cmd := "npx --yes " + flags + " @aziontech/bundler " + params
				if captured != nil {
					*captured = cmd
				}
				return cmd
			},
		}
	}
}

func TestPresets_run(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("happy path streams output without CLEAN_OUTPUT_MODE", func(t *testing.T) {
		f, stdout, _ := testutils.NewFactory(&httpmock.Registry{})

		var capturedCmd string
		list := &ListCmd{
			Io: f.IOStreams,
			F:  f,
			CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
				return "1.0.0", nil
			},
			CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
				capturedCmd = comm
				return nil
			},
			Vulcan: newMockVulcan(nil),
		}

		cmd := NewCobraCmd(list, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Contains(t, capturedCmd, "presets ls")
		assert.Contains(t, capturedCmd, "--loglevel=error --no-update-notifier")
		assert.NotContains(t, capturedCmd, "CLEAN_OUTPUT_MODE")
		assert.Contains(t, stdout.String(), "Fetching available presets")
	})

	t.Run("propagates error from npm show edge-functions version", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})

		list := &ListCmd{
			Io: f.IOStreams,
			F:  f,
			CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
				return "", errors.New("npm failed")
			},
			CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
				t.Fatalf("CommandRunInteractive should not be called when version lookup fails")
				return nil
			},
			Vulcan: newMockVulcan(nil),
		}

		cmd := NewCobraCmd(list, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "npm failed")
	})

	t.Run("propagates error from CheckVulcanMajor", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})

		list := &ListCmd{
			Io: f.IOStreams,
			F:  f,
			CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
				return "1.0.0", nil
			},
			CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
				t.Fatalf("CommandRunInteractive should not be called when vulcan major check fails")
				return nil
			},
			Vulcan: func() *vulcanPkg.VulcanPkg {
				return &vulcanPkg.VulcanPkg{
					CheckVulcanMajor: func(currentVersion string, f *cmdutil.Factory, vulcan *vulcanPkg.VulcanPkg) error {
						return errors.New("vulcan major mismatch")
					},
					Command: func(flags, params string, f *cmdutil.Factory) string { return "" },
				}
			},
		}

		cmd := NewCobraCmd(list, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "vulcan major mismatch")
	})

	t.Run("propagates error from bundler invocation", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})

		list := &ListCmd{
			Io: f.IOStreams,
			F:  f,
			CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
				return "1.0.0", nil
			},
			CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
				return errors.New("bundler failed")
			},
			Vulcan: newMockVulcan(nil),
		}

		cmd := NewCobraCmd(list, f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "bundler failed")
	})
}

func TestPresets_NewCobraCmd_metadata(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	list := &ListCmd{
		Io: f.IOStreams,
		F:  f,
		CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
			return "1.0.0", nil
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error { return nil },
		Vulcan:                newMockVulcan(nil),
	}

	cmd := NewCobraCmd(list, f)

	assert.Equal(t, "presets", cmd.Use)
	assert.True(t, cmd.SilenceUsage)
	assert.True(t, cmd.SilenceErrors)
	assert.NotNil(t, cmd.Flags().Lookup("help"))
	assert.True(t, strings.Contains(cmd.Example, "azion list presets"))
}
