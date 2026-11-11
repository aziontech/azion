package deployments

import (
	"context"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

const rel = "/deployments/DEP00001/releases"

func TestReleaseLifecycle(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", rel, v6test.JSON(v6test.Data(map[string]any{"id": "REL00001", "state": "draft"})))
	api.On("POST", rel+"/REL00001/build", v6test.JSON(v6test.Data(map[string]any{"id": "REL00001", "state": "queued"})))
	api.On("POST", rel+"/REL00001/activate", v6test.JSON(v6test.Data(map[string]any{"id": "REL00001", "state": "ready"})))
	api.On("GET", rel+"/REL00001", v6test.JSON(v6test.Data(map[string]any{"id": "REL00001", "state": "ready", "traffic_role": "ACTIVE", "trace_id": "T"})))
	c := NewClient(api.Client(), v6test.URL, "t")
	ctx := context.Background()

	r, err := c.CreateRelease(ctx, "DEP00001", []ReleaseResource{{ResourceID: 1, ResourceType: "application", VersionID: "latest"}}, OriginCLI)
	require.NoError(t, err)
	assert.Equal(t, version.Draft, r.State)
	body := api.Bodies("POST", rel)[0]
	assert.Equal(t, []any{map[string]any{"resource_id": float64(1), "resource_type": "application", "version_id": "latest"}}, body["resources"])
	assert.Equal(t, "cli", body["origin"].(map[string]any)["type"])

	r, err = c.BuildRelease(ctx, "DEP00001", "REL00001")
	require.NoError(t, err)
	assert.Equal(t, version.Queued, r.State)

	_, err = c.ActivateRelease(ctx, "DEP00001", "REL00001")
	require.NoError(t, err)

	r, err = c.GetRelease(ctx, "DEP00001", "REL00001")
	require.NoError(t, err)
	assert.Equal(t, RoleActive, r.TrafficRole)
	assert.Equal(t, "T", r.TraceID)
}

func TestHierarchyFlattensEveryResourceType(t *testing.T) {
	api := &v6test.API{}
	entry := func(kind, id, v string) map[string]any {
		return map[string]any{"resource_type": kind, "resource_id": id, "resource_version": v, "product_version": nil}
	}
	api.On("GET", "/deployments/DEP00001/resource_hierarchy", v6test.JSON(v6test.Data(map[string]any{
		"deployment_id": "DEP00001", "active_release_id": "REL00001",
		"resources": map[string]any{
			"application": entry("application", "1", "A1"), "firewall": nil, "custom_page": nil,
			"function": []any{entry("function", "4", "F1")}, "waf": []any{}, "connector": []any{}, "network_list": []any{},
		},
	})))
	c := NewClient(api.Client(), v6test.URL, "t")

	h, err := c.Hierarchy(context.Background(), "DEP00001")

	require.NoError(t, err)
	assert.Equal(t, "REL00001", h.ActiveReleaseID)
	assert.Equal(t, []HierarchyEntry{{"application", "1", "A1"}, {"function", "4", "F1"}}, h.Resources)
}
