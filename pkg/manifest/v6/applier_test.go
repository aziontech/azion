package v6

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/v6/applications"
	apiConnectors "github.com/aziontech/azion-cli/pkg/api/v6/connectors"
	apiCustomPages "github.com/aziontech/azion-cli/pkg/api/v6/custompages"
	apiFirewalls "github.com/aziontech/azion-cli/pkg/api/v6/firewalls"
	apiFunctions "github.com/aziontech/azion-cli/pkg/api/v6/functions"
	apiNetworkLists "github.com/aziontech/azion-cli/pkg/api/v6/networklists"
	. "github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	apiWafs "github.com/aziontech/azion-cli/pkg/api/v6/wafs"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

func clients(api *API) *Clients {
	c := api.Client()
	return &Clients{
		Applications: apiApplications.NewClient(c, URL, "t"),
		Functions:    apiFunctions.NewClient(c, URL, "t"),
		Connectors:   apiConnectors.NewClient(c, URL, "t"),
		Firewalls:    apiFirewalls.NewClient(c, URL, "t"),
		CustomPages:  apiCustomPages.NewClient(c, URL, "t"),
		NetworkLists: apiNetworkLists.NewClient(c, URL, "t"),
		WAFs:         apiWafs.NewClient(c, URL, "t"),
	}
}

func applier(api *API, conf *contracts.AzionApplicationOptionsV6, m *contracts.ManifestV6) (*Applier, *int) {
	writes := 0
	return &Applier{
		Clients: clients(api), Conf: conf, Manifest: m,
		WriteConfig: func() error { writes++; return nil },
	}, &writes
}

var networkList = sdk.NetworkListRequest{Name: "blocked", Type: "ip_cidr", Items: []string{"10.0.0.0/8"}}

const nl = "/workspace/network_lists"

// One case per row of the Draft table in doc/V6-DEPLOY-PLAN.md section 5.
func TestDraftResolution(t *testing.T) {
	tests := []struct {
		name     string
		recorded []contracts.AzionJsonV6Resource
		routes   func(*API)
		wantID   int64
		wantVID  string
		calls    map[string]int
	}{
		{
			name: "no id and no resource by that name: create it",
			routes: func(a *API) {
				a.On("GET", nl, JSON(Page()))
				a.On("POST", nl, JSON(Data(Version(7, "blocked", "V1", "draft"))))
				a.On("GET", nl+"/7/versions/V1", JSON(Version(7, "blocked", "V1", "draft")))
			},
			wantID: 7, wantVID: "V1",
			calls: map[string]int{"POST " + nl: 1, "PUT " + nl + "/7/versions/V1": 0},
		},
		{
			name: "no id but a resource by that name exists: add a version to it",
			routes: func(a *API) {
				a.On("GET", nl, JSON(Page(Version(9, "blocked-old", "X", "ready"), Version(7, "blocked", "R1", "ready"))))
				a.On("POST", nl+"/7/versions", JSON(Version(7, "blocked", "V2", "draft")))
				a.On("PUT", nl+"/7/versions/V2", JSON(Version(7, "blocked", "V2", "draft")))
				a.On("GET", nl+"/7/versions/V2", JSON(Version(7, "blocked", "V2", "draft")))
			},
			wantID: 7, wantVID: "V2",
			calls: map[string]int{"POST " + nl: 0, "POST " + nl + "/7/versions": 1},
		},
		{
			name:     "recorded Draft from a previous run: reuse it",
			recorded: []contracts.AzionJsonV6Resource{{ID: 7, Name: "blocked", VersionID: "V1", VersionState: "error"}},
			routes: func(a *API) {
				a.On("GET", nl+"/7/versions/V1", JSON(Version(7, "blocked", "V1", "error")), JSON(Version(7, "blocked", "V1", "draft")))
				a.On("PUT", nl+"/7/versions/V1", JSON(Version(7, "blocked", "V1", "draft")))
			},
			wantID: 7, wantVID: "V1",
			calls: map[string]int{"POST " + nl + "/7/versions": 0, "PUT " + nl + "/7/versions/V1": 1},
		},
		{
			name:     "recorded version is ready: open a Draft from the latest Ready",
			recorded: []contracts.AzionJsonV6Resource{{ID: 7, Name: "blocked", VersionID: "V1", VersionState: "ready"}},
			routes: func(a *API) {
				a.On("GET", nl+"/7/versions/V1", JSON(Version(7, "blocked", "V1", "ready")))
				a.On("POST", nl+"/7/versions", JSON(Version(7, "blocked", "V2", "draft")))
				a.On("PUT", nl+"/7/versions/V2", JSON(Version(7, "blocked", "V2", "draft")))
				a.On("GET", nl+"/7/versions/V2", JSON(Version(7, "blocked", "V2", "draft")))
			},
			wantID: 7, wantVID: "V2",
			calls: map[string]int{"POST " + nl + "/7/versions": 1},
		},
		{
			name:     "recorded resource no longer exists: create it again",
			recorded: []contracts.AzionJsonV6Resource{{ID: 7, Name: "blocked"}},
			routes: func(a *API) {
				a.On("POST", nl+"/7/versions", Status(404, map[string]any{"detail": "not found"}))
				a.On("POST", nl, JSON(Data(Version(8, "blocked", "V1", "draft"))))
				a.On("GET", nl+"/8/versions/V1", JSON(Version(8, "blocked", "V1", "draft")))
			},
			wantID: 8, wantVID: "V1",
			calls: map[string]int{"POST " + nl: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{}
			tt.routes(api)
			conf := &contracts.AzionApplicationOptionsV6{NetworkLists: tt.recorded}
			a, writes := applier(api, conf, &contracts.ManifestV6{NetworkLists: []sdk.NetworkListRequest{networkList}})

			pending, err := a.ApplyNetworkLists(context.Background())

			require.NoError(t, err, api.Summary())
			require.Len(t, pending, 1)
			assert.Equal(t, tt.wantID, pending[0].ID)
			assert.Equal(t, tt.wantVID, pending[0].VersionID)
			assert.Equal(t, []contracts.AzionJsonV6Resource{{ID: tt.wantID, Name: "blocked", VersionID: tt.wantVID, VersionState: "draft"}}, conf.NetworkLists)
			assert.Positive(t, *writes, "azion.json is saved after each Draft")
			for call, n := range tt.calls {
				parts := splitCall(call)
				assert.Equal(t, n, api.Count(parts[0], parts[1]), call)
			}
		})
	}
}

