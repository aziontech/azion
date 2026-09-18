package apiversion

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// TTL is how long a resolved version stays good before the CLI re-checks with
// the SSO service. It matches the window already used for the update/metrics
// check in pkg/cmd/root/pre_command.go.
const TTL = 24 * time.Hour

// Cache is a resolved version as persisted in a profile: what was resolved,
// when, and which credential produced it.
//
// It is deliberately free of any dependency on pkg/token; callers map it to and
// from token.Settings.
type Cache struct {
	Version   Version
	CheckedAt time.Time
	TokenHash string
}

// TokenHash binds a cached version to the credential that produced it. Only the
// hash is ever stored, never the token itself.
func TokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CacheStatus is the outcome of a cache lookup. It doubles as the human-readable
// reason a lookup missed, so callers can log why they are re-checking.
type CacheStatus string

const (
	// CacheHit means the cached version was used and no lookup is needed.
	CacheHit CacheStatus = "hit"
	// CacheMissEmpty means nothing has been cached for this profile yet.
	CacheMissEmpty CacheStatus = "no version cached yet"
	// CacheMissCredential means the entry was produced by a different
	// credential: a new login, a profile switch, or another token.
	CacheMissCredential CacheStatus = "cached for a different credential"
	// CacheMissExpired means the entry is older than the TTL.
	CacheMissExpired CacheStatus = "cached version older than the TTL"
)

// Hit reports whether the lookup produced a usable version.
func (s CacheStatus) Hit() bool { return s == CacheHit }

// Lookup returns the cached version when it is still usable for this token,
// along with the reason when it is not.
func (c Cache) Lookup(token string, now time.Time) (Version, CacheStatus) {
	if !c.exists() {
		return "", CacheMissEmpty
	}
	if c.TokenHash != TokenHash(token) {
		return "", CacheMissCredential
	}
	if c.Age(now) >= TTL {
		return "", CacheMissExpired
	}
	return c.Version, CacheHit
}

// Age reports how long ago the entry was resolved. It is zero for an empty
// cache.
func (c Cache) Age(now time.Time) time.Duration {
	if c.CheckedAt.IsZero() {
		return 0
	}
	return now.Sub(c.CheckedAt)
}

// Stale reports the cached version regardless of age, along with whether there
// is one at all. A version that is merely old is still a better answer than
// silently downgrading the account when the SSO lookup cannot be completed.
func (c Cache) Stale() (Version, bool) {
	if !c.exists() {
		return "", false
	}
	return c.Version, true
}

func (c Cache) exists() bool {
	return c.Version != "" && !c.CheckedAt.IsZero()
}

// Refreshed returns the entry to persist after a successful lookup.
func Refreshed(v Version, token string, now time.Time) Cache {
	return Cache{Version: v, CheckedAt: now, TokenHash: TokenHash(token)}
}
