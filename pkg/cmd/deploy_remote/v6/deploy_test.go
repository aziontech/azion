package v6

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	manifestmsg "github.com/aziontech/azion-cli/messages/manifest"
	apiDeployments "github.com/aziontech/azion-cli/pkg/api/v6/deployments"
	apiEnvironments "github.com/aziontech/azion-cli/pkg/api/v6/environments"
	. "github.com/aziontech/azion-cli/pkg/api/v6/v6test"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	apiWorkloads "github.com/aziontech/azion-cli/pkg/api/v6/workloads"
	deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/pkg/poll"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	apiApplications "github.com/aziontech/azion-cli/pkg/api/v6/applications"
	apiConnectors "github.com/aziontech/azion-cli/pkg/api/v6/connectors"
	apiCustomPages "github.com/aziontech/azion-cli/pkg/api/v6/custompages"
	apiFirewalls "github.com/aziontech/azion-cli/pkg/api/v6/firewalls"
	apiFunctions "github.com/aziontech/azion-cli/pkg/api/v6/functions"
	apiNetworkLists "github.com/aziontech/azion-cli/pkg/api/v6/networklists"
	apiWafs "github.com/aziontech/azion-cli/pkg/api/v6/wafs"
)

func TestMain(m *testing.M) {
	logger.New(zapcore.DebugLevel)
	os.Exit(m.Run())
}

// instant is a clock where every wait between reads fires at once and the
// deadline of a wait never does.
func instant(d time.Duration) <-chan time.Time {
	if d >= poll.DefaultTimeout {
		return make(chan time.Time)
	}
	ch := make(chan time.Time, 1)
	ch <- time.Time{}
	return ch
}

func fakeClients(api *API) *Clients {
	c := api.Client()
	return &Clients{
		Resources: &manifestv6.Clients{
			Applications: apiApplications.NewClient(c, URL, "t"),
			Functions:    apiFunctions.NewClient(c, URL, "t"),
			Connectors:   apiConnectors.NewClient(c, URL, "t"),
			Firewalls:    apiFirewalls.NewClient(c, URL, "t"),
			CustomPages:  apiCustomPages.NewClient(c, URL, "t"),
			NetworkLists: apiNetworkLists.NewClient(c, URL, "t"),
			WAFs:         apiWafs.NewClient(c, URL, "t"),
		},
		Environments: apiEnvironments.NewClient(c, URL, "t"),
		Deployments:  apiDeployments.NewClient(c, URL, "t"),
		Workloads:    apiWorkloads.NewClient(c, URL, "t"),
	}
}

func factory() (*cmdutil.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &cmdutil.Factory{
		IOStreams: &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}},
		Config:    viper.New(),
	}, out
}

// project writes a manifest into a fresh working directory.
func project(t *testing.T, m map[string]any) {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".edge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".edge", "worker.js"), []byte("export default {}"), 0644))
	data, err := json.Marshal(m)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".edge", "manifest.json"), data, 0644))
}

// siteManifest is a manifest in the format the bundler writes with
// --config-version 6.
var siteManifest = map[string]any{
	"version":   6,
	"functions": []any{map[string]any{"name": "fn", "path": "./worker.js"}},
	"applications": []any{map[string]any{
		"name":                "site",
		"functions_instances": []any{map[string]any{"name": "handler", "function": "fn"}},
		"rules": []any{map[string]any{
			"phase": "request",
			"rule": map[string]any{
				"name":      "run",
				"criteria":  [][]any{{map[string]any{"variable": "${uri}", "operator": "starts_with", "conditional": "if", "argument": "/"}}},
				"behaviors": []any{map[string]any{"type": "run_function", "attributes": map[string]any{"value": "fn"}}},
			},
		}},
	}},
	"environments": []any{map[string]any{"name": "Production", "deployment_policy": "single_version"}},
	"deployments":  []any{map[string]any{"name": "site-deploy", "binding_policy": "FLEXIBLE"}},
	"workloads": []any{map[string]any{
		"name": "site",
		"bindings": []any{map[string]any{
			"environment_name": "Production", "deployment_name": "site-deploy",
			"domains": []any{map[string]any{"domain": "www.example.com", "certificate": nil}},
		}},
	}},
}

