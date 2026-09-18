package token

import "github.com/aziontech/azion-cli/pkg/apiversion"

// APIVersionCache reads the profile's cached API generation.
func (s Settings) APIVersionCache() apiversion.Cache {
	return apiversion.Cache{
		Version:   apiversion.Version(s.APIVersion),
		CheckedAt: s.APIVersionCheckedAt,
		TokenHash: s.APIVersionTokenHash,
	}
}

// SetAPIVersionCache records a freshly resolved API generation on the profile.
func (s *Settings) SetAPIVersionCache(c apiversion.Cache) {
	s.APIVersion = c.Version.String()
	s.APIVersionCheckedAt = c.CheckedAt
	s.APIVersionTokenHash = c.TokenHash
}