func splitCall(call string) [2]string {
	for i := range call {
		if call[i] == ' ' {
			return [2]string{call[:i], call[i+1:]}
		}
	}
	return [2]string{call, ""}
}

// The first check of section 6: if someone builds the Draft while the CLI is
// writing it, the deploy stops instead of building over it.
func TestDraftTakenByAnotherProcess(t *testing.T) {
	api := &API{}
	api.On("GET", nl, JSON(Page()))
	api.On("POST", nl, JSON(Data(Version(7, "blocked", "V1", "draft"))))
	api.On("GET", nl+"/7/versions/V1", JSON(Version(7, "blocked", "V1", "building")))
	a, _ := applier(api, &contracts.AzionApplicationOptionsV6{}, &contracts.ManifestV6{NetworkLists: []sdk.NetworkListRequest{networkList}})

	_, err := a.ApplyNetworkLists(context.Background())

	assert.ErrorIs(t, err, msg.ErrorV6DraftTaken)
}

const app = "/workspace/applications/1/versions/V1"

// The Application's sub-resources are reconciled by name inside the Draft, and
// rule references by name become ids.
func TestApplicationReconcilesSubResourcesAndResolvesReferences(t *testing.T) {
	api := &API{}
	api.On("GET", "/workspace/applications", JSON(Page()))
	api.On("POST", "/workspace/applications", JSON(Data(Version(1, "site", "V1", "draft"))))
	api.On("GET", app, JSON(Version(1, "site", "V1", "draft")))
	api.On("GET", app+"/device_groups", JSON(Page()))
	// "old" is no longer in the manifest; "static" is updated in place.
	api.On("GET", app+"/cache_settings", JSON(Page(map[string]any{"id": 50, "name": "static"}, map[string]any{"id": 51, "name": "old"})))
	api.On("PUT", app+"/cache_settings/50", JSON(Data(map[string]any{"id": 50, "name": "static"})))
	api.On("DELETE", app+"/cache_settings/51", JSON(map[string]any{}))
	api.On("GET", app+"/functions", JSON(Page()))
	api.On("POST", app+"/functions", JSON(Data(map[string]any{"id": 60, "name": "handler"})))
	api.On("GET", app+"/request_rules", JSON(Page()))
	api.On("POST", app+"/request_rules", JSON(Data(map[string]any{"id": 70, "name": "cache"})), JSON(Data(map[string]any{"id": 71, "name": "fn"})))
	api.On("PUT", app+"/request_rules/order", JSON(Page()))
	api.On("GET", app+"/response_rules", JSON(Page()))

	m := &contracts.ManifestV6{Applications: []contracts.ManifestV6Application{{
		Name:              "site",
		CacheSettings:     []sdk.CacheSettingRequest{{Name: "static"}},
		FunctionInstances: []contracts.ManifestV6FunctionInstance{{Name: "handler", Function: contracts.FunctionReference{Name: "fn"}}},
		Rules: []contracts.ManifestV6Rule{
			{Phase: "request", Rule: json.RawMessage(`{"name":"cache","criteria":[[{"variable":"${uri}","operator":"starts_with","conditional":"if","argument":"/"}]],"behaviors":[{"type":"set_cache_policy","attributes":{"value":"static"}}]}`)},
			// The bundler names the Function; the rule runs its instance.
			{Phase: "request", Rule: json.RawMessage(`{"name":"fn","criteria":[[{"variable":"${uri}","operator":"starts_with","conditional":"if","argument":"/"}]],"behaviors":[{"type":"run_function","attributes":{"value":"fn"}}]}`)},
		},
	}}}
	conf := &contracts.AzionApplicationOptionsV6{Functions: []contracts.AzionJsonV6Resource{{ID: 5, Name: "fn"}}}
	a, _ := applier(api, conf, m)

	pending, err := a.ApplyApplications(context.Background())

	require.NoError(t, err, api.Summary())
	require.Len(t, pending, 1)

	inst := api.Bodies("POST", app+"/functions")
	require.Len(t, inst, 1)
	assert.EqualValues(t, 5, inst[0]["function"], "the instance runs the function by id")

	rules := api.Bodies("POST", app+"/request_rules")
	require.Len(t, rules, 2)
	assert.EqualValues(t, 50, behaviorValue(rules[0]), "set_cache_policy gets the cache setting id")
	assert.EqualValues(t, 60, behaviorValue(rules[1]), "run_function gets the function instance id")

	assert.Equal(t, 1, api.Count("DELETE", app+"/cache_settings/51"), "orphaned sub-resources leave the Draft")
	order := api.Bodies("PUT", app+"/request_rules/order")
	require.Len(t, order, 1)
	assert.Equal(t, []any{float64(70), float64(71)}, order[0]["order"])
}

