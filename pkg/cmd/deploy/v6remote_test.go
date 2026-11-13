package deploy

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/apiversion"
	v6deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy/v6"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// v6Project writes azion/azion.json in a fresh working directory.
func v6Project(t *testing.T, azionJson string) {
	t.Helper()
	logger.New(zapcore.DebugLevel)
	dir := t.TempDir()
	t.Chdir(dir)
	if azionJson != "" {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "azion"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "azion", "azion.json"), []byte(azionJson), 0o644))
	}
}

// What a v6 deploy leaves in azion.json.
const deployedV6Json = `{
  "name": "site", "bucket": "site", "preset": "react", "prefix": "20261111150000", "not-first-run": true,
  "deployments": [{"id": "DEP1", "name": "site", "release": {"id": "R1", "state": "ready", "activated": true}}],
  "workloads": [{"id": 30, "name": "site", "url": "https://old.azion.app"}]
}`

// The result the script runner returns once the deploy on Azion's side is
// done: its azion.json is in the v6 format.
const remoteResultV6 = `{"result": {"errors": null, "azion": {
  "name": "site", "bucket": "site", "preset": "react", "prefix": "P2", "not-first-run": true,
  "deployments": [{"id": "DEP1", "name": "site", "release": {"id": "R2", "state": "ready", "activated": true}}],
  "workloads": [{"id": 30, "name": "site", "url": "https://site.azion.app"}]
}}}`

// On a v6 account, init and link get the v6 deploy with the remote deploy
// wired in, so a deploy from them runs remotely as `azion deploy` does.
func TestExternalDeployerGivesV6TheRemoteDeploy(t *testing.T) {
	d, ok := NewExternalDeployer(&cmdutil.Factory{APIVersion: apiversion.V6}).(*v6deploy.DeployCmd)
	require.True(t, ok)
	assert.NotNil(t, d.Remote)
}

// The v6 remote deploy runs the v4 steps (upload, script, logs) on a v6
// azion.json, and saves the result in the v6 format.
func TestRemoteDeployV6(t *testing.T) {
	v6Project(t, deployedV6Json)
	f, out, _ := testutils.NewFactory(nil)

	d := v6deploy.NewDeployCmd(f)
	d.ProjectConf = "azion"
	d.ConfigFile = ""
	cmd := remoteDeployV6(f, d)

	var uploadedPrefix, scriptPrefix, scriptBucket, scriptConfDir string
	cmd.CheckToken = func(*cmdutil.Factory) error { return nil }
	cmd.ReadSettings = func(string) (token.Settings, error) {
		return token.Settings{Token: "tok", S3AccessKey: "key", S3SecretKey: "secret", S3Bucket: "site-bucket"}, nil
	}
	cmd.VersionID = func() string { return "P2" }
	cmd.UploadFiles = func(_ *cmdutil.Factory, conf *contracts.AzionApplicationOptions, _ *[]string, _, _ string, _ *DeployCmd, _ token.Settings) error {
		uploadedPrefix = conf.Prefix
		return nil
	}
	cmd.CallScript = func(_, _, _, prefix, name, confDir string, _ *DeployCmd) (string, error) {
		scriptPrefix, scriptBucket, scriptConfDir = prefix, name, confDir
		return "exec-1", nil
	}
	cmd.OpenBrowser = func(*cmdutil.Factory, string, *DeployCmd) error { return nil }
	cmd.CaptureLogs = func(_ string, _ string, c *DeployCmd) error {
		return c.SaveResult([]byte(remoteResultV6), c)
	}
	cmd.WriteAzionConfig = func(*contracts.AzionConfig) error { return nil }

	require.NoError(t, cmd.Run(f))

	// The prefix rotates as in v4, and the same one is uploaded and sent.
	assert.Equal(t, "P2", uploadedPrefix)
	assert.Equal(t, "P2", scriptPrefix)
	assert.Equal(t, "site-bucket", scriptBucket)
	assert.Equal(t, "azion", scriptConfDir)

	// azion.json is the deploy's, still in the v6 format.
	conf, err := utils.GetAzionJsonContentV6("azion")
	require.NoError(t, err)
	assert.Equal(t, "P2", conf.Prefix)
	assert.Equal(t, "R2", conf.Deployments[0].Release.ID)
	assert.Equal(t, []contracts.AzionJsonV6Workload{{ID: 30, Name: "site", Url: "https://site.azion.app"}}, conf.Workloads)
	var raw map[string]json.RawMessage
	data, err := os.ReadFile(filepath.Join("azion", "azion.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, byte('['), bytes.TrimSpace(raw["workloads"])[0])

	// The URL shown is the deployed Workload's.
	assert.Contains(t, out.String(), "https://site.azion.app")
}

func TestRemoteDeployV6ReportsTheScriptError(t *testing.T) {
	v6Project(t, deployedV6Json)
	f, _, _ := testutils.NewFactory(nil)
	cmd := remoteDeployV6(f, v6deploy.NewDeployCmd(f))
	cmd.WriteAzionConfig = func(*contracts.AzionConfig) error { return nil }

	err := cmd.SaveResult([]byte(`{"result": {"errors": {"stack": "build failed"}}}`), cmd)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "build failed")
	conf, err := utils.GetAzionJsonContentV6("azion")
	require.NoError(t, err)
	assert.Equal(t, "20261111150000", conf.Prefix, "azion.json is left as it was")
}

// Without azion.json the remote deploy stops as the v4 one does.
func TestRemoteDeployV6NeedsAzionJson(t *testing.T) {
	v6Project(t, "")

	_, err := azionJsonForRemoteV6("azion")

	assert.ErrorIs(t, err, utils.ErrorOpeningAzionJsonFile)
}

// The project name and prefix come from a v6 azion.json.
func TestRemoteDeployV6ReadsTheV6AzionJson(t *testing.T) {
	v6Project(t, deployedV6Json)

	conf, err := azionJsonForRemoteV6("azion")

	require.NoError(t, err)
	assert.Equal(t, "site", conf.Name)
	assert.Equal(t, "20261111150000", conf.Prefix)
	assert.True(t, conf.NotFirstRun)
}
