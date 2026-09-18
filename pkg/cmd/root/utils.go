package root

import (
	"time"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap"
)

// resolveAPIVersion decides which command tree to build for this invocation.
//
// The resolved generation is cached in the active profile for apiversion.TTL,
// so a normal invocation makes no network call at all. The cache is bound to
// the credential that produced it, so a new login, a profile switch or a
// different token re-resolves instead of trusting a stale answer.
func (fact *factoryRoot) resolveAPIVersion() apiversion.Version {
	tok := fact.factory.Config.GetString("token")

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
			age: cache.Age(now), checkedAt: cache.CheckedAt,
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
				age: cache.Age(now), checkedAt: cache.CheckedAt, err: err,
			}
			return cached
		}

		// Nothing cached: fall back as the CLI always has, but say so out loud
		// instead of only in a debug log.
		fact.apiVersionSource = apiVersionSource{source: sourceFallback, version: apiversion.V3, status: status, err: err}
		logger.FInfoFlags(fact.factory.IOStreams.Out, msg.APIVersionLookupFailed,
			fact.factory.Format, fact.factory.Out)
		return apiversion.V3
	}

	fact.apiVersionSource = apiVersionSource{
		source: sourceLookup, version: version, previous: cache.Version, status: status,
		age: cache.Age(now), checkedAt: cache.CheckedAt,
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
