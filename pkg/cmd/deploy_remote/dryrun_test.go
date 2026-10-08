package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// refusingTransport fails the test if the dry run tries to reach the network.
type refusingTransport struct{ t *testing.T }

func (rt refusingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.t.Errorf("a dry run must not make requests, but one went to %s", r.URL)
	return nil, fmt.Errorf("refused")
}

// dryRunCmd builds a deploy command wired for a dry run: a config supplied in
// memory, a working directory of its own, and a client that refuses to make
// requests.
func dryRunCmd(t *testing.T, conf *contracts.AzionApplicationOptions) (*DeployCmd, *bytes.Buffer) {
	t.Helper()
	out := &bytes.Buffer{}
	f := &cmdutil.Factory{
		HttpClient: &http.Client{Transport: refusingTransport{t}},
		IOStreams:  &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}},
		Config:     viper.New(),
	}
	cmd := NewDeployCmd(f)
	cmd.F = f
	cmd.Io = f.IOStreams
	cmd.GetAzionJsonContent = func(string) (*contracts.AzionApplicationOptions, error) {
		return conf, nil
	}
	cmd.WriteAzionJsonContent = func(*contracts.AzionApplicationOptions, string) error {
		t.Error("a dry run must not write azion.json")
		return nil
	}
	return cmd, out
}

// inTempProject runs fn in an empty working directory, optionally containing a
// manifest.
func inTempProject(t *testing.T, manifest *contracts.ManifestV4, fn func()) {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	dir := t.TempDir()
	if manifest != nil {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, ".edge"), 0o755))
		b, err := json.Marshal(manifest)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(dir, ".edge", "manifest.json"), b, 0o644))
	}
	require.NoError(t, os.Chdir(dir))
	defer func() { require.NoError(t, os.Chdir(wd)) }()
	fn()
}

func TestDryRunMakesNoRequestsAndWritesNothing(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj"}
	cmd, _ := dryRunCmd(t, conf)

	// The refusing transport and the writing stub fail the test themselves if
	// the dry run reaches the network or touches azion.json.
	inTempProject(t, &contracts.ManifestV4{
		Functions: []contracts.Function{{}},
		Workloads: []contracts.WorkloadManifest{{}},
		Storage:   []contracts.StorageManifest{{}},
		Purge:     []contracts.PurgeManifest{{}},
	}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})
}

func TestDryRunDescribesAFirstDeploy(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj"}
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.ManifestV4{
		Storage: []contracts.StorageManifest{{}},
	}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	// A project whose azion.json has no application id is a creation.
	require.Contains(t, got, "Creating Application named 'proj'")
	// Storage declared in the manifest is what makes the bucket step apply —
	// the simulation this replaces gated it on the project preset instead.
	require.Contains(t, got, "Creating Bucket named 'proj'")
	// The manifest steps are named, which the old simulation never reached.
	require.Contains(t, got, "Applying the resources declared in your manifest")
	require.Contains(t, got, "application")
}

func TestDryRunReportsAnUpdateWhenTheProjectExists(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj"}
	conf.Application.ID = 42
	conf.NotFirstRun = true
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.ManifestV4{}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	require.Contains(t, got, "Updating Application with ID '42'")
	require.NotContains(t, got, "Creating Application")
}

func TestDryRunWithoutAManifestSaysSo(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj"}
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, nil, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	require.Contains(t, got, "has not been built yet")
	// No storage is declared, so no bucket is announced.
	require.NotContains(t, got, "Creating Bucket")
}

// TestDryRunDescribesOnlyStepsThatWouldRun is the property that matters: the
// description comes from the pipeline, so a step whose predicate says no must
// not be announced.
func TestDryRunDescribesOnlyStepsThatWouldRun(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Bucket: "proj"}
	conf.NotFirstRun = true
	conf.RulesEngine.Rules = []contracts.AzionJsonDataRules{{Name: "existing"}}
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.ManifestV4{}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	// Rules already exist, so the rules step does not apply.
	require.NotContains(t, got, "Creating the default Rules Engine")
	// No storage declared, so no bucket.
	require.NotContains(t, got, "Creating Bucket")
	// And nothing from the v3 vocabulary, which the old simulation printed even
	// on v4 projects.
	require.NotContains(t, strings.ToLower(got), "single origin")
	require.NotContains(t, got, "default Rule Engine")
}

// TestV4DryRunDoesNotOfferACacheSetting guards against describing v4 with v3
// behaviour.
//
// v3's doRulesDeploy prompts the user to create a cache setting; v4's only
// creates the preset's default rules. An earlier version of this renderer
// reused the v3 message for the v4 step, which is the same drift that made the
// simulation this replaces wrong in the first place.
func TestV4DryRunDoesNotOfferACacheSetting(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptions{Name: "proj", Preset: "vite"}
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.ManifestV4{}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	require.Contains(t, got, "Creating the default Rules Engine for the 'vite' preset")
	require.NotContains(t, got, "Cache Setting")
	require.NotContains(t, got, "Presenting the option")
}

// TestV4DryRunReflectsTheManifestContents is the v4 counterpart: its steps are
// conditional, so the list already varies, and the counts say how much.
func TestV4DryRunReflectsTheManifestContents(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	render := func(m *contracts.ManifestV4) string {
		conf := &contracts.AzionApplicationOptions{Name: "p", Bucket: "p", Preset: "vite"}
		conf.NotFirstRun = true
		cmd, out := dryRunCmd(t, conf)
		inTempProject(t, m, func() { require.NoError(t, cmd.DryRun(cmd.F)) })
		return out.String()
	}

	onlyFunctions := render(&contracts.ManifestV4{Functions: []contracts.Function{{}, {}}})
	require.Contains(t, onlyFunctions, "functions (2)")
	require.NotContains(t, onlyFunctions, "firewalls")
	require.NotContains(t, onlyFunctions, "purge")

	withFirewalls := render(&contracts.ManifestV4{
		Firewalls: []contracts.FirewallManifest{{}},
		Purge:     []contracts.PurgeManifest{{}, {}},
	})
	require.Contains(t, withFirewalls, "firewalls (1)")
	require.Contains(t, withFirewalls, "purge (2)")
	require.NotContains(t, withFirewalls, "functions")

	require.NotEqual(t, onlyFunctions, withFirewalls)
}