const (
	fnV  = "/workspace/functions/4/versions/F1"
	appV = "/workspace/applications/1/versions/A1"
	rel  = "/deployments/DEP00001/releases/REL00001"
)

func release(state, role string) Response {
	return JSON(Data(map[string]any{"id": "REL00001", "deployment_id": "DEP00001", "state": state, "traffic_role": role}))
}

func hierarchy(appVersion string) Response {
	return JSON(Data(map[string]any{
		"deployment_id": "DEP00001", "active_release_id": "REL00001",
		"resources": map[string]any{
			"application": map[string]any{"resource_type": "application", "resource_id": "1", "resource_version": appVersion, "product_version": nil},
			"firewall":    nil, "custom_page": nil,
			"function": []any{map[string]any{"resource_type": "function", "resource_id": "4", "resource_version": "F1", "product_version": nil}},
			"waf":      []any{}, "connector": []any{}, "network_list": []any{},
		},
	}))
}

// siteAPI answers a first deploy of siteManifest to Production.
func siteAPI() *API {
	api := &API{}
	api.On("GET", "/workspace/functions", JSON(Page()))
	api.On("POST", "/workspace/functions", JSON(Data(Version(4, "fn", "F1", "draft"))))
	api.On("GET", fnV, JSON(Version(4, "fn", "F1", "draft")), JSON(Version(4, "fn", "F1", "queued")), JSON(Version(4, "fn", "F1", "ready")))
	api.On("POST", fnV+"/build", Status(202))

	api.On("GET", "/workspace/applications", JSON(Page()))
	api.On("POST", "/workspace/applications", JSON(Data(Version(1, "site", "A1", "draft"))))
	api.On("GET", appV, JSON(Version(1, "site", "A1", "draft")), JSON(Version(1, "site", "A1", "building")), JSON(Version(1, "site", "A1", "ready")))
	for _, sub := range []string{"device_groups", "cache_settings", "functions", "request_rules", "response_rules"} {
		api.On("GET", appV+"/"+sub, JSON(Page()))
	}
	api.On("POST", appV+"/functions", JSON(Data(map[string]any{"id": 60, "name": "handler"})))
	api.On("POST", appV+"/request_rules", JSON(Data(map[string]any{"id": 70, "name": "run"})))
	api.On("POST", appV+"/build", Status(202))

	api.On("GET", "/environments", JSON(Page(map[string]any{"id": "ENV00001", "name": "Production", "deployment_policy": "single_version"})))
	api.On("GET", "/deployments", JSON(Page()))
	api.On("POST", "/deployments", JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "binding_policy": "FLEXIBLE", "deployment_policy": "single_version", "version_state": "ready"})))

	api.On("POST", "/deployments/DEP00001/releases", JSON(Data(map[string]any{"id": "REL00001", "state": "draft"})))
	api.On("POST", rel+"/build", JSON(Data(map[string]any{"id": "REL00001", "state": "queued"})))
	api.On("GET", rel, release("building", "INACTIVE"), release("ready", "INACTIVE"), release("ready", "ACTIVE"))
	api.On("POST", rel+"/activate", JSON(Data(map[string]any{"id": "REL00001", "state": "ready"})))
	api.On("GET", "/deployments/DEP00001/resource_hierarchy", hierarchy("A1"))

	workload := map[string]any{
		"id": 30, "name": "site", "workload_domain": "abc.azion.app",
		"bindings": []any{map[string]any{"environment_id": "ENV00001", "deployment_id": "DEP00001", "domains": []any{map[string]any{"domain": "www.example.com"}}}},
	}
	api.On("GET", "/workspace/workloads", JSON(Page()))
	api.On("POST", "/workspace/workloads", JSON(Data(workload)))
	api.On("GET", "/workspace/workloads/30", JSON(Data(workload)))
	// Read again right before the Workload binds them.
	api.On("GET", "/environments/ENV00001", JSON(Data(map[string]any{"id": "ENV00001", "name": "Production", "deployment_policy": "single_version", "version_state": "ready"})))
	api.On("GET", "/deployments/DEP00001", JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "deployment_policy": "single_version", "version_id": "DV1", "version_state": "ready"})))
	return api
}

