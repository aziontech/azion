package root

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

const testProfile = "default"

// newRootFactory builds a factoryRoot whose profile settings live in a temp
// HOME, so tests never touch the developer's real ~/.azion.
func newRootFactory(t *testing.T, mock *httpmock.Registry, tok string) (*factoryRoot, *bytes.Buffer) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())

	out := &bytes.Buffer{}
	cfg := viper.New()
	cfg.Set("token", tok)

	f := &cmdutil.Factory{
		HttpClient: &http.Client{Transport: mock},
		IOStreams:  &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}},
		Config:     cfg,
	}
	return &factoryRoot{factory: f}, out
}

func stubAccountInfo(mock *httpmock.Registry, flags string) {
	mock.Register(
		httpmock.REST("GET", "account/info"),
		httpmock.JSONFromString(`{"client_flags":`+flags+`}`),
	)
}

func writeCache(t *testing.T, c apiversion.Cache) {
	t.Helper()
	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	settings.SetAPIVersionCache(c)
	require.NoError(t, token.WriteSettings(settings, testProfile))
}

func readCache(t *testing.T) apiversion.Cache {
	t.Helper()
	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	return settings.APIVersionCache()
}

func TestResolveAPIVersionQueriesAndCaches(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `["`+apiversion.BlockAPIV4IncompatibleEndpoints+`"]`)
	fact, _ := newRootFactory(t, mock, "a-token")

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V3, got, "a blocked account resolves to the legacy generation")
	assert.Len(t, mock.Requests, 1, "a cold cache asks the SSO service once")

	cached := readCache(t)
	assert.Equal(t, apiversion.V3, cached.Version)
	assert.Equal(t, apiversion.TokenHash("a-token"), cached.TokenHash)
	assert.WithinDuration(t, time.Now(), cached.CheckedAt, time.Minute)
}

// D2: a warm cache must make no network call at all.
func TestResolveAPIVersionUsesCacheWithoutNetwork(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	fact, _ := newRootFactory(t, mock, "a-token")
	writeCache(t, apiversion.Refreshed(apiversion.V4, "a-token", time.Now()))

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V4, got)
	assert.Empty(t, mock.Requests, "a warm cache must not call the SSO service")
}

func TestResolveAPIVersionRechecksWhenCacheIsStale(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `[]`)
	fact, _ := newRootFactory(t, mock, "a-token")
	writeCache(t, apiversion.Refreshed(apiversion.V3, "a-token", time.Now().Add(-apiversion.TTL-time.Minute)))

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V4, got, "an expired entry is re-resolved")
	assert.Len(t, mock.Requests, 1)
	assert.Equal(t, apiversion.V4, readCache(t).Version)
}

// The cache is bound to the credential that produced it: an account that
// migrated under a new token must not be served the previous answer.
func TestResolveAPIVersionRechecksWhenCredentialChanges(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `[]`)
	fact, _ := newRootFactory(t, mock, "new-token")
	writeCache(t, apiversion.Refreshed(apiversion.V3, "old-token", time.Now()))

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V4, got)
	assert.Len(t, mock.Requests, 1, "a different credential re-resolves even within the TTL")
	assert.Equal(t, apiversion.TokenHash("new-token"), readCache(t).TokenHash)
}

// D3, narrowed: with something cached, a failed lookup reuses it rather than
// silently downgrading the account.
func TestResolveAPIVersionReusesStaleCacheWhenLookupFails(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
	fact, out := newRootFactory(t, mock, "a-token")
	writeCache(t, apiversion.Refreshed(apiversion.V4, "a-token", time.Now().Add(-apiversion.TTL-time.Minute)))

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V4, got, "an unreachable lookup must not downgrade a known v4 account")
	assert.NotContains(t, out.String(), "Warning")
}

// D3, narrowed: the fallback now applies only with nothing cached, and it says
// so on stdout instead of only in a debug log.
func TestResolveAPIVersionWarnsWhenFallingBack(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
	fact, out := newRootFactory(t, mock, "a-token")

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V3, got)
	assert.Contains(t, out.String(), "could not verify which Azion API version")
	assert.Contains(t, out.String(), "azion profiles --refresh")
	assert.Empty(t, readCache(t).Version, "a failed lookup must not be cached")
}

