package deploy

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// selected returns the steps V4DeployPipeline picks for a given state.
//
// It runs the pipeline as a dry run, which evaluates every predicate and
// executes nothing. The state must therefore be pre-populated with what the
// earlier steps would have produced — notably Manifest, which the read step
// fills in a real run.
func selected(t *testing.T, s *DeployState) []string {
	t.Helper()
	res, err := pipeline.Run(context.Background(), V4DeployPipeline{}, s, &pipeline.Plan{DryRun: true})
	require.NoError(t, err)
	names := make([]string, 0, len(res.Steps))
	for _, st := range res.Steps {
		names = append(names, st.Name)
	}
	return names
}

// withStaticDir runs fn in a temp working directory, optionally containing the
// static-files directory the upload step looks for.
func withStaticDir(t *testing.T, present bool, fn func()) {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	dir := t.TempDir()
	if present {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, PathStatic), 0o755))
	}
	require.NoError(t, os.Chdir(dir))
	defer func() { require.NoError(t, os.Chdir(wd)) }()
	fn()
}

func state(mutate func(*DeployState)) *DeployState {
	s := &DeployState{
		Cmd:      &DeployCmd{},
		Conf:     &contracts.AzionApplicationOptions{},
		Manifest: &contracts.ManifestV4{},
	}
	if mutate != nil {
		mutate(s)
	}
	return s
}

// always lists the steps with no predicate, which run for every deploy.
var always = []string{"ManifestPath", "Application", "ReadManifest", "ReadManifestAgain", "Manifest"}

// TestV4DeployPipelineSelection pins each step's condition against the if
// blocks of the sequence it replaces.
func TestV4DeployPipelineSelection(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	tests := []struct {
		name   string
		static bool
		state  *DeployState
		want   []string
	}{
		{
			// First run, nothing declared: no build (NotFirstRun is false), no
			// bucket, no upload. Rules and workload both apply.
			name:  "first run with an empty manifest",
			state: state(nil),
			want: []string{"ManifestPath", "Application", "ReadManifest", "BundlerInit",
				"ReadManifestAgain", "RulesEngine", "Manifest", "Workload"},
		},
		{
			name: "repeat run builds and skips bundler init and rules",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
			}),
			want: []string{"Build", "ManifestPath", "Application", "ReadManifest",
				"ReadManifestAgain", "Manifest", "Workload"},
		},
		{
			name: "--skip-build on a repeat run skips the build",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
				s.Cmd.SkipBuild = true
			}),
			want: []string{"ManifestPath", "Application", "ReadManifest",
				"ReadManifestAgain", "Manifest", "Workload"},
		},
		{
			name: "declared storage adds the bucket step",
			state: state(func(s *DeployState) {
				s.Manifest.Storage = []contracts.StorageManifest{{}}
			}),
			want: []string{"ManifestPath", "Application", "ReadManifest", "Bucket",
				"BundlerInit", "ReadManifestAgain", "RulesEngine", "Manifest", "Workload"},
		},
		{
			name:   "static files present adds the upload step",
			static: true,
			state:  state(nil),
			want: []string{"ManifestPath", "Application", "ReadManifest", "BundlerInit",
				"ReadManifestAgain", "Credentials", "UploadStaticFiles", "RulesEngine",
				"Manifest", "Workload"},
		},
		{
			name:   "--skip-build suppresses the upload even when static files exist",
			static: true,
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
				s.Cmd.SkipBuild = true
			}),
			want: []string{"ManifestPath", "Application", "ReadManifest",
				"ReadManifestAgain", "Manifest", "Workload"},
		},
		{
			name: "existing rules suppress the rules step",
			state: state(func(s *DeployState) {
				s.Conf.RulesEngine.Rules = []contracts.AzionJsonDataRules{{Name: "r"}}
			}),
			want: []string{"ManifestPath", "Application", "ReadManifest", "BundlerInit",
				"ReadManifestAgain", "Manifest", "Workload"},
		},
		{
			name: "a named workload in the manifest suppresses the workload step",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
				s.Manifest.Workloads = []contracts.WorkloadManifest{{Name: "w"}}
			}),
			want: []string{"Build", "ManifestPath", "Application", "ReadManifest",
				"ReadManifestAgain", "Manifest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withStaticDir(t, tt.static, func() {
				require.Equal(t, tt.want, selected(t, tt.state))
			})
		})
	}
}

