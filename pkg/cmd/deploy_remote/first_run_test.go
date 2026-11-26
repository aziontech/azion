package deploy

import (
	"testing"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	dryrunmsg "github.com/aziontech/azion-cli/messages/dryrun"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestIgnoreSkipsOnFirstRun(t *testing.T) {
	skipBuild, skipFramework := true, true
	assert.True(t, IgnoreSkipsOnFirstRun(false, &skipBuild, &skipFramework))
	assert.False(t, skipBuild)
	assert.False(t, skipFramework)

	skipBuild, skipFramework = false, false
	assert.False(t, IgnoreSkipsOnFirstRun(false, &skipBuild, &skipFramework), "nothing to ignore")

	skipBuild, skipFramework = true, true
	assert.False(t, IgnoreSkipsOnFirstRun(true, &skipBuild, &skipFramework))
	assert.True(t, skipBuild)
	assert.True(t, skipFramework)
}

func TestFirstDeployIgnoresTheSkipFlags(t *testing.T) {
	for _, flags := range []struct {
		name             string
		build, framework bool
	}{
		{"skip-build", true, false},
		{"skip-framework-build", false, true},
		{"both", true, true},
	} {
		t.Run(flags.name, func(t *testing.T) {
			logger.New(zapcore.InfoLevel)
			cmd, out := dryRunCmd(t, &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj"})
			cmd.SkipBuild, cmd.SkipFramework = flags.build, flags.framework

			inTempProject(t, &contracts.ManifestV4{Storage: []contracts.StorageManifest{{}}}, func() {
				require.NoError(t, cmd.DryRun(cmd.F))
			})

			assert.False(t, cmd.SkipBuild)
			assert.False(t, cmd.SkipFramework)
			assert.Contains(t, out.String(), msg.FirstRunIgnoresSkips)
			assert.Contains(t, out.String(), dryrunmsg.BundlerInit)
		})
	}
}

func TestLaterDeployKeepsTheSkipFlags(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj", NotFirstRun: true}
	cmd, out := dryRunCmd(t, conf)
	cmd.SkipBuild, cmd.SkipFramework = true, true

	inTempProject(t, &contracts.ManifestV4{}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	assert.True(t, cmd.SkipBuild)
	assert.True(t, cmd.SkipFramework)
	assert.NotContains(t, out.String(), msg.FirstRunIgnoresSkips)
}
