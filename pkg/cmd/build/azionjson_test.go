package build

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// azionJson writes azion/azion.json in a fresh working directory.
func azionJson(t *testing.T, content string) string {
	t.Helper()
	logger.New(zapcore.DebugLevel)
	dir := t.TempDir()
	t.Chdir(dir)
	if content != "" {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, "azion"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "azion", "azion.json"), []byte(content), 0o644))
	}
	return filepath.Join(dir, "azion", "azion.json")
}

// What a v6 deploy leaves: workloads is a list, which the v4 contract cannot
// read.
const deployedV6 = `{
  "name": "gatsby-teste-1-prod-v6",
  "bucket": "gatsby-teste-1-prod-v6",
  "preset": "gatsby",
  "prefix": "20261111150000",
  "not-first-run": true,
  "deployments": [{"id": "A2YDGYHL", "name": "gatsby-teste-1-prod-v6", "release": {"id": "R1", "state": "ready", "activated": true}}],
  "workloads": [{"id": 1794417699, "name": "gatsby-teste-1-prod-v6", "url": "https://abc.zone.azionedge.net"}],
  "applications": [{"id": 1794413898, "name": "gatsby-teste-1-prod-v6", "version-id": "AGZJ7H9R", "version-state": "ready"}]
}`

func TestBuildReadsADeployedV6Project(t *testing.T) {
	azionJson(t, deployedV6)

	conf, err := getAzionJsonContent("azion")

	require.NoError(t, err)
	assert.Equal(t, "gatsby", conf.Preset)
	assert.Equal(t, "gatsby-teste-1-prod-v6", conf.Name)
}

// Writing back changes the preset and keeps everything the deploy recorded.
func TestBuildWritesBackOnlyThePresetOfAV6Project(t *testing.T) {
	file := azionJson(t, deployedV6)

	require.NoError(t, writeAzionJsonContent(&contracts.AzionApplicationOptions{Name: "ignored", Preset: "next"}, "azion"))

	conf, err := utils.GetAzionJsonContentV6("azion")
	require.NoError(t, err)
	assert.Equal(t, "next", conf.Preset)
	assert.Equal(t, "gatsby-teste-1-prod-v6", conf.Name)
	assert.Equal(t, "20261111150000", conf.Prefix)
	assert.Equal(t, []contracts.AzionJsonV6Workload{{ID: 1794417699, Name: "gatsby-teste-1-prod-v6", Url: "https://abc.zone.azionedge.net"}}, conf.Workloads)
	require.Len(t, conf.Deployments, 1)
	assert.True(t, conf.Deployments[0].Release.Activated)

	var raw map[string]json.RawMessage
	data, err := os.ReadFile(file)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, byte('['), raw["workloads"][0], "still the v6 format")
}

// A v4 azion.json takes the v4 path, as before.
func TestBuildKeepsTheV4FormatForAV4Project(t *testing.T) {
	file := azionJson(t, `{"name": "site", "preset": "react", "application": {"id": 1}, "workloads": {"id": 2}}`)

	conf, err := getAzionJsonContent("azion")
	require.NoError(t, err)
	assert.Equal(t, "react", conf.Preset)
	assert.EqualValues(t, 1, conf.Application.ID)

	conf.Preset = "vue"
	require.NoError(t, writeAzionJsonContent(conf, "azion"))
	data, err := os.ReadFile(file)
	require.NoError(t, err)
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.Equal(t, byte('{'), raw["workloads"][0], "still the v4 format")
}

// Without azion.json the build fails as it did before.
func TestBuildWithoutAzionJsonFailsAsBefore(t *testing.T) {
	azionJson(t, "")

	_, err := getAzionJsonContent("azion")

	assert.ErrorIs(t, err, utils.ErrorOpeningAzionJsonFile)
}
