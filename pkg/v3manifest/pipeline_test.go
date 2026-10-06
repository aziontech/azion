package manifest

import (
	"context"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// selected returns the steps V3Pipeline picks for a manifest, by running the
// pipeline as a dry run: every predicate is evaluated and nothing executes.
func selected(t *testing.T, m *contracts.Manifest) []string {
	t.Helper()
	rc := &ResourceContext{manifest: m}
	res, err := pipeline.Run(context.Background(), V3Pipeline{}, rc, &pipeline.Plan{DryRun: true})
	require.NoError(t, err)
	names := make([]string, 0, len(res.Steps))
	for _, s := range res.Steps {
		names = append(names, s.Name)
	}
	return names
}

// TestV3PipelineSelection pins which steps run for a given manifest.
//
// Only the domain step is conditional. The rest run for every manifest, which
// is what the sequence this replaces did: they were plain loops with no
// enclosing length check.
func TestV3PipelineSelection(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	unconditional := []string{
		"ManifestOrigins", "ManifestCacheSettings", "ManifestRulesEngine",
		"ManifestPurge", "ManifestDeleteOrphanedResources",
	}

	tests := []struct {
		name     string
		manifest *contracts.Manifest
		want     []string
	}{
		{
			name:     "empty manifest still runs every unconditional step",
			manifest: &contracts.Manifest{},
			want:     unconditional,
		},
		{
			name:     "a nil domain skips only the domain step",
			manifest: &contracts.Manifest{Origins: []contracts.Origin{{}}},
			want:     unconditional,
		},
		{
			name:     "a domain with an empty name is skipped",
			manifest: &contracts.Manifest{Domain: &contracts.Domains{Name: ""}},
			want:     unconditional,
		},
		{
			name:     "a named domain is applied first",
			manifest: &contracts.Manifest{Domain: &contracts.Domains{Name: "example"}},
			want:     append([]string{"ManifestDomain"}, unconditional...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, selected(t, tt.manifest))
		})
	}
}

// TestV3EmptyManifestStillRewritesConfig covers the reason the origins, cache
// and rules steps carry no predicate.
//
// Each one ends by assigning its result back to the project config and writing
// azion.json, outside the loop. With an empty manifest the loops do nothing but
// those assignments still run, so azion.json is rewritten with those sections
// emptied. Gating the steps on "the manifest declares some" would silently stop
// that, leaving stale entries in azion.json.
func TestV3EmptyManifestStillRewritesConfig(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	msgs := []string{}
	conf := &contracts.AzionApplicationOptionsV3{
		Name:        "project",
		Application: contracts.AzionJsonDataApplication{ID: 1},
		Origin:      []contracts.AzionJsonDataOrigin{{Name: "stale", OriginId: 9, OriginKey: "k"}},
		CacheSettings: []contracts.AzionJsonDataCacheSettings{
			{Name: "stale-cache", Id: 8},
		},
	}
	conf.RulesEngine.Rules = []contracts.AzionJsonDataRules{{Name: "stale-rule", Id: 7, Phase: "request"}}

	writes := 0
	interpreter := NewManifestInterpreter()
	interpreter.WriteAzionJsonContent = func(*contracts.AzionApplicationOptionsV3, string) error {
		writes++
		return nil
	}

	// SkipDeletion keeps orphan removal from calling the API for the stale ids.
	skip := true
	conf.SkipDeletion = &skip

	err := interpreter.CreateResources(conf, &contracts.Manifest{}, f, "azion", &msgs)
	require.NoError(t, err)

	require.Equal(t, 3, writes, "origins, cache and rules each rewrite azion.json")
	require.Empty(t, conf.Origin, "an empty manifest empties the origins in azion.json")
	require.Empty(t, conf.CacheSettings, "an empty manifest empties the cache settings in azion.json")
	require.Empty(t, conf.RulesEngine.Rules, "an empty manifest empties the rules in azion.json")
}

// TestV3FailedPurgeReportsSuccessAndSkipsCleanup pins a deliberate oddity of
// the original sequence, reproduced by errPurgeAborted.
//
// A failing purge was logged at debug level and then returned nil from
// CreateResources: the caller saw success, and orphan removal never ran. It is
// almost certainly a bug, but it is what users get today, so the port keeps it
// and this test states it out loud.
//
// The stale rule in the config is what makes the two outcomes distinguishable:
// if orphan removal did run, it would try to delete that rule, find no stub and
// fail, so a nil error here proves it was skipped.
func TestV3FailedPurgeReportsSuccessAndSkipsCleanup(t *testing.T) {
	logger.New(zapcore.InfoLevel)

	// The purge endpoint answers 500, so the request fails with a response in
	// hand. A transport-level failure would instead panic inside
	// utils.LogAndRewindBody, which dereferences a nil *http.Response — a
	// pre-existing bug on the v3 purge path, untouched by this change.
	mock := &httpmock.Registry{}
	mock.Register(
		httpmock.REST("POST", "purge/url"),
		httpmock.StatusStringResponse(500, "{}"),
	)
	f, _, _ := testutils.NewFactory(mock)
	msgs := []string{}
	conf := &contracts.AzionApplicationOptionsV3{
		Name:        "project",
		Application: contracts.AzionJsonDataApplication{ID: 1},
	}
	conf.RulesEngine.Rules = []contracts.AzionJsonDataRules{{Name: "stale-rule", Id: 7, Phase: "request"}}

	interpreter := NewManifestInterpreter()
	interpreter.WriteAzionJsonContent = func(*contracts.AzionApplicationOptionsV3, string) error { return nil }

	m := &contracts.Manifest{
		Purge: []contracts.Purges{{Type: "url", Urls: []string{"https://example.com/a"}}},
	}

	err := interpreter.CreateResources(conf, m, f, "azion", &msgs)
	require.NoError(t, err, "a failed purge is reported as success, as it was before")
}

// TestV3PipelineOrder guards the order, which encodes real dependencies: the
// id maps are seeded before anything runs, rules are built from the cache and
// origin ids the earlier steps produced, and orphan removal reads what is left.
func TestV3PipelineOrder(t *testing.T) {
	want := []string{
		"ManifestDomain", "ManifestOrigins", "ManifestCacheSettings",
		"ManifestRulesEngine", "ManifestPurge", "ManifestDeleteOrphanedResources",
	}
	got := make([]string, 0, len(want))
	for _, s := range (V3Pipeline{}).Steps() {
		got = append(got, s.Name)
	}
	require.Equal(t, want, got)
}