func runner(t *testing.T, api *API, opts Options) (*Runner, *contracts.AzionApplicationOptionsV6, *bytes.Buffer) {
	f, out := factory()
	conf := &contracts.AzionApplicationOptionsV6{Name: "site", NotFirstRun: true}
	opts.SkipBuild, opts.SkipFramework = true, true
	cmd := deploy.NewDeployCmd(f)
	cmd.SkipBuild, cmd.SkipFramework = true, true
	return &Runner{
		F: f, Cmd: cmd.ForV6(), Clients: fakeClients(api), Opts: opts,
		ReadConf:  func(string) (*contracts.AzionApplicationOptionsV6, error) { return conf, nil },
		WriteConf: func(*contracts.AzionApplicationOptionsV6, string) error { return nil },
		NewPrefix: func() string { return "P1" },
		WorkDir:   os.Getwd,
		Poll:      poll.Options{After: instant},
	}, conf, out
}

func TestDeployEndToEnd(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	r, conf, out := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	assert.Equal(t, []contracts.AzionJsonV6Resource{{ID: 4, Name: "fn", VersionID: "F1", VersionState: "ready"}}, conf.Functions)
	// Each build that reaches ready says so.
	assert.Contains(t, out.String(), `function "fn" built successfully: version F1 is ready`)
	assert.Contains(t, out.String(), `application "site" built successfully: version A1 is ready`)
	assert.Equal(t, []contracts.AzionJsonV6Resource{{ID: 1, Name: "site", VersionID: "A1", VersionState: "ready"}}, conf.Applications)
	assert.Equal(t, []contracts.AzionJsonV6Environment{{ID: "ENV00001", Name: "Production", Policy: "single_version"}}, conf.Environments)
	assert.Equal(t, []contracts.AzionJsonV6Deployment{{
		ID: "DEP00001", Name: "site-deploy", Policy: "single_version",
		Release: contracts.AzionJsonV6Release{ID: "REL00001", State: "ready", Activated: true},
	}}, conf.Deployments)
	assert.Equal(t, []contracts.AzionJsonV6Workload{{ID: 30, Name: "site", Url: "https://www.example.com"}}, conf.Workloads)

	// The Environment exists, so it is not created; the Deployment is, with
	// the manifest's binding policy and the Environment's policy.
	assert.Zero(t, api.Count("POST", "/environments"))
	dep := api.Bodies("POST", "/deployments")[0]
	assert.Equal(t, "FLEXIBLE", dep["binding_policy"])
	assert.Equal(t, "single_version", dep["deployment_policy"])

	// The release names the version this run built of what it serves and of
	// its dependencies (D7): the API does not add a missing dependency.
	body := api.Bodies("POST", "/deployments/DEP00001/releases")[0]
	assert.Equal(t, []any{
		map[string]any{"resource_id": float64(1), "resource_type": "application", "version_id": "A1"},
		map[string]any{"resource_id": float64(4), "resource_type": "function", "version_id": "F1"},
	}, body["resources"])
	assert.Equal(t, "cli", body["origin"].(map[string]any)["type"])

	// Dependencies are built before what depends on them.
	assert.Less(t, callIndex(api, "POST", fnV+"/build"), callIndex(api, "POST", appV+"/build"))

	// run_function names the Function; the rule runs the instance of it.
	rule := api.Bodies("POST", appV+"/request_rules")[0]
	assert.EqualValues(t, 60, rule["behaviors"].([]any)[0].(map[string]any)["attributes"].(map[string]any)["value"])

	w := api.Bodies("POST", "/workspace/workloads")[0]
	assert.Equal(t, []any{map[string]any{"environment_id": "ENV00001", "deployment_id": "DEP00001", "domains": []any{map[string]any{"domain": "www.example.com"}}}}, w["bindings"])

	assert.Contains(t, out.String(), "queued → ready")
	assert.Contains(t, out.String(), msg.V6DeploySuccessful)
}

func callIndex(api *API, method, path string) int {
	for i, c := range api.Calls {
		if c.Method == method && c.Path == path {
			return i
		}
	}
	return -1
}