func behaviorValue(rule map[string]any) any {
	b := rule["behaviors"].([]any)[0].(map[string]any)
	return b["attributes"].(map[string]any)["value"]
}

func TestUnknownReferenceFails(t *testing.T) {
	raw := json.RawMessage(`{"name":"r","behaviors":[{"type":"set_cache_policy","attributes":{"value":"missing"}}]}`)
	_, err := resolveRule(raw, refs{behaviors: map[string]target{"set_cache_policy": {attribute: "value", kind: "cache setting", ids: map[string]int64{}}}})
	assert.ErrorIs(t, err, msg.ErrorV6UnknownRef)
}

func TestReferenceByIDIsKept(t *testing.T) {
	raw := json.RawMessage(`{"name":"r","behaviors":[{"type":"set_cache_policy","attributes":{"value":12}}]}`)
	out, err := resolveRule(raw, refs{behaviors: map[string]target{"set_cache_policy": {attribute: "value", ids: map[string]int64{}}}})
	require.NoError(t, err)
	assert.JSONEq(t, string(raw), string(out))
}

func TestValidate(t *testing.T) {
	v6 := func(m contracts.ManifestV6) contracts.ManifestV6 { m.Version = 6; return m }
	two := []contracts.ManifestV6Application{{Name: "a"}, {Name: "b"}}
	tests := []struct {
		name string
		m    contracts.ManifestV6
		want error
	}{
		{"valid", v6(contracts.ManifestV6{Functions: []contracts.ManifestV6Function{{Name: "f", Path: "f.js"}}}), nil},
		{"a v4 manifest", contracts.ManifestV6{Version: 4}, msg.ErrorV6WrongVersion},
		{"no version", contracts.ManifestV6{}, msg.ErrorV6WrongVersion},
		{"two applications", v6(contracts.ManifestV6{Applications: two}), msg.ErrorV6TooMany},
		{"duplicate names", v6(contracts.ManifestV6{NetworkLists: []sdk.NetworkListRequest{networkList, networkList}}), msg.ErrorV6DuplicateName},
		{"no name", v6(contracts.ManifestV6{WAFs: []contracts.ManifestV6WAF{{}}}), msg.ErrorV6NoName},
		{"function without code", v6(contracts.ManifestV6{Functions: []contracts.ManifestV6Function{{Name: "f"}}}), msg.ErrorV6FunctionCode},
		{"instance of an undeclared function", v6(contracts.ManifestV6{Applications: []contracts.ManifestV6Application{{
			Name: "a", FunctionInstances: []contracts.ManifestV6FunctionInstance{{Name: "i", Function: contracts.FunctionReference{Name: "nope"}}},
		}}}), msg.ErrorV6UnknownRef},
		{"rule with an unknown phase", v6(contracts.ManifestV6{Applications: []contracts.ManifestV6Application{{
			Name: "a", Rules: []contracts.ManifestV6Rule{{Phase: "middle", Rule: json.RawMessage(`{}`)}},
		}}}), msg.ErrorV6InvalidRule},
		{"binding to an undeclared environment", v6(contracts.ManifestV6{
			Deployments: []sdk.DeploymentCreateRequest{{Name: "d", BindingPolicy: "STRICT"}},
			Workloads:   []contracts.ManifestV6Workload{{Name: "w", Bindings: []contracts.ManifestV6Binding{{EnvironmentName: "x", DeploymentName: "d"}}}},
		}), msg.ErrorV6UnknownBinding},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.m)
			if tt.want == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.want)
		})
	}
}

