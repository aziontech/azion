package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestAzionJsonV6RoundTrip(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Chdir(t.TempDir())

	empty, err := GetAzionJsonContentV6("azion")
	require.NoError(t, err, "a first v6 deploy starts without azion.json")
	assert.Equal(t, &contracts.AzionApplicationOptionsV6{}, empty)

	conf := &contracts.AzionApplicationOptionsV6{
		Name:         "site",
		Applications: []contracts.AzionJsonV6Resource{{ID: 1, Name: "site", VersionID: "A1", VersionState: "ready"}},
		Deployments:  []contracts.AzionJsonV6Deployment{{Name: "site-deploy", ID: "DEP00001", Release: contracts.AzionJsonV6Release{ID: "R", Activated: true}}},
	}
	require.NoError(t, WriteAzionJsonContentV6(conf, "azion"))

	got, err := GetAzionJsonContentV6("azion")
	require.NoError(t, err)
	assert.Equal(t, conf.Applications, got.Applications)
	assert.Equal(t, conf.Deployments, got.Deployments)

	require.NoError(t, os.WriteFile(filepath.Join("azion", "azion.json"), []byte("{"), 0644))
	_, err = GetAzionJsonContentV6("azion")
	assert.ErrorIs(t, err, ErrorUnmarshalAzionJsonFile)
}

// A project made by the v4 "azion init" (every account got one before init
// learned about v6) is read as a fresh v6 project with the same name, preset,
// bucket and prefix.
func TestAzionJsonV6ReadsAV4Project(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll("azion", 0755))
	v4 := `{"name":"genuinely-chevalgrand","bucket":"b1","preset":"eleventy","env":"production","prefix":"p1",
	"rotate-prefix":true,"not-first-run":true,"function":[],"application":{"id":7,"name":"__DEFAULT__"},
	"domain":{"id":0,"name":"__DEFAULT__"},"rules-engine":{"created":false,"rules":null},
	"workloads":{"id":0,"name":"","domains":null,"url":"","deployment_id":null},"connectors":null}`
	require.NoError(t, os.WriteFile(filepath.Join("azion", "azion.json"), []byte(v4), 0644))

	got, err := GetAzionJsonContentV6("azion")
	require.NoError(t, err)
	rotate := true
	assert.Equal(t, &contracts.AzionApplicationOptionsV6{
		Name: "genuinely-chevalgrand", Bucket: "b1", Preset: "eleventy", Env: "production", Prefix: "p1",
		RotatePrefix: &rotate,
	}, got, "v4 ids are dropped and not-first-run is reset so the v6 placeholders get replaced")
}
