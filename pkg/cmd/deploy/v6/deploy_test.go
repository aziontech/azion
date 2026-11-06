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
