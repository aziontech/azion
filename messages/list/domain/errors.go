package domain

import "errors"

// Used only by the v3 command tree (bug-fix-only; see doc/plan.md).
var (
	ErrorGetDomains = errors.New("Failed to list your domains. Check your settings and try again. If the error persists, contact Azion support.")
)