func TestResolveAPIVersionWithoutToken(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	fact, out := newRootFactory(t, mock, "")

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V4, got, "with no credential there is no block flag, so the current generation applies")
	assert.Empty(t, mock.Requests, "there is nothing to ask about without a credential")
	assert.Contains(t, out.String(), "Please remember to login")
	assert.Empty(t, readCache(t).Version, "a logged-out CLI caches nothing")
}

// CmdRoot records the generation on the factory and uses it to pick the tree,
// which is what setCmds/setV3Cmds consume.
func TestCmdRootResolvesVersionAndPicksTree(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name        string
		token       string
		clientFlags string
		wantVersion apiversion.Version
		wantWarmup  bool
	}{
		{"no token builds the v4 tree", "", "", apiversion.V4, true},
		{"account with no flags builds the v4 tree", "a-token", `[]`, apiversion.V4, true},
		{"account blocked from v4 builds the v3 tree", "a-token", `["` + apiversion.BlockAPIV4IncompatibleEndpoints + `"]`, apiversion.V3, false},
		{"account blocked from v3 builds the v4 tree", "a-token", `["` + apiversion.BlockAPIV3Access + `"]`, apiversion.V4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			if tt.token != "" {
				stubAccountInfo(mock, tt.clientFlags)
			}
			fact, _ := newRootFactory(t, mock, tt.token)

			cmd := fact.CmdRoot().(*cobra.Command)

			assert.Equal(t, tt.wantVersion, fact.factory.APIVersion)

			names := map[string]bool{}
			for _, c := range cmd.Commands() {
				names[c.Name()] = true
			}
			assert.Equal(t, tt.wantWarmup, names["warmup"], "warmup exists only in the v4 tree")
		})
	}
}

// The debug line is emitted from persistentPreRunE, because resolution runs
// before Cobra has parsed --debug. These pin what it will report.
func TestAPIVersionSourceIsRecorded(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("a warm cache records the cache as the source", func(t *testing.T) {
		mock := &httpmock.Registry{}
		fact, _ := newRootFactory(t, mock, "a-token")
		checkedAt := time.Now().Add(-2 * time.Hour)
		writeCache(t, apiversion.Refreshed(apiversion.V4, "a-token", checkedAt))

		fact.resolveAPIVersion()

		src := fact.apiVersionSource
		assert.Equal(t, sourceCache, src.source)
		assert.Equal(t, apiversion.V4, src.version)
		assert.Equal(t, apiversion.CacheHit, src.status)
		assert.InDelta(t, (2 * time.Hour).Seconds(), src.age.Seconds(), 60)
		assert.WithinDuration(t, checkedAt, src.checkedAt, time.Second)
	})

	t.Run("an expired entry records the lookup and why it happened", func(t *testing.T) {
		mock := &httpmock.Registry{}
		stubAccountInfo(mock, `[]`)
		fact, _ := newRootFactory(t, mock, "a-token")
		writeCache(t, apiversion.Refreshed(apiversion.V3, "a-token", time.Now().Add(-apiversion.TTL-time.Minute)))

		fact.resolveAPIVersion()

		src := fact.apiVersionSource
		assert.Equal(t, sourceLookup, src.source)
		assert.Equal(t, apiversion.V4, src.version)
		assert.Equal(t, apiversion.V3, src.previous)
		assert.Equal(t, apiversion.CacheMissExpired, src.status, "the reason names the 24h TTL, not a generic miss")
	})

	t.Run("a changed credential is recorded as the reason", func(t *testing.T) {
		mock := &httpmock.Registry{}
		stubAccountInfo(mock, `["`+apiversion.BlockAPIV4IncompatibleEndpoints+`"]`)
		fact, _ := newRootFactory(t, mock, "new-token")
		writeCache(t, apiversion.Refreshed(apiversion.V4, "old-token", time.Now()))

		fact.resolveAPIVersion()

		assert.Equal(t, sourceLookup, fact.apiVersionSource.source)
		assert.Equal(t, apiversion.CacheMissCredential, fact.apiVersionSource.status)
	})

	t.Run("a failed lookup over a stale entry is recorded with the error", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
		fact, _ := newRootFactory(t, mock, "a-token")
		writeCache(t, apiversion.Refreshed(apiversion.V4, "a-token", time.Now().Add(-apiversion.TTL-time.Minute)))

		fact.resolveAPIVersion()

		src := fact.apiVersionSource
		assert.Equal(t, sourceStaleCache, src.source)
		assert.Equal(t, apiversion.V4, src.version)
		assert.Error(t, src.err)
	})

	t.Run("the fallback is recorded as such", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
		fact, _ := newRootFactory(t, mock, "a-token")

		fact.resolveAPIVersion()

		assert.Equal(t, sourceFallback, fact.apiVersionSource.source)
		assert.Equal(t, apiversion.V3, fact.apiVersionSource.version)
		assert.Error(t, fact.apiVersionSource.err)
	})

	// Nothing recorded means nothing logged, rather than a misleading line.
	t.Run("logging an empty source is a no-op", func(t *testing.T) {
		fact, _ := newRootFactory(t, &httpmock.Registry{}, "")
		assert.NotPanics(t, fact.logAPIVersionSource)
	})
}