func TestNoActivateStopsAfterTheReleaseBuild(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	r, conf, out := runner(t, api, Options{NoActivate: true})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	assert.Zero(t, api.Count("POST", rel+"/activate"))
	assert.Zero(t, api.Count("GET", "/deployments/DEP00001/resource_hierarchy"))
	assert.False(t, conf.Deployments[0].Release.Activated)
	assert.Zero(t, api.Count("POST", "/workspace/workloads"), "no Workload change: nothing new is live")
	assert.Contains(t, out.String(), "was built but not activated")
}

// An Environment the manifest declares and the account lacks is created.
func TestMissingEnvironmentIsCreated(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/environments", JSON(Page()))
	api.On("POST", "/environments", JSON(Data(map[string]any{"id": "ENV00001", "name": "Production", "deployment_policy": "single_version"})))
	r, _, _ := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	assert.Equal(t, "single_version", api.Bodies("POST", "/environments")[0]["deployment_policy"])
}

// An Environment's deployment policy cannot change; a different one in the
// manifest stops the deploy instead of being ignored.
func TestEnvironmentPolicyMismatchFails(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/environments", JSON(Page(map[string]any{"id": "ENV00001", "name": "Production", "deployment_policy": "versioned_urls"})))
	r, _, _ := runner(t, api, Options{})

	assert.ErrorIs(t, r.Run(context.Background()), msg.ErrorV6PolicyImmutable)
}

// --environment releases only the Deployments bound to that Environment and
// changes only its bindings.
func TestEnvironmentFlagFilters(t *testing.T) {
	m := map[string]any{}
	for k, v := range siteManifest {
		m[k] = v
	}
	m["environments"] = []any{
		map[string]any{"name": "Production", "deployment_policy": "single_version"},
		map[string]any{"name": "Staging", "deployment_policy": "single_version"},
	}
	m["deployments"] = []any{
		map[string]any{"name": "site-deploy", "binding_policy": "FLEXIBLE"},
		map[string]any{"name": "site-staging", "binding_policy": "FLEXIBLE"},
	}
	m["workloads"] = []any{map[string]any{"name": "site", "bindings": []any{
		map[string]any{"environment_name": "Production", "deployment_name": "site-deploy"},
		map[string]any{"environment_name": "Staging", "deployment_name": "site-staging"},
	}}}
	project(t, m)
	api := siteAPI()
	r, conf, _ := runner(t, api, Options{Environment: "Production"})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	require.Len(t, api.Bodies("POST", "/deployments"), 1)
	assert.Equal(t, "site-deploy", api.Bodies("POST", "/deployments")[0]["name"])
	require.Len(t, conf.Deployments, 1)
	w := api.Bodies("POST", "/workspace/workloads")[0]
	assert.Len(t, w["bindings"], 1)
}

// A manifest without Environments and Deployments still deploys: to
// Production, which must exist, with a Deployment named after the project.
func TestManifestWithoutEnvironmentsUsesDefaults(t *testing.T) {
	m := map[string]any{}
	for k, v := range siteManifest {
		if k != "environments" && k != "deployments" && k != "workloads" {
			m[k] = v
		}
	}
	m["workloads"] = []any{map[string]any{"name": "site"}}
	project(t, m)
	api := siteAPI()
	r, _, _ := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	dep := api.Bodies("POST", "/deployments")[0]
	assert.Equal(t, "site", dep["name"])
	assert.Equal(t, "FLEXIBLE", dep["binding_policy"], "a STRICT Deployment can only take its first release")

	api = siteAPI()
	api.On("GET", "/environments", JSON(Page()))
	r, _, _ = runner(t, api, Options{})
	assert.ErrorIs(t, r.Run(context.Background()), msg.ErrorV6EnvironmentNotFound, "the default Environment is never created")
}

