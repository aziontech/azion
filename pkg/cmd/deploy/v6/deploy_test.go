package deploy

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	remotemsg "github.com/aziontech/azion-cli/messages/deploy-remote"
	deployV6 "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote/v6"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

func TestFlagsReachThePipeline(t *testing.T) {
	out := &bytes.Buffer{}
	f := &cmdutil.Factory{IOStreams: &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}}, Config: viper.New()}
	d := NewDeployCmd(f)
	var got deployV6.Options
	d.run = func(_ context.Context, r *deployV6.Runner) error { got = r.Opts; return nil }

	c := NewCobraCmd(d)
	c.SetArgs([]string{"--environment", "Staging", "--no-activate", "--timeout", "2m", "--config-dir", "cfg", "--skip-build"})
	require.NoError(t, c.Execute())

	assert.Equal(t, "Staging", got.Environment)
	assert.True(t, got.NoActivate)
	assert.Equal(t, 2*time.Minute, got.Timeout)
	assert.Equal(t, "cfg", got.ConfigDir)
	assert.True(t, got.SkipBuild)
	// Without --local, the deploy says it runs locally (D1).
	assert.Contains(t, out.String(), remotemsg.V6RemoteNotAvailable)
}

func TestDefaultTimeoutIsThreeMinutes(t *testing.T) {
	f := &cmdutil.Factory{IOStreams: &iostreams.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}}, Config: viper.New()}
	d := NewDeployCmd(f)
	var got deployV6.Options
	d.run = func(_ context.Context, r *deployV6.Runner) error { got = r.Opts; return nil }
	c := NewCobraCmd(d)
	c.SetArgs([]string{"--local"})
	require.NoError(t, c.Execute())
	assert.Equal(t, 3*time.Minute, got.Timeout)
}

// pathTaken runs the v6 command with args and reports whether it went remote
// and whether it ran the pipeline here.
func pathTaken(t *testing.T, args ...string) (remote, local bool, err error) {
	t.Helper()
	f := &cmdutil.Factory{IOStreams: &iostreams.IOStreams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}}, Config: viper.New()}
	d := NewDeployCmd(f)
	d.Remote = func(context.Context, *DeployCmd) error { remote = true; return nil }
	d.run = func(context.Context, *deployV6.Runner) error { local = true; return nil }
	c := NewCobraCmd(d)
	c.SetArgs(args)
	err = c.Execute()
	return remote, local, err
}

// Without --local the deploy runs on the script runner, as in v4.
func TestDeployRunsRemotely(t *testing.T) {
	remote, local, err := pathTaken(t)
	require.NoError(t, err)
	assert.True(t, remote)
	assert.False(t, local)
}

// --local and --dry-run run here, as in v4.
func TestLocalAndDryRunRunHere(t *testing.T) {
	for _, flag := range []string{"--local", "--dry-run"} {
		remote, local, err := pathTaken(t, flag)
		require.NoError(t, err, flag)
		assert.False(t, remote, flag)
		assert.True(t, local, flag)
	}
}

// The script is given none of the v6 pipeline's options, so they need --local.
func TestPipelineFlagsNeedLocal(t *testing.T) {
	for _, args := range [][]string{{"--environment", "Staging"}, {"--no-activate"}, {"--timeout", "2m"}} {
		remote, local, err := pathTaken(t, args...)
		assert.ErrorIs(t, err, remotemsg.ErrorV6LocalOnlyFlags, args)
		assert.False(t, remote || local, args)

		_, local, err = pathTaken(t, append(args, "--local")...)
		require.NoError(t, err, args)
		assert.True(t, local, args)
	}
}
