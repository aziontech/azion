package applications

import (
	"context"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

const versions = "/workspace/applications/42/versions"

func client(api *v6test.API) *Client { return NewClient(api.Client(), v6test.URL, "token") }

func TestCreate(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", "/workspace/applications", v6test.JSON(v6test.Data(v6test.Version(42, "site", "V1", "draft"))))

	app, err := client(api).Create(context.Background(), sdk.ApplicationRequest{Name: "site"})

	require.NoError(t, err)
	assert.Equal(t, int64(42), app.GetId())
	assert.Equal(t, "V1", app.GetVersionId())
	assert.Equal(t, "site", api.Bodies("POST", "/workspace/applications")[0]["name"])
}

// The list filter may return near matches; only the exact name counts.
func TestFindByNameIsExact(t *testing.T) {
	api := &v6test.API{}
	api.On("GET", "/workspace/applications", v6test.JSON(v6test.Page(v6test.Version(1, "site-old", "V", "ready"), v6test.Version(2, "site", "V", "ready"))))

	app, found, err := client(api).FindByName(context.Background(), "site")

	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, int64(2), app.GetId())
}

func TestVersionOperations(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", versions, v6test.JSON(v6test.Version(42, "site", "V2", "draft")))
	api.On("PUT", versions+"/V2", v6test.JSON(v6test.Version(42, "site", "V2", "draft")))
	api.On("POST", versions+"/V2/build", v6test.Status(202))
	api.On("POST", versions+"/V2/cancel", v6test.Status(202))
	api.On("GET", versions+"/V2", v6test.JSON(v6test.Version(42, "site", "V2", "building")))
	c := client(api)
	ctx := context.Background()

	vid, err := c.CreateVersion(ctx, 42, "", "azion-cli deploy")
	require.NoError(t, err)
	assert.Equal(t, "V2", vid)
	assert.Equal(t, map[string]any{"comment": "azion-cli deploy"}, api.Bodies("POST", versions)[0], "no source_version: the API clones the latest Ready")

	_, err = c.UpdateVersion(ctx, 42, "V2", sdk.ApplicationRequest{Name: "site"})
	require.NoError(t, err)
	require.NoError(t, c.Build(ctx, 42, "V2", "deploy"))
	require.NoError(t, c.CancelBuild(ctx, 42, "V2"))

	state, err := c.VersionState(ctx, 42, "V2")
	require.NoError(t, err)
	assert.Equal(t, version.Building, state)
}

// Errors go through the same handling as v4: the same messages and sentinels.
func TestErrorsAreTheV4Ones(t *testing.T) {
	api := &v6test.API{}
	api.On("POST", "/workspace/applications", v6test.Status(409, map[string]any{"detail": "name in use"}))
	api.On("GET", versions+"/GONE", v6test.Status(404, map[string]any{"detail": "not found"}))
	api.On("PUT", versions+"/V1", v6test.Status(400, map[string]any{"detail": "modules: invalid"}))
	c := client(api)
	ctx := context.Background()

	_, err := c.Create(ctx, sdk.ApplicationRequest{Name: "site"})
	assert.ErrorIs(t, err, utils.ErrorNameInUse)

	_, err = c.GetVersion(ctx, 42, "GONE")
	assert.ErrorIs(t, err, utils.ErrorNotFound404)

	_, err = c.UpdateVersion(ctx, 42, "V1", sdk.ApplicationRequest{Name: "site"})
	assert.ErrorContains(t, err, "modules: invalid", "the 400 body reaches the user, as in v4")
}

// The reads a wait polls also carry the status, so the wait can retry a 429
// after Retry-After and stop on a 404.
func TestPolledReadsCarryTheStatus(t *testing.T) {
	api := &v6test.API{}
	api.On("GET", versions+"/V1", v6test.Response{Status: 429, Body: map[string]any{"detail": "slow down"}, Header: map[string]string{"Retry-After": "2"}})
	api.On("GET", versions+"/GONE", v6test.Status(404, map[string]any{"detail": "not found"}))
	c := client(api)

	_, err := c.VersionState(context.Background(), 42, "V1")
	retry, after := version.Retryable(err)
	assert.True(t, retry)
	assert.Equal(t, "2s", after.String())

	_, err = c.VersionState(context.Background(), 42, "GONE")
	retry, _ = version.Retryable(err)
	assert.False(t, retry)
	assert.ErrorIs(t, err, utils.ErrorNotFound404)
}

func TestSubResources(t *testing.T) {
	api := &v6test.API{}
	api.On("GET", versions+"/V1/cache_settings", v6test.JSON(v6test.Page(map[string]any{"id": 5, "name": "static"})))
	api.On("POST", versions+"/V1/request_rules", v6test.JSON(v6test.Data(map[string]any{"id": 7, "name": "r"})))
	api.On("PUT", versions+"/V1/request_rules/order", v6test.JSON(v6test.Page()))
	api.On("DELETE", versions+"/V1/device_groups/3", v6test.Status(200, map[string]any{}))
	c := client(api)
	ctx := context.Background()

	list, err := c.ListCacheSettings(ctx, 42, "V1", &contracts.ListOptions{Page: 1, PageSize: 100})
	require.NoError(t, err)
	assert.Equal(t, "static", list.GetResults()[0].GetName())

	rule, err := c.CreateRequestRule(ctx, 42, "V1", sdk.RequestPhaseRuleRequest{Name: "r"})
	require.NoError(t, err)
	assert.Equal(t, int64(7), rule.GetId())

	require.NoError(t, c.OrderRequestRules(ctx, 42, "V1", []int64{7, 8}))
	assert.Equal(t, []any{float64(7), float64(8)}, api.Bodies("PUT", versions+"/V1/request_rules/order")[0]["order"])

	require.NoError(t, c.DeleteDeviceGroup(ctx, 42, "V1", 3))
}