// What the bundler writes with --config-version 6 is a valid manifest.
func TestBundlerOutputIsValid(t *testing.T) {
	m, err := Read("testdata/bundler666_full.json")
	require.NoError(t, err)
	assert.Equal(t, 6, m.Version)
}

// Firewall criteria name a Network List, set_waf names a WAF and run_function
// names a Function; all become ids.
func TestFirewallReferences(t *testing.T) {
	raw := json.RawMessage(`{"name":"r","criteria":[[{"variable":"${network}","operator":"is_in_list","conditional":"if","argument":"blocked"}]],
		"behaviors":[{"type":"set_waf","attributes":{"waf_id":"my-waf","mode":"blocking"}}]}`)
	out, err := resolveRule(raw, refs{
		behaviors:    map[string]target{"set_waf": {attribute: "waf_id", kind: "waf", ids: map[string]int64{"my-waf": 9}}},
		networkLists: map[string]int64{"blocked": 3},
	})
	require.NoError(t, err)
	var rule map[string]any
	require.NoError(t, json.Unmarshal(out, &rule))
	assert.EqualValues(t, 3, rule["criteria"].([]any)[0].([]any)[0].(map[string]any)["argument"])
	assert.EqualValues(t, 9, rule["behaviors"].([]any)[0].(map[string]any)["attributes"].(map[string]any)["waf_id"])

	_, err = resolveRule(raw, refs{networkLists: map[string]int64{}})
	assert.ErrorIs(t, err, msg.ErrorV6UnknownRef)
}

func TestRunFunctionNamesTheFunction(t *testing.T) {
	instances := []contracts.ManifestV6FunctionInstance{
		{Name: "a", Function: contracts.FunctionReference{Name: "fn"}},
		{Name: "b", Function: contracts.FunctionReference{Name: "fn"}},
	}
	got := runFunctionIDs(instances, map[string]int64{"a": 1, "b": 2})
	assert.Equal(t, map[string]int64{"fn": 1, "a": 1, "b": 2}, got, "the first instance of a Function wins; each instance is reachable by its own name")
}