func TestPreParseGlobalFlags(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name       string
		args       []string
		wantToken  string
		wantConfig string
	}{
		{"long form with a space", []string{"--token", "tok123", "list", "workload"}, "tok123", ""},
		{"long form with equals", []string{"--token=tok123", "list"}, "tok123", ""},
		{"shorthand with a space", []string{"-t", "tok123"}, "tok123", ""},
		{"shorthand glued", []string{"-ttok123"}, "tok123", ""},
		{"after the subcommand", []string{"list", "workload", "--token", "tok123"}, "tok123", ""},
		{"config long form", []string{"--config", "/tmp/x.toml"}, "", "/tmp/x.toml"},
		{"config shorthand", []string{"-c", "/tmp/x.toml"}, "", "/tmp/x.toml"},
		{"both together", []string{"--config", "/tmp/x.toml", "--token", "tok123"}, "tok123", "/tmp/x.toml"},
		{"alongside other global flags", []string{"--debug", "--format", "json", "-t", "tok123"}, "tok123", ""},
		{"alongside unknown subcommand flags", []string{"deploy", "--path", "dist", "--token", "tok123"}, "tok123", ""},
		{"absent", []string{"list", "workload"}, "", ""},
		{"no arguments at all", []string{}, "", ""},
		{"after a terminator is not a flag", []string{"--", "--token", "tok123"}, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fact, _ := newRootFactory(t, &httpmock.Registry{}, "")

			fact.preParseGlobalFlags(tt.args)

			assert.Equal(t, tt.wantToken, fact.tokenFlag)
			assert.Equal(t, tt.wantConfig, fact.configFlag)
		})
	}
}

// Malformed input must not panic or abort here; Cobra reports it moments later.
func TestPreParseGlobalFlagsToleratesBadInput(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	for _, args := range [][]string{
		{"--token"},              // value missing
		{"--config"},             // value missing
		{"--", "--"},             // nothing after the terminator
		{"-"},                    // a bare dash
		{"--unknown-flag", "vv"}, // a flag this build does not know
	} {
		fact, _ := newRootFactory(t, &httpmock.Registry{}, "")
		assert.NotPanics(t, func() { fact.preParseGlobalFlags(args) }, "args=%v", args)
	}
}

func TestEffectiveTokenPrefersTheFlag(t *testing.T) {
	fact, _ := newRootFactory(t, &httpmock.Registry{}, "stored-token")
	assert.Equal(t, "stored-token", fact.effectiveToken(), "without the flag, the stored credential is used")

	fact.tokenFlag = "flag-token"
	assert.Equal(t, "flag-token", fact.effectiveToken(), "--token wins")
}

func TestCmdRootPicksTreeFromTheTokenFlag(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(
		func(req *http.Request) bool {
			return req.Header.Get("Authorization") == "Token flag-token" &&
				httpmock.REST("GET", "account/info")(req)
		},
		httpmock.JSONFromString(`{"client_flags":[]}`),
	)
	// the stored credential belongs to a legacy account
	fact, _ := newRootFactory(t, mock, "stored-token")
	writeCache(t, apiversion.Refreshed(apiversion.V3, "stored-token", time.Now()))

	// CmdRoot pre-parses os.Args itself, after setFlags has reset the bindings.
	defer func(saved []string) { os.Args = saved }(os.Args)
	os.Args = []string{"azion", "--token", "flag-token", "list", "workload"}

	cmd := fact.CmdRoot().(*cobra.Command)

	assert.Equal(t, apiversion.V4, fact.factory.APIVersion,
		"the account behind --token decides the generation")

	names := map[string]bool{}
	for _, c := range cmd.Commands() {
		names[c.Name()] = true
	}
	assert.True(t, names["warmup"], "the v4 tree must be built for a v4 --token")
}

