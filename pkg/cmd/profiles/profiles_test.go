package profiles

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

const testProfile = "default"

// newFactory isolates profile settings in a temp HOME so tests never touch the
// developer's real ~/.azion.
func newFactory(t *testing.T, mock *httpmock.Registry) (*cmdutil.Factory, *bytes.Buffer) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())

	out := &bytes.Buffer{}
	return &cmdutil.Factory{
		HttpClient: &http.Client{Transport: mock},
		IOStreams:  &iostreams.IOStreams{Out: out, Err: &bytes.Buffer{}},
		Config:     viper.New(),
	}, out
}

func seedProfile(t *testing.T, settings token.Settings) {
	t.Helper()
	require.NoError(t, token.WriteProfiles(token.Profile{Name: testProfile}))
	require.NoError(t, token.WriteSettings(settings, testProfile))
}

func stubAccountInfo(mock *httpmock.Registry, flags string) {
	mock.Register(
		httpmock.REST("GET", "account/info"),
		httpmock.JSONFromString(`{"client_flags":`+flags+`}`),
	)
}

func TestRefreshRecordsTheCurrentVersion(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `["`+apiversion.BlockAPIV4IncompatibleEndpoints+`"]`)
	f, out := newFactory(t, mock)
	seedProfile(t, token.Settings{Token: "a-token"})

	require.NoError(t, runRefresh(f))

	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	cached := settings.APIVersionCache()
	assert.Equal(t, apiversion.V3, cached.Version)
	assert.Equal(t, apiversion.TokenHash("a-token"), cached.TokenHash)
	assert.WithinDuration(t, time.Now(), cached.CheckedAt, time.Minute)
	assert.Contains(t, out.String(), "is on Azion API v3")
}

// The point of the flag: an account that migrated should not have to wait out
// the TTL.
func TestRefreshReportsAChangedGeneration(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `[]`)
	f, out := newFactory(t, mock)
	seedProfile(t, token.Settings{Token: "a-token"})

	// a still-valid cache saying v3, which the refresh must override
	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	settings.SetAPIVersionCache(apiversion.Refreshed(apiversion.V3, "a-token", time.Now()))
	require.NoError(t, token.WriteSettings(settings, testProfile))

	require.NoError(t, runRefresh(f))

	assert.Contains(t, out.String(), "moved from Azion API v3 to v4")

	settings, err = token.ReadSettings(testProfile)
	require.NoError(t, err)
	assert.Equal(t, apiversion.V4, settings.APIVersionCache().Version)
}

func TestRefreshWithoutCredential(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	f, _ := newFactory(t, mock)
	seedProfile(t, token.Settings{})

	err := runRefresh(f)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no credential configured")
	assert.Empty(t, mock.Requests, "nothing to ask without a credential")
}

// A failed refresh is an error the user asked for, not a silent downgrade, and
// it must not overwrite what is already cached.
func TestRefreshFailureKeepsTheExistingCache(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(httpmock.REST("GET", "account/info"), httpmock.StatusStringResponse(500, "boom"))
	f, _ := newFactory(t, mock)
	seedProfile(t, token.Settings{Token: "a-token"})

	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	settings.SetAPIVersionCache(apiversion.Refreshed(apiversion.V4, "a-token", time.Now()))
	require.NoError(t, token.WriteSettings(settings, testProfile))

	err = runRefresh(f)

	require.Error(t, err)
	settings, readErr := token.ReadSettings(testProfile)
	require.NoError(t, readErr)
	assert.Equal(t, apiversion.V4, settings.APIVersionCache().Version)
}

func TestRefreshFlagIsRegistered(t *testing.T) {
	f, _ := newFactory(t, &httpmock.Registry{})
	cmd := NewCmd(f)

	flag := cmd.Flags().Lookup("refresh")
	require.NotNil(t, flag, "--refresh must exist on azion profiles")
	assert.Equal(t, "false", flag.DefValue, "refresh is opt-in")
}

// The refresh must bind the cache to the credential the CLI actually uses, not
// whatever happens to be in the settings file, or root would never match the
// hash and would re-resolve on every invocation.
func TestRefreshUsesTheEffectiveCredential(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	stubAccountInfo(mock, `[]`)
	f, _ := newFactory(t, mock)
	f.Config.(*viper.Viper).Set("token", "env-token")
	seedProfile(t, token.Settings{Token: "stored-token"})

	require.NoError(t, runRefresh(f))

	settings, err := token.ReadSettings(testProfile)
	require.NoError(t, err)
	assert.Equal(t, apiversion.TokenHash("env-token"), settings.APIVersionCache().TokenHash)
}
