package manifest

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/pipeline"
)

// V4Pipeline is the deploy pipeline for the v4 API.
//
// It replaces the sequence of if blocks that CreateResources used to spell out,
// where the same six lines — length check, time.Now(), call, error check,
// timing callback — were repeated for every resource. What is left here is the
// part that carries meaning: which steps exist, in which order, and what each
// one needs in the manifest before it is worth running.
//
// The pipeline is stateless. ResourceContext is the state: it already holds the
// manifest, the API clients and the id maps that the Apply methods read and
// write, so it is handed to pipeline.Run and reaches every step from there.
type V4Pipeline struct{}

func (V4Pipeline) Name() string { return "v4" }

// hasApplication reports whether the manifest declares an application.
//
// Cache settings, connectors and rules engine are gated on this in addition to
// their own condition, because in the sequence this replaces they lived inside
// the "if len(manifest.Applications) > 0" block. That nesting is easy to miss
// when reading the original top to bottom, and it is load-bearing: a manifest
// declaring connectors but no application applies no connectors today, and must
// keep behaving that way.
func hasApplication(rc *ResourceContext) bool {
	return len(rc.Manifest.Applications) > 0
}

// application returns the application the manifest declares. Only valid where
// hasApplication holds, which is why every caller is gated on it.
func application(rc *ResourceContext) contracts.Applications {
	return rc.Manifest.Applications[0]
}

// withApplication gates a condition on the manifest declaring an application.
func withApplication(cond func(*ResourceContext) bool) func(*ResourceContext) bool {
	return func(rc *ResourceContext) bool { return hasApplication(rc) && cond(rc) }
}

// Steps returns the v4 deploy steps in dependency order.
//
// The names are a contract with pkg/cmd/deploy_remote/timer.go, which switches
// on them to fill its timing summary and silently ignores anything it does not
// recognise. TestStepNamesAreRecordedByTheTimingSummary there guards the pair.
func (V4Pipeline) Steps() []pipeline.Step[ResourceContext] {
	return []pipeline.Step[ResourceContext]{
		{
			Name:    "ManifestFunctions",
			Applies: func(rc *ResourceContext) bool { return len(rc.Manifest.Functions) > 0 },
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyFunctions(rc.Manifest.Functions)
			},
		},
		{
			Name: "ManifestFunctionInstances",
			Applies: withApplication(func(rc *ResourceContext) bool {
				return len(application(rc).FunctionsInstances) > 0
			}),
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyFunctionInstances(application(rc).FunctionsInstances)
			},
		},
		{
			Name:    "ManifestEdgeApplication",
			Applies: hasApplication,
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyEdgeApplication(application(rc))
			},
		},
		{
			Name: "ManifestCacheSettings",
			Applies: withApplication(func(rc *ResourceContext) bool {
				return len(application(rc).CacheSettings) > 0
			}),
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyCacheSettings(application(rc).CacheSettings)
			},
		},
		{
			// Connectors are read from the manifest root, not from the
			// application, but are gated on the application existing — that is
			// how the original nesting behaves.
			Name: "ManifestConnectors",
			Applies: withApplication(func(rc *ResourceContext) bool {
				return len(rc.Manifest.Connectors) > 0
			}),
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyConnectors(rc.Manifest.Connectors)
			},
		},
		{
			Name: "ManifestRulesEngine",
			Applies: withApplication(func(rc *ResourceContext) bool {
				return len(application(rc).Rules) > 0
			}),
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyRulesEngine(application(rc).Rules)
			},
		},
		{
			Name:    "ManifestWorkloads",
			Applies: func(rc *ResourceContext) bool { return len(rc.Manifest.Workloads) > 0 },
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyWorkloads(rc.Manifest.Workloads)
			},
		},
		{
			Name:    "ManifestWorkloadDeployments",
			Applies: func(rc *ResourceContext) bool { return len(rc.Manifest.WorkloadDeployments) > 0 },
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyWorkloadDeployments(rc.Manifest.WorkloadDeployments)
			},
		},
		{
			Name:    "ManifestFirewalls",
			Applies: func(rc *ResourceContext) bool { return len(rc.Manifest.Firewalls) > 0 },
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyFirewalls(rc.Manifest.Firewalls)
			},
		},
		{
			Name:    "ManifestPurge",
			Applies: func(rc *ResourceContext) bool { return len(rc.Manifest.Purge) > 0 },
			Run: func(_ context.Context, rc *ResourceContext) error {
				return rc.ApplyPurge(rc.Manifest.Purge)
			},
		},
		{
			// Orphan removal is the trailing step rather than a deferred hook,
			// so that a failure earlier in the run skips it — which is what the
			// original sequence did by returning early.
			//
			// It has no entry in the timing summary, by inheritance: the
			// original never timed it.
			Name: "ManifestDeleteOrphanedResources",
			Run: func(_ context.Context, rc *ResourceContext) error {
				// deleteResources reads these package-level maps, and reads them
				// minus whatever this run consumed: a cache setting still
				// referenced by a rule is not an orphan.
				CacheIds = rc.CacheIds
				RuleIds = rc.RuleIds
				return rc.DeleteOrphanedResources()
			},
		},
	}
}
