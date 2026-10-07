package deploy

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// selected returns the steps V3DeployPipeline picks for a given state, by
// running the pipeline as a dry run: every predicate is evaluated and nothing
// executes. The state is pre-populated with what earlier steps would produce.
func selected(t *testing.T, s *DeployState) []string {
	t.Helper()
	res, err := pipeline.Run(context.Background(), V3DeployPipeline{}, s, &pipeline.Plan{DryRun: true})
	require.NoError(t, err)
	names := make([]string, 0, len(res.Steps))
	for _, st := range res.Steps {
		names = append(names, st.Name)
	}
	return names
}

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
		Conf:     &contracts.AzionApplicationOptionsV3{},
		Manifest: &contracts.Manifest{},
	}
	if mutate != nil {
		mutate(s)
	}
	return s
}

// TestV3DeployPipelineSelection pins each step's condition against the if
// blocks of the sequence it replaces.
func TestV3DeployPipelineSelection(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	tests := []struct {
		name   string
		static bool
		state  *DeployState
		want   []string
	}{
		{
			name:  "first run creates the origin and default rules",
			state: state(nil),
			want: []string{"ManifestPath", "Application", "FirstRunOriginAndRules",
				"Bucket", "Function", "ReadManifest", "Manifest", "Domain"},
		},
		{
			name: "repeat run skips the first-run block",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
			}),
			want: []string{"ManifestPath", "Application", "Bucket", "Function",
				"ReadManifest", "Manifest", "Domain"},
		},
		{
			name:   "static files present adds the upload step",
			static: true,
			state:  state(func(s *DeployState) { s.Conf.NotFirstRun = true }),
			want: []string{"ManifestPath", "Application", "Bucket", "UploadStaticFiles",
				"Function", "ReadManifest", "Manifest", "Domain"},
		},
		{
			// A manifest that declares a domain with an empty name is the one
			// case where deploy does not create the project's own domain.
			name: "a domain with an empty name suppresses the domain step",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
				s.Manifest.Domain = &contracts.Domains{Name: ""}
			}),
			want: []string{"ManifestPath", "Application", "Bucket", "Function",
				"ReadManifest", "Manifest"},
		},
		{
			name: "a named domain still gets the domain step",
			state: state(func(s *DeployState) {
				s.Conf.NotFirstRun = true
				s.Manifest.Domain = &contracts.Domains{Name: "example"}
			}),
			want: []string{"ManifestPath", "Application", "Bucket", "Function",
				"ReadManifest", "Manifest", "Domain"},
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

// TestV3DeployPipelineOrder guards the order, which encodes the dependencies:
// the application exists before its origin and rules, the bucket and upload
// precede the function that serves them, and the manifest is read before it is
// applied.
func TestV3DeployPipelineOrder(t *testing.T) {
	want := []string{
		"ManifestPath", "Application", "FirstRunOriginAndRules", "Bucket",
		"UploadStaticFiles", "Function", "ReadManifest", "Manifest", "Domain",
	}
	got := make([]string, 0, len(want))
	for _, s := range (V3DeployPipeline{}).Steps() {
		got = append(got, s.Name)
	}
	require.Equal(t, want, got)
}
