package root

import (
	"errors"
	"io"
	"strings"
	"time"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// resolveAPIVersion decides which command tree to build for this invocation.
//
// The resolved generation is cached in the active profile for apiversion.TTL,
// so a normal invocation makes no network call at all. The cache is bound to
// the credential that produced it, so a new login, a profile switch or a
// different token re-resolves instead of trusting a stale answer.
func (fact *factoryRoot) resolveAPIVersion() apiversion.Version {
	tok := fact.effectiveToken()

	// With no credential there is nothing to ask about; this is the logged-out
	// CLI, which has always been served the default generation.
	if tok == "" {
		logger.Debug("Token is not configured")
		logger.FInfoFlags(fact.factory.IOStreams.Out, msg.LoginMessage, fact.factory.Format, fact.factory.Out)
		return apiversion.Select(apiversion.AccountInfo{})
	}

	activeProfile := fact.factory.GetActiveProfile()
	settings, err := token.ReadSettings(activeProfile)
	if err != nil {
		logger.Debug("Could not read profile settings for the API version cache", zap.Error(err))
	}
	cache := settings.APIVersionCache()

	now := time.Now()
	version, status := cache.Lookup(tok, now)
	if status.Hit() {
		fact.apiVersionSource = apiVersionSource{
			source: sourceCache, version: version, status: status,
			age: cache.Age(now), checkedAt: cache.CheckedAt, tokenHash: cache.TokenHash,
		}
		return version
	}

	version, err = apiversion.Resolve(fact.factory.HttpClient, constants.AuthURL, tok)
	if err != nil {
		// A cached answer, even an expired one, beats silently downgrading an
		// account that the lookup simply could not reach right now.
		if cached, ok := cache.Stale(); ok {
			fact.apiVersionSource = apiVersionSource{
				source: sourceStaleCache, version: cached, status: status,
				age: cache.Age(now), checkedAt: cache.CheckedAt, tokenHash: cache.TokenHash, err: err,
			}
			return cached
		}

		fact.apiVersionSource = apiVersionSource{
			source: sourceFallback, version: apiversion.V3, status: status,
			tokenHash: apiversion.TokenHash(tok), err: err,
		}

		// A rejected credential is not a failed lookup: the account's generation
		// is simply unknowable until the user logs in again, and the command
		// being run says so in its own terms. Warning here would only repeat it
		// on every invocation and point at `azion profiles --refresh`, which
		// cannot work either while the credential is bad.
		if !errors.Is(err, apiversion.ErrUnauthorized) {
			// Nothing cached and the lookup itself could not be completed: fall
			// back as the CLI always has, but say so out loud instead of only in
			// a debug log.
			logger.FInfoFlags(fact.factory.IOStreams.Out, msg.APIVersionLookupFailed,
				fact.factory.Format, fact.factory.Out)
		}
		return apiversion.V3
	}

	fact.apiVersionSource = apiVersionSource{
		source: sourceLookup, version: version, previous: cache.Version, status: status,
		age: cache.Age(now), checkedAt: cache.CheckedAt, tokenHash: apiversion.TokenHash(tok),
	}
	fact.cacheAPIVersion(&settings, version, tok, activeProfile)
	return version
}

// Where the API version for this invocation came from. Recorded during
// resolution and logged later: resolution runs inside CmdRoot, before Cobra has
// parsed --debug, so anything logged there would be dropped at the default
// level.
type versionSource string

const (
	sourceCache      versionSource = "profile cache"
	sourceLookup     versionSource = "SSO lookup"
	sourceStaleCache versionSource = "expired profile cache (SSO lookup failed)"
	sourceFallback   versionSource = "fallback (SSO lookup failed, nothing cached)"
)

type apiVersionSource struct {
	source    versionSource
	version   apiversion.Version
	previous  apiversion.Version
	status    apiversion.CacheStatus
	age       time.Duration
	checkedAt time.Time
	tokenHash string
	err       error
}

// logAPIVersionSource reports how the API version was obtained. It is called
// from persistentPreRunE, once the requested log level is in effect.
func (fact *factoryRoot) logAPIVersionSource() {
	src := fact.apiVersionSource
	if src.source == "" {
		return
	}

	fields := []zap.Field{
		zap.String("version", src.version.String()),
		zap.String("source", string(src.source)),
		zap.Duration("ttl", apiversion.TTL),
	}
	if !src.checkedAt.IsZero() {
		fields = append(fields, zap.Time("cached_at", src.checkedAt), zap.Duration("cache_age", src.age))
	}
	if src.err != nil {
		fields = append(fields, zap.Error(src.err))
	}

	switch src.source {
	case sourceCache:
		logger.Debug("API version read from the profile cache; no SSO lookup was made", fields...)
	case sourceLookup:
		fields = append(fields, zap.String("reason", string(src.status)))
		if src.previous != "" {
			fields = append(fields, zap.String("previous", src.previous.String()))
		}
		logger.Debug("API version refreshed from the SSO service and re-cached", fields...)
	default:
		logger.Debug("API version could not be refreshed from the SSO service", fields...)
	}
}

// cacheAPIVersion persists a freshly resolved generation on the active profile.
// Failing to write is not fatal: the CLI just re-resolves next time.
func (fact *factoryRoot) cacheAPIVersion(settings *token.Settings, version apiversion.Version, tok, activeProfile string) {
	settings.SetAPIVersionCache(apiversion.Refreshed(version, tok, time.Now()))
	if err := token.WriteSettings(*settings, activeProfile); err != nil {
		logger.Debug("Could not cache the resolved API version", zap.Error(err))
	}
}

// preParseGlobalFlags reads --token and --config straight from the argument
// list, before the command tree exists.
//
// CmdRoot builds the tree — and therefore resolves the API generation — before
// Cobra has parsed anything, so neither flag is available through the normal
// binding at the moment it is needed. Without this pass, `azion --token <v4>`
// builds its tree from whatever credential happens to be on disk, which is
// defect: a CI pipeline that authenticates only through --token can be shown
// the commands of the wrong generation entirely.
//
// The pass is best effort by design. It reuses the real root flag set, so it
// never drifts from it, and anything it cannot parse simply leaves the previous
// behaviour in place; Cobra parses the arguments properly moments later and is
// still the one that reports malformed input to the user.
func (fact *factoryRoot) preParseGlobalFlags(args []string) {
	probe := &cobra.Command{}
	fact.setFlags(probe)

	flags := probe.PersistentFlags()
	flags.ParseErrorsAllowlist.UnknownFlags = true
	flags.SetOutput(io.Discard)
	flags.Usage = func() {}

	if err := flags.Parse(args); err != nil {
		logger.Debug("Could not pre-parse global flags; falling back to stored configuration", zap.Error(err))
	}
}

// applyConfigFlag points the CLI at the configuration folder given by --config
// before any profile is read.
//
// Cobra applies this flag in doPreCommandCheck, which runs after the tree is
// built, so version resolution would otherwise read the cached version and the
// stored credential from the default profile even when the user asked for
// another one. Errors are ignored here on purpose: the same call runs again in
// doPreCommandCheck, which is where a bad path is reported to the user.
func (fact *factoryRoot) applyConfigFlag() {
	if fact.configFlag == "" || strings.HasPrefix(fact.configFlag, PREFIX_FLAG) {
		return
	}
	if err := config.SetPath(fact.configFlag); err != nil {
		logger.Debug("Could not apply --config before resolving the API version", zap.Error(err))
	}
}

// effectiveToken reports the credential this invocation authenticates with.
//
// An explicit --token wins, which is the whole point of the pre-parse pass.
// Otherwise the value Viper already holds is used, so the environment keeps the
// precedence it has always had.
func (fact *factoryRoot) effectiveToken() string {
	if fact.tokenFlag != "" {
		return fact.tokenFlag
	}
	return fact.factory.Config.GetString("token")
}

// credentialStatus is what the API version lookup already established about the
// credential, before any command asks about it separately.
type credentialStatus int

const (
	// credentialUnknown means nothing usable was established and the credential
	// has to be checked on its own.
	credentialUnknown credentialStatus = iota
	// credentialAccepted means the authentication service accepted it, either
	// on this invocation or within the cache TTL.
	credentialAccepted
	// credentialRejected means the service refused it.
	credentialRejected
)

// credentialStatusFor reports what resolving the API version already proved
// about the given credential.
//
// Resolution authenticates against the same service with the same token, so a
// successful lookup is itself a validation, and a cached result is one that was
// valid within apiversion.TTL. That is what lets the CLI stop validating the
// token on every single invocation.
//
// The answer only counts for the credential the lookup actually used: with a
// different token in play the caller has to check for itself.
func (fact *factoryRoot) credentialStatusFor(tok string) credentialStatus {
	src := fact.apiVersionSource
	if src.source == "" || tok == "" || src.tokenHash != apiversion.TokenHash(tok) {
		return credentialUnknown
	}

	switch src.source {
	case sourceCache, sourceLookup:
		return credentialAccepted
	case sourceFallback:
		if errors.Is(src.err, apiversion.ErrUnauthorized) {
			return credentialRejected
		}
	}

	// sourceStaleCache means the service could not be reached, so nothing new
	// was established about the credential.
	return credentialUnknown
}
