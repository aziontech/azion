package manifest

import (
	"context"
	"testing"

	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/pipeline"
)

// selected returns the steps V4Pipeline picks for a manifest.
//
// It runs the pipeline as a dry run, which evaluates every predicate and
// executes nothing, so the selection can be tested without API clients.
func selected(t *testing.T, m *contracts.ManifestV4) []string {
	t.Helper()
	rc := &ResourceContext{Manifest: m}
	res, err := pipeline.Run(context.Background(), V4Pipeline{}, rc, &pipeline.Plan{DryRun: true})
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
	names := make([]string, 0, len(res.Steps))
	for _, s := range res.Steps {
		names = append(names, s.Name)
	}
	return names
}

// cleanup is the trailing step, which has no predicate and so is selected for
// every manifest.
const cleanup = "ManifestDeleteOrphanedResources"

// TestV4PipelineSelection pins which steps run for a given manifest, matching
// the sequence of if blocks that CreateResources used to spell out.
func TestV4PipelineSelection(t *testing.T) {
	app := func(mutate func(*contracts.Applications)) []contracts.Applications {
		a := contracts.Applications{}
		if mutate != nil {
			mutate(&a)
		}
		return []contracts.Applications{a}
	}

	tests := []struct {
		name     string
		manifest *contracts.ManifestV4
		want     []string
	}{
		{
			name:     "empty manifest runs only the cleanup",
			manifest: &contracts.ManifestV4{},
			want:     []string{cleanup},
		},
		{
			name: "everything declared runs every step in order",
			manifest: &contracts.ManifestV4{
				Functions: []contracts.Function{{}},
				Applications: app(func(a *contracts.Applications) {
					a.FunctionsInstances = []contracts.FunctionInstance{{}}
					a.CacheSettings = []contracts.ManifestCacheSetting{{}}
					a.Rules = []contracts.ManifestRulesEngine{{}}
				}),
				Connectors:          []edgesdk.ConnectorRequest{{}},
				Workloads:           []contracts.WorkloadManifest{{}},
				WorkloadDeployments: []contracts.WorkloadDeployment{{}},
				Firewalls:           []contracts.FirewallManifest{{}},
				Purge:               []contracts.PurgeManifest{{}},
			},
			want: []string{
				"ManifestFunctions", "ManifestFunctionInstances", "ManifestEdgeApplication",
				"ManifestCacheSettings", "ManifestConnectors", "ManifestRulesEngine",
				"ManifestWorkloads", "ManifestWorkloadDeployments", "ManifestFirewalls",
				"ManifestPurge", cleanup,
			},
		},
		{
			// The nesting that is easy to miss: connectors live at the manifest
			// root but their block sat inside the application check, so a
			// manifest with connectors and no application applies none.
			name: "connectors without an application are not applied",
			manifest: &contracts.ManifestV4{
				Connectors: []edgesdk.ConnectorRequest{{}},
			},
			want: []string{cleanup},
		},
		{
			name: "an application with no children still applies the application",
			manifest: &contracts.ManifestV4{
				Applications: app(nil),
			},
			want: []string{"ManifestEdgeApplication", cleanup},
		},
		{
			name: "function instances need the application",
			manifest: &contracts.ManifestV4{
				Functions: []contracts.Function{{}},
				Applications: app(func(a *contracts.Applications) {
					a.FunctionsInstances = []contracts.FunctionInstance{{}}
				}),
			},
			want: []string{
				"ManifestFunctions", "ManifestFunctionInstances",
				"ManifestEdgeApplication", cleanup,
			},
		},
		{
			// Functions, workloads, deployments, firewalls and purge are not
			// nested: they run with no application declared.
			name: "root-level resources need no application",
			manifest: &contracts.ManifestV4{
				Functions:           []contracts.Function{{}},
				Workloads:           []contracts.WorkloadManifest{{}},
				WorkloadDeployments: []contracts.WorkloadDeployment{{}},
				Firewalls:           []contracts.FirewallManifest{{}},
				Purge:               []contracts.PurgeManifest{{}},
			},
			want: []string{
				"ManifestFunctions", "ManifestWorkloads", "ManifestWorkloadDeployments",
				"ManifestFirewalls", "ManifestPurge", cleanup,
			},
		},
		{
			name: "cache settings and rules are nested under the application too",
			manifest: &contracts.ManifestV4{
				Applications: app(func(a *contracts.Applications) {
					a.CacheSettings = []contracts.ManifestCacheSetting{{}}
					a.Rules = []contracts.ManifestRulesEngine{{}}
				}),
			},
			want: []string{
				"ManifestEdgeApplication", "ManifestCacheSettings",
				"ManifestRulesEngine", cleanup,
			},
		},
		{
			// Storage is in the manifest but has never been part of the deploy
			// pipeline — ApplyStorage is called by azion config apply instead.
			name: "storage is not a deploy step",
			manifest: &contracts.ManifestV4{
				Storage: []contracts.StorageManifest{{}},
			},
			want: []string{cleanup},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selected(t, tt.manifest)
			if len(got) != len(tt.want) {
				t.Fatalf("selected %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("selected %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestV4PipelineCleanupIsLast guards the one ordering constraint that is not
// about resource dependencies: orphan removal reads the ids the earlier steps
// populated, so it cannot move.
func TestV4PipelineCleanupIsLast(t *testing.T) {
	steps := V4Pipeline{}.Steps()
	if got := steps[len(steps)-1].Name; got != cleanup {
		t.Errorf("last step is %q, want %q", got, cleanup)
	}
	if steps[len(steps)-1].Applies != nil {
		t.Error("the cleanup step must have no predicate, so it always runs")
	}
}