// TestV4DeployPipelineAlwaysRunsTheCoreSteps guards the steps that carry no
// predicate: the manifest path, the application, both reads and the manifest
// apply happen on every deploy.
func TestV4DeployPipelineAlwaysRunsTheCoreSteps(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	withStaticDir(t, false, func() {
		got := selected(t, state(func(s *DeployState) { s.Conf.NotFirstRun = true }))
		for _, name := range always {
			require.Contains(t, got, name)
		}
	})
}

// TestV4DeployPipelineOrder guards the order, which encodes the data
// dependencies: the path is needed by both reads, the second read follows the
// build that may regenerate the manifest, and the manifest apply follows the
// upload that its rules point at.
func TestV4DeployPipelineOrder(t *testing.T) {
	want := []string{
		"Build", "ManifestPath", "Application", "ReadManifest", "Bucket",
		"BundlerInit", "ReadManifestAgain", "Credentials", "UploadStaticFiles",
		"RulesEngine", "Manifest", "Workload",
	}
	got := make([]string, 0, len(want))
	for _, s := range (V4DeployPipeline{}).Steps() {
		got = append(got, s.Name)
	}
	require.Equal(t, want, got)
}

// TestDeployTimingCallbackFillsTheSummary checks the mapping from step name to
// summary field, including the two reads that accumulate onto one field.
//
// The sequence this replaces assigned ReadManifestTime on the first read and
// added to it on the second. Both steps now add, which is equivalent because
// InitTimingSummary starts the field at zero.
func TestDeployTimingCallbackFillsTheSummary(t *testing.T) {
	InitTimingSummary()

	HandleDeployTimingCallback("ReadManifest", 2*time.Second)
	HandleDeployTimingCallback("ReadManifestAgain", 3*time.Second)
	HandleDeployTimingCallback("Bucket", 5*time.Second)
	HandleDeployTimingCallback("Credentials", 7*time.Second)
	HandleDeployTimingCallback("UploadStaticFiles", 11*time.Second)
	HandleDeployTimingCallback("Manifest", 13*time.Second)
	// A step with no field must be ignored rather than panic or misfile.
	HandleDeployTimingCallback("Application", 17*time.Second)

	require.Equal(t, 5*time.Second, GlobalTimingSummary.ReadManifestTime,
		"both reads accumulate onto one field")
	require.Equal(t, 5*time.Second, GlobalTimingSummary.BucketCreateTime)
	require.Equal(t, 7*time.Second, GlobalTimingSummary.CredentialsTime)
	require.Equal(t, 11*time.Second, GlobalTimingSummary.UploadStaticFilesTime)
	require.Equal(t, 13*time.Second, GlobalTimingSummary.ManifestCreateTime)
}

// TestEveryTimedStepNameExists guards the mapping against a step rename: each
// name the callback knows must still be produced by the pipeline.
func TestEveryTimedStepNameExists(t *testing.T) {
	timed := map[string]bool{
		"ReadManifest": true, "ReadManifestAgain": true, "Bucket": true,
		"Credentials": true, "UploadStaticFiles": true, "Manifest": true,
	}
	for _, s := range (V4DeployPipeline{}).Steps() {
		delete(timed, s.Name)
	}
	for name := range timed {
		t.Errorf("the timing summary records %q but no pipeline step emits it", name)
	}
}

// TestSkippedUploadLogsItsReasonOnce guards the memoisation in shouldUpload.
//
// Two steps share the predicate, so without memoisation a skipped upload would
// log its reason twice where the if/else-if chain it replaces logged it once.
func TestSkippedUploadLogsItsReasonOnce(t *testing.T) {
	logger.New(zapcore.InfoLevel)
	withStaticDir(t, false, func() {
		s := state(nil)
		evaluations := 0
		for _, step := range (V4DeployPipeline{}).Steps() {
			if step.Name == "Credentials" || step.Name == "UploadStaticFiles" {
				require.False(t, step.Applies(s))
				evaluations++
			}
		}
		require.Equal(t, 2, evaluations, "both steps ask the predicate")
		require.True(t, s.uploadDecided, "the answer is computed once and remembered")
	})
}
