package apiversion

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheLookup(t *testing.T) {
	now := time.Now()
	const tok = "a-token"

	fresh := Refreshed(V4, tok, now)

	tests := []struct {
		name       string
		cache      Cache
		token      string
		at         time.Time
		want       Version
		wantStatus CacheStatus
	}{
		{
			name:       "a fresh entry for the same credential hits",
			cache:      fresh,
			token:      tok,
			at:         now.Add(time.Hour),
			want:       V4,
			wantStatus: CacheHit,
		},
		{
			name:       "an entry just under the TTL still hits",
			cache:      fresh,
			token:      tok,
			at:         now.Add(TTL - time.Minute),
			want:       V4,
			wantStatus: CacheHit,
		},
		{
			name:       "an entry exactly at the TTL expires",
			cache:      fresh,
			token:      tok,
			at:         now.Add(TTL),
			wantStatus: CacheMissExpired,
		},
		{
			name:       "an entry past the TTL expires",
			cache:      fresh,
			token:      tok,
			at:         now.Add(TTL + time.Minute),
			wantStatus: CacheMissExpired,
		},
		{
			name:       "a different credential misses",
			cache:      fresh,
			token:      "another-token",
			at:         now,
			wantStatus: CacheMissCredential,
		},
		{
			name:       "an empty cache misses",
			cache:      Cache{},
			token:      tok,
			at:         now,
			wantStatus: CacheMissEmpty,
		},
		{
			name:       "a version with no timestamp misses",
			cache:      Cache{Version: V4, TokenHash: TokenHash(tok)},
			token:      tok,
			at:         now,
			wantStatus: CacheMissEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, status := tt.cache.Lookup(tt.token, tt.at)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantStatus == CacheHit, status.Hit())
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCacheStale(t *testing.T) {
	old := Refreshed(V4, "a-token", time.Now().Add(-90*24*time.Hour))

	got, ok := old.Stale()
	assert.True(t, ok, "an expired entry is still an answer")
	assert.Equal(t, V4, got)

	_, ok = Cache{}.Stale()
	assert.False(t, ok, "an empty cache has no answer")
}

// The token itself must never be written to disk.
func TestTokenHashDoesNotLeakTheToken(t *testing.T) {
	const tok = "azion-super-secret-token"
	h := TokenHash(tok)

	assert.NotContains(t, h, tok)
	assert.Len(t, h, 64)
	assert.Equal(t, h, TokenHash(tok), "hashing is stable")
	assert.NotEqual(t, h, TokenHash(tok+"x"))
}

// The miss reason is what the debug line reports, so it must name the actual
// cause rather than a generic miss.
func TestCacheStatusExplainsTheMiss(t *testing.T) {
	now := time.Now()
	const tok = "a-token"

	_, status := Cache{}.Lookup(tok, now)
	assert.Equal(t, CacheMissEmpty, status)
	assert.Contains(t, string(status), "no version cached")

	_, status = Refreshed(V4, "other", now).Lookup(tok, now)
	assert.Equal(t, CacheMissCredential, status)
	assert.Contains(t, string(status), "different credential")

	_, status = Refreshed(V4, tok, now.Add(-TTL)).Lookup(tok, now)
	assert.Equal(t, CacheMissExpired, status)
	assert.Contains(t, string(status), "older than the TTL")
}

func TestCacheAge(t *testing.T) {
	now := time.Now()
	assert.Equal(t, 3*time.Hour, Refreshed(V4, "t", now.Add(-3*time.Hour)).Age(now))
	assert.Zero(t, Cache{}.Age(now), "an empty cache has no age")
}