// A rejected credential is reported by the command being run, not by the
// version resolver. Warning here would repeat on every invocation and point at
// `azion profiles --refresh`, which cannot work while the credential is bad.
func TestResolveAPIVersionStaysQuietWhenCredentialIsRejected(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	for _, status := range []int{401, 403} {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(status, `{"detail":"Invalid token."}`))
		fact, out := newRootFactory(t, mock, "expired-token")

		got := fact.resolveAPIVersion()

		assert.Equal(t, apiversion.V3, got)
		assert.NotContains(t, out.String(), "could not verify which Azion API version",
			"status %d must not print the lookup warning", status)
		assert.Equal(t, sourceFallback, fact.apiVersionSource.source)
	}
}

// The warning must still fire when the lookup itself could not be completed,
// which is the case D3 narrowed it to.
func TestResolveAPIVersionStillWarnsOnLookupFailure(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
	fact, out := newRootFactory(t, mock, "a-token")

	got := fact.resolveAPIVersion()

	assert.Equal(t, apiversion.V3, got)
	assert.Contains(t, out.String(), "could not verify which Azion API version")
}

// Resolving the API version authenticates with the same credential against the
// same service, so its outcome answers "is this token still valid?" without a
// second round trip. These pin when that answer may be reused and when it may
// not.
func TestCredentialStatusFromVersionLookup(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	const tok = "a-token"

	t.Run("a fresh lookup proves the credential", func(t *testing.T) {
		mock := &httpmock.Registry{}
		stubAccountInfo(mock, `[]`)
		fact, _ := newRootFactory(t, mock, tok)

		fact.resolveAPIVersion()

		assert.Equal(t, credentialAccepted, fact.credentialStatusFor(tok))
	})

	t.Run("a cache hit proves it too, with no network at all", func(t *testing.T) {
		mock := &httpmock.Registry{}
		fact, _ := newRootFactory(t, mock, tok)
		writeCache(t, apiversion.Refreshed(apiversion.V4, tok, time.Now()))

		fact.resolveAPIVersion()

		require.Empty(t, mock.Requests, "precondition: the cache was used")
		assert.Equal(t, credentialAccepted, fact.credentialStatusFor(tok))
	})

	t.Run("a refused credential is reported as rejected", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(401, `{"detail":"Invalid token."}`))
		fact, _ := newRootFactory(t, mock, tok)

		fact.resolveAPIVersion()

		assert.Equal(t, credentialRejected, fact.credentialStatusFor(tok))
	})

	// An unreachable service establishes nothing, so the caller must not treat
	// a reused stale version as proof the token is still good.
	t.Run("an unreachable service leaves it unknown", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
		fact, _ := newRootFactory(t, mock, tok)
		writeCache(t, apiversion.Refreshed(apiversion.V4, tok, time.Now().Add(-apiversion.TTL-time.Minute)))

		fact.resolveAPIVersion()

		require.Equal(t, sourceStaleCache, fact.apiVersionSource.source)
		assert.Equal(t, credentialUnknown, fact.credentialStatusFor(tok))
	})

	// The answer belongs to the credential the lookup used; a different one has
	// to be checked on its own.
	t.Run("a different credential is unknown", func(t *testing.T) {
		mock := &httpmock.Registry{}
		stubAccountInfo(mock, `[]`)
		fact, _ := newRootFactory(t, mock, tok)

		fact.resolveAPIVersion()

		assert.Equal(t, credentialUnknown, fact.credentialStatusFor("some-other-token"))
		assert.Equal(t, credentialUnknown, fact.credentialStatusFor(""))
	})

	t.Run("nothing resolved means nothing is known", func(t *testing.T) {
		fact, _ := newRootFactory(t, &httpmock.Registry{}, "")
		assert.Equal(t, credentialUnknown, fact.credentialStatusFor(tok))
	})
}