// One failed build does not stop the others of its level, and the error lists
// every failure.
func TestBuildLevelWaitsForAllAndJoinsFailures(t *testing.T) {
	api := &API{}
	for _, id := range []string{"1", "2", "3"} {
		api.On("POST", "/workspace/network_lists/"+id+"/versions/V/build", Status(202))
	}
	api.On("GET", "/workspace/network_lists/1/versions/V", JSON(Version(1, "a", "V", "building")), JSON(Version(1, "a", "V", "error")))
	api.On("GET", "/workspace/network_lists/2/versions/V", JSON(Version(2, "b", "V", "queued")), JSON(Version(2, "b", "V", "ready")))
	api.On("GET", "/workspace/network_lists/3/versions/V", JSON(Version(3, "c", "V", "canceled")))

	f, _ := factory()
	conf := &contracts.AzionApplicationOptionsV6{NetworkLists: []contracts.AzionJsonV6Resource{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}, {ID: 3, Name: "c"}}}
	c := fakeClients(api).Resources.NetworkLists
	level := []manifestv6.Pending{}
	for _, r := range conf.NetworkLists {
		level = append(level, manifestv6.Pending{Kind: "network list", Name: r.Name, ID: r.ID, VersionID: "V", Client: c, Record: &conf.NetworkLists})
	}
	s := &DeployStateV6{F: f, Conf: conf, Poll: poll.Options{After: instant}, msgs: &[]string{}}

	err := s.buildLevel(context.Background(), level)

	require.Error(t, err)
	assert.ErrorIs(t, err, version.ErrBuildFailed)
	var be *version.BuildError
	require.True(t, errors.As(err, &be))
	assert.Contains(t, err.Error(), `network list "a"`)
	assert.Contains(t, err.Error(), `network list "c"`)
	assert.NotContains(t, err.Error(), `network list "b"`)
	assert.Equal(t, []string{"error", "ready", "canceled"}, []string{conf.NetworkLists[0].VersionState, conf.NetworkLists[1].VersionState, conf.NetworkLists[2].VersionState})
}

// Ctrl-C during a wait asks the API to stop the build this run started.
func TestInterruptCancelsTheBuild(t *testing.T) {
	api := &API{}
	api.On("POST", "/workspace/network_lists/1/versions/V/build", Status(202))
	api.On("GET", "/workspace/network_lists/1/versions/V", JSON(Version(1, "a", "V", "building")))
	api.On("POST", "/workspace/network_lists/1/versions/V/cancel", Status(202))

	ctx, cancel := context.WithCancel(context.Background())
	f, _ := factory()
	conf := &contracts.AzionApplicationOptionsV6{NetworkLists: []contracts.AzionJsonV6Resource{{ID: 1, Name: "a"}}}
	s := &DeployStateV6{F: f, Conf: conf, msgs: &[]string{}, Poll: poll.Options{After: func(d time.Duration) <-chan time.Time {
		if d < poll.DefaultTimeout {
			cancel() // the user presses Ctrl-C while the CLI waits
		}
		return make(chan time.Time)
	}}}
	p := manifestv6.Pending{Kind: "network list", Name: "a", ID: 1, VersionID: "V", Client: fakeClients(api).Resources.NetworkLists, Record: &conf.NetworkLists}

	err := s.buildLevel(ctx, []manifestv6.Pending{p})

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, api.Count("POST", "/workspace/network_lists/1/versions/V/cancel"))
}

// A release whose build fails is deleted, so failed runs do not use up the
// Deployment's release Drafts, and nothing goes live.
func TestFailedReleaseIsDeleted(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", rel, release("building", "INACTIVE"), JSON(Data(map[string]any{"id": "REL00001", "state": "error", "trace_id": "T-1"})))
	api.On("DELETE", rel, Status(204))
	r, conf, _ := runner(t, api, Options{})

	err := r.Run(context.Background())

	assert.ErrorIs(t, err, msg.ErrorV6ReleaseFailed)
	assert.Contains(t, err.Error(), "trace id T-1")
	assert.Equal(t, 1, api.Count("DELETE", rel))
	assert.Zero(t, api.Count("POST", rel+"/activate"))
	assert.Equal(t, contracts.AzionJsonV6Release{}, conf.Deployments[0].Release)
}

// The release names the version this run built, so serving another version
// of a released resource is a difference, not a warning.
func TestHierarchyWithAnotherVersionFails(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/deployments/DEP00001/resource_hierarchy", hierarchy("A9"))
	r, _, _ := runner(t, api, Options{})

	err := r.Run(context.Background())

	assert.ErrorIs(t, err, msg.ErrorV6VerifyFailed)
	assert.Contains(t, err.Error(), `application "site" is served at version A9, not A1`)
}