func TestCustomPageConnectorIsResolved(t *testing.T) {
	raw := json.RawMessage(`{"name":"p","pages":[{"code":"404","page":{"type":"page_connector","attributes":{"connector":"my-origin","uri":"/404.html"}}}]}`)
	out, err := resolveCustomPage(raw, map[string]int64{"my-origin": 11})
	require.NoError(t, err)
	var cp sdk.CustomPageRequest
	require.NoError(t, json.Unmarshal(out, &cp))
	assert.Equal(t, "p", cp.Name)
	assert.Contains(t, string(out), `"connector":11`)

	_, err = resolveCustomPage(raw, map[string]int64{})
	assert.ErrorIs(t, err, msg.ErrorV6UnknownRef)
}

func TestOrphansAreReportedNotDeleted(t *testing.T) {
	conf := &contracts.AzionApplicationOptionsV6{Connectors: []contracts.AzionJsonV6Resource{{ID: 3, Name: "gone"}}}
	a, _ := applier(&API{}, conf, &contracts.ManifestV6{})
	orphans := a.Orphans()
	require.Len(t, orphans, 1)
	assert.Contains(t, orphans[0], `connector "gone"`)
}

func TestFunctionCodeComesFromItsFile(t *testing.T) {
	api := &API{}
	api.On("GET", "/workspace/functions", JSON(Page()))
	api.On("POST", "/workspace/functions", JSON(Data(Version(4, "fn", "V1", "draft"))))
	api.On("GET", "/workspace/functions/4/versions/V1", JSON(Version(4, "fn", "V1", "draft")))
	a, _ := applier(api, &contracts.AzionApplicationOptionsV6{}, &contracts.ManifestV6{Functions: []contracts.ManifestV6Function{{Name: "fn", Path: ".edge/worker.js"}}})
	a.ReadFile = func(path string) ([]byte, error) {
		if path != ".edge/worker.js" {
			return nil, errors.New("unexpected path " + path)
		}
		return []byte("export default {}"), nil
	}

	_, err := a.ApplyFunctions(context.Background())

	require.NoError(t, err, api.Summary())
	body := api.Bodies("POST", "/workspace/functions")
	require.Len(t, body, 1)
	assert.Equal(t, "export default {}", body[0]["code"])
}

// Every rule, firewall rule and custom page the bundler writes decodes into
// its SDK request once the names are resolved: the SDK's decoders are strict,
// so this is where a format mismatch would show.
func TestBundlerOutputDecodesIntoSDKRequests(t *testing.T) {
	m, err := Read("testdata/bundler666_full.json")
	require.NoError(t, err)

	app := m.Applications[0]
	request, response, err := splitRules(app.Rules)
	require.NoError(t, err)
	appRefs := refs{behaviors: map[string]target{
		"run_function":     {attribute: "value", ids: runFunctionIDs(app.FunctionInstances, map[string]int64{"my-instance": 60})},
		"set_cache_policy": {attribute: "value", ids: map[string]int64{"my-cache": 50}},
		"set_connector":    {attribute: "value", ids: map[string]int64{"my-connector": 70, "my-origin": 71}},
	}}
	req, err := decodeRules[sdk.RequestPhaseRuleRequest](request, appRefs)
	require.NoError(t, err)
	assert.Len(t, req, 2)
	resp, err := decodeRules[sdk.ResponsePhaseRuleRequest](response, appRefs)
	require.NoError(t, err)
	assert.Len(t, resp, 1)

	fw := m.Firewalls[0]
	fwRules, err := decodeRules[sdk.FirewallRuleRequest](fw.Rules, refs{
		behaviors:    map[string]target{"set_waf": {attribute: "waf_id", ids: map[string]int64{"my-waf": 9}}},
		networkLists: map[string]int64{"blocked": 3},
	})
	require.NoError(t, err)
	assert.Len(t, fwRules, 2)

	page, err := resolveCustomPage(m.CustomPages[0], map[string]int64{"my-origin": 71})
	require.NoError(t, err)
	var cp sdk.CustomPageRequest
	require.NoError(t, json.Unmarshal(page, &cp))
	assert.Equal(t, "my-pages", cp.Name)
}
