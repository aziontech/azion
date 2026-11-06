package contracts

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The manifest the bundler writes with --config-version 6 decodes into
// ManifestV6 without losing any resource.
func TestManifestV6DecodesTheBundlerOutput(t *testing.T) {
	data, err := os.ReadFile("../manifest/v6/testdata/bundler666_full.json")
	require.NoError(t, err)

	var m ManifestV6
	require.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, 6, m.Version)
	assert.Len(t, m.Functions, 1)
	assert.Len(t, m.NetworkLists, 1)
	assert.Len(t, m.WAFs, 1)
	assert.Len(t, m.Connectors, 2)
	// SDK 0.4.0 has active on every Connector type again, so it is kept.
	for _, c := range m.Connectors {
		inst, ok := c.GetActualInstance().(interface{ GetActive() bool })
		require.True(t, ok, "connector type without active")
		assert.True(t, inst.GetActive())
	}
	assert.Len(t, m.CustomPages, 1)
	require.Len(t, m.Applications, 1)
	assert.Len(t, m.Applications[0].Rules, 3)
	assert.Equal(t, PhaseResponse, m.Applications[0].Rules[2].Phase)
	assert.Equal(t, "my-func", m.Applications[0].FunctionInstances[0].Function.Name)
	require.Len(t, m.Firewalls, 1)
	assert.Len(t, m.Firewalls[0].Rules, 2)
	require.Len(t, m.Environments, 1)
	assert.Equal(t, "single_version", m.Environments[0].DeploymentPolicy)
	require.Len(t, m.Deployments, 1)
	assert.Equal(t, "FLEXIBLE", m.Deployments[0].BindingPolicy)
	require.Len(t, m.Workloads, 1)
	assert.Equal(t, ManifestV6Binding{
		EnvironmentName: "production", DeploymentName: "my-deployment",
		Domains: []ManifestV6Domain{{Domain: "example.com"}},
	}, m.Workloads[0].Bindings[0])
	assert.Len(t, m.Purge, 1)
}
