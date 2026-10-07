package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func dryRunCmd(t *testing.T, conf *contracts.AzionApplicationOptionsV3) (*DeployCmd, *bytes.Buffer) {
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
	cmd.GetAzionJsonContent = func(string) (*contracts.AzionApplicationOptionsV3, error) {
		return conf, nil
	}
	cmd.WriteAzionJsonContent = func(*contracts.AzionApplicationOptionsV3, string) error {
		t.Error("a dry run must not write azion.json")
		return nil
	}
	return cmd, out
}

func inTempProject(t *testing.T, manifest *contracts.Manifest, fn func()) {
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

func TestV3DryRunMakesNoRequestsAndWritesNothing(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptionsV3{Name: "proj", Bucket: "proj"}
	cmd, _ := dryRunCmd(t, conf)
	inTempProject(t, &contracts.Manifest{
		Origins: []contracts.Origin{{}},
		Rules:   []contracts.RuleEngine{{}},
	}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})
}

func TestV3DryRunDescribesAFirstDeploy(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptionsV3{Name: "proj", Bucket: "proj"}
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.Manifest{Origins: []contracts.Origin{{}}}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	require.Contains(t, got, "Creating Application named 'proj'")
	// The first-run step announces all three things it does.
	require.Contains(t, got, "Creating single Origin named 'proj_single'")
	require.Contains(t, got, "Updating default Rule Engine")
	require.Contains(t, got, "Creating Function named 'proj'")
	require.Contains(t, got, "Creating Domain named 'proj'")
	require.Contains(t, got, "Applying the resources declared in your manifest")
}

func TestV3DryRunOnARepeatDeploySkipsTheFirstRunStep(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	conf := &contracts.AzionApplicationOptionsV3{Name: "proj", Bucket: "proj"}
	conf.Application.ID = 7
	conf.NotFirstRun = true
	conf.Domain.Id = 9
	cmd, out := dryRunCmd(t, conf)

	inTempProject(t, &contracts.Manifest{}, func() {
		require.NoError(t, cmd.DryRun(cmd.F))
	})

	got := out.String()
	require.Contains(t, got, "Updating Application with ID '7'")
	require.Contains(t, got, "Updating Domain with ID '9'")
	require.NotContains(t, got, "Creating single Origin")
	require.NotContains(t, got, "Updating default Rule Engine")
}