func TestVerifyFailsWhenTheReleaseIsNotActive(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/deployments/DEP00001/resource_hierarchy", JSON(Data(map[string]any{
		"deployment_id": "DEP00001", "active_release_id": "OTHER001",
		"resources": map[string]any{"application": nil, "firewall": nil, "custom_page": nil, "function": []any{}, "waf": []any{}, "connector": []any{}, "network_list": []any{}},
	})))
	r, _, out := runner(t, api, Options{})

	err := r.Run(context.Background())

	assert.ErrorIs(t, err, msg.ErrorV6VerifyFailed)
	assert.Contains(t, err.Error(), `the active release is "OTHER001"`)
	assert.Contains(t, err.Error(), `application "site" (id 1) is not served`)
	assert.Contains(t, out.String(), "The new release is live")
}

func TestRecordedDeploymentIsReused(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/deployments/DEP00001", JSON(Data(map[string]any{"id": "DEP00001", "name": "mine", "deployment_policy": "single_version", "version_state": "ready"})))
	r, conf, _ := runner(t, api, Options{})
	conf.Deployments = []contracts.AzionJsonV6Deployment{{Name: "site-deploy", ID: "DEP00001"}}

	require.NoError(t, r.Run(context.Background()), api.Summary())

	assert.Zero(t, api.Count("POST", "/deployments"))
}

func TestDryRunCallsNothing(t *testing.T) {
	project(t, siteManifest)
	api := &API{}
	r, _, out := runner(t, api, Options{DryRun: true})

	require.NoError(t, r.Run(context.Background()))

	assert.Empty(t, api.Calls)
	assert.Contains(t, out.String(), "Functions: 1 Draft(s)")
	assert.Contains(t, out.String(), "Environments: Production")
	assert.Contains(t, out.String(), "Deployments: site-deploy")
	assert.Contains(t, out.String(), "each at the version this run builds")
}

// A dry run writes nothing, azion.json included.
func TestDryRunDoesNotWriteAzionJson(t *testing.T) {
	project(t, siteManifest)
	r, _, _ := runner(t, &API{}, Options{DryRun: true})
	r.WriteConf = func(*contracts.AzionApplicationOptionsV6, string) error {
		t.Error("a dry run saved azion.json")
		return nil
	}

	require.NoError(t, r.Run(context.Background()))
}

// The manifest a v4 build left behind does not stop a dry run: the deploy
// builds the manifest again with the v6 bundler before reading it.
func TestDryRunWithAV4Manifest(t *testing.T) {
	project(t, map[string]any{"application": []any{map[string]any{"name": "site"}}})
	api := &API{}
	r, _, out := runner(t, api, Options{DryRun: true})

	require.NoError(t, r.Run(context.Background()))

	assert.Empty(t, api.Calls)
	assert.Contains(t, out.String(), "The current manifest cannot be used")
	assert.Contains(t, out.String(), "Dry run: nothing was sent to the API")
}

// A deploy to one Environment keeps the Workload's bindings to the others.
func TestWorkloadKeepsOtherEnvironmentsBindings(t *testing.T) {
	existing := &sdk.Workload{Bindings: []sdk.WorkloadBinding{
		{EnvironmentId: "STAGING1", DeploymentId: "DEPSTAGE"},
		{EnvironmentId: "ENV00001", DeploymentId: "OLD"},
	}}
	existing.SetName("site")
	ours := []sdk.WorkloadBindingRequest{{EnvironmentId: "ENV00001", DeploymentId: "DEP00001"}}

	req, err := workloadRequest(existing, contracts.ManifestV6Workload{Name: "site"}, ours)

	require.NoError(t, err)
	require.Len(t, req.Bindings, 2)
	assert.Equal(t, "STAGING1", req.Bindings[0].EnvironmentId)
	assert.Equal(t, "DEPSTAGE", req.Bindings[0].DeploymentId)
	assert.Equal(t, "DEP00001", req.Bindings[1].DeploymentId)
}

// Step names are a contract with the timing summary: a renamed step must not
// silently disappear from it.
func TestEveryStepHasATimingLabel(t *testing.T) {
	for _, step := range (V6DeployPipeline{}).Steps() {
		assert.NotEmpty(t, stepLabels[step.Name], "step %q has no timing label", step.Name)
	}
}

// A first run generates the manifest again when the one on disk is not v6
// (a v4 build or init left it), and keeps a v6 one.
func TestFirstRunRegeneratesANonV6Manifest(t *testing.T) {
	var generate pipeline.Step[DeployStateV6]
	for _, step := range (V6DeployPipeline{}).Steps() {
		if step.Name == StepGenerateManifest {
			generate = step
		}
	}
	first := &DeployStateV6{Conf: &contracts.AzionApplicationOptionsV6{}}

	project(t, map[string]any{"application": []any{}})
	assert.True(t, generate.Applies(first), "a v4 manifest is generated again")

	project(t, siteManifest)
	assert.False(t, generate.Applies(first), "a v6 manifest is kept")

	require.NoError(t, os.Remove(".edge/manifest.json"))
	assert.True(t, generate.Applies(first), "a missing manifest is generated")

	later := &DeployStateV6{Conf: &contracts.AzionApplicationOptionsV6{NotFirstRun: true}}
	assert.False(t, generate.Applies(later), "later runs build instead")
}

// A panic in a step fails the deploy with that step's error instead of
// crashing the CLI.
func TestStepPanicIsAnError(t *testing.T) {
	run := guarded(StepReleaseCreate, func(context.Context, *DeployStateV6) error {
		var o version.Object
		_ = o.GetVersionId() // nil interface
		return nil
	})

	var err error
	assert.NotPanics(t, func() { err = run(context.Background(), &DeployStateV6{}) })
	require.Error(t, err)
	assert.ErrorIs(t, err, manifestmsg.ErrorV6Internal)
	assert.Contains(t, err.Error(), "Release create")
}

// templateManifest is what a v6 template's manifest looks like before the
// first deploy replaces its placeholders: no Deployment, and a Workload bound
// to the placeholder names.
func templateManifest(env, deployment, workload string) map[string]any {
	return map[string]any{
		"version":      6,
		"environments": []any{map[string]any{"name": env, "deployment_policy": "single_version"}},
		"workloads": []any{map[string]any{
			"name":     workload,
			"bindings": []any{map[string]any{"environment_name": env, "deployment_name": deployment}},
		}},
	}
}

func stepRun(t *testing.T, name string) func(context.Context, *DeployStateV6) error {
	t.Helper()
	for _, step := range (V6DeployPipeline{}).Steps() {
		if step.Name == name {
			return step.Run
		}
	}
	t.Fatalf("no step %s", name)
	return nil
}

// On a first deploy the manifest is read once before BundlerInit replaces the
// template's placeholders. That read must not fail on them; the read after
// the replace validates, with the Deployment the CLI adds by default.
func TestFirstRunReadsTheTemplateBeforeReplace(t *testing.T) {
	s := &DeployStateV6{
		Clients: &Clients{},
		Conf:    &contracts.AzionApplicationOptionsV6{Name: "site"},
	}

	project(t, templateManifest("$ENVIRONMENT_NAME", "$DEPLOYMENT_NAME", "$WORKLOAD_NAME"))
	require.NoError(t, stepRun(t, StepReadManifest)(context.Background(), s))

	// What BundlerInit leaves: the placeholders replaced, the manifest built
	// again.
	data, err := json.Marshal(templateManifest("Production", "site", "site"))
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(".edge/manifest.json", data, 0644))
	require.NoError(t, stepRun(t, StepReadManifestAgain)(context.Background(), s))

	require.Len(t, s.Manifest.Deployments, 1)
	assert.Equal(t, "site", s.Manifest.Deployments[0].Name)
}

// Placeholders left after the replace are still an error.
func TestPlaceholdersAfterReplaceFail(t *testing.T) {
	s := &DeployStateV6{
		Clients: &Clients{},
		Conf:    &contracts.AzionApplicationOptionsV6{Name: "site"},
	}
	project(t, templateManifest("$ENVIRONMENT_NAME", "$DEPLOYMENT_NAME", "$WORKLOAD_NAME"))

	err := stepRun(t, StepReadManifestAgain)(context.Background(), s)

	assert.ErrorIs(t, err, manifestmsg.ErrorV6UnknownBinding)
}

// Later runs have no replace step, so their first read validates.
func TestLaterRunsValidateTheFirstRead(t *testing.T) {
	s := &DeployStateV6{
		Clients: &Clients{},
		Conf:    &contracts.AzionApplicationOptionsV6{Name: "site", NotFirstRun: true},
	}
	project(t, templateManifest("$ENVIRONMENT_NAME", "$DEPLOYMENT_NAME", "$WORKLOAD_NAME"))

	err := stepRun(t, StepReadManifest)(context.Background(), s)

	assert.ErrorIs(t, err, manifestmsg.ErrorV6UnknownBinding)
}

// A release the API refuses to build (here, one missing a dependency) is
// deleted too, and the error says which release failed and why.
func TestRefusedReleaseBuildIsDeleted(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("POST", rel+"/build", Status(400, map[string]any{"errors": []any{map[string]any{
		"code": "43010", "title": "Dependency Mismatch", "status": "400",
		"detail": "Validation failed on field 'resources': missing required subresource dependency of type: connector",
	}}}))
	api.On("DELETE", rel, Status(204))
	r, conf, _ := runner(t, api, Options{})

	err := r.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "REL00001")
	assert.Contains(t, err.Error(), "missing required subresource dependency")
	assert.Equal(t, 1, api.Count("DELETE", rel))
	assert.Zero(t, api.Count("POST", rel+"/activate"))
	assert.Equal(t, contracts.AzionJsonV6Release{}, conf.Deployments[0].Release)
}

// The Workload is saved only once the Environment and the Deployment it binds
// are queued, building or ready: the API refuses a binding to a draft one, and
// a Deployment leaves draft only after its release is activated.
func TestWorkloadWaitsUntilItsBindingsAreBindable(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/deployments/DEP00001",
		// Read by the activation, then polled before the Workload binds it.
		JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "version_id": "DV1", "version_state": "draft"})),
		JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "version_state": "draft"})),
		JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "version_state": "DRAFT"})),
		JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "version_state": "building"})))
	r, _, out := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	assert.Equal(t, 4, api.Count("GET", "/deployments/DEP00001"), "one read for the activation, then polled until building")
	assert.Equal(t, 1, api.Count("GET", "/environments/ENV00001"), "already ready: read once")
	assert.Equal(t, 1, api.Count("POST", "/workspace/workloads"))
	assert.Contains(t, out.String(), `Waiting for Deployment "site-deploy" (DEP00001) before binding it to the Workload: it is draft`)
}

func TestWorkloadIsNotSavedWhenABindingFailed(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/environments/ENV00001", JSON(Data(map[string]any{"id": "ENV00001", "name": "Production", "version_state": "error"})))
	r, _, _ := runner(t, api, Options{})

	err := r.Run(context.Background())

	require.ErrorIs(t, err, msg.ErrorV6NotBindable)
	assert.Contains(t, err.Error(), `Environment "Production" (ENV00001) is error`)
	assert.Zero(t, api.Count("POST", "/workspace/workloads"))
}

// The activation names the Deployment's current version, read right before,
// as the API's own deploy scripts do.
func TestActivationSendsTheDeploymentVersion(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	r, _, _ := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	act := api.Bodies("POST", rel+"/activate")
	require.Len(t, act, 1)
	assert.Equal(t, "DV1", act[0]["deployment_version_id"])
}

// The field is optional: a failed read activates without it rather than
// failing the deploy.
func TestActivationWithoutTheDeploymentVersionWhenItCannotBeRead(t *testing.T) {
	project(t, siteManifest)
	api := siteAPI()
	api.On("GET", "/deployments/DEP00001", Status(500), JSON(Data(map[string]any{"id": "DEP00001", "name": "site-deploy", "version_state": "ready"})))
	r, _, _ := runner(t, api, Options{})

	require.NoError(t, r.Run(context.Background()), api.Summary())

	act := api.Bodies("POST", rel+"/activate")
	require.Len(t, act, 1)
	assert.NotContains(t, act[0], "deployment_version_id")
}
