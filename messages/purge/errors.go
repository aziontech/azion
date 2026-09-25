package purge

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorTooManyUrls = errors.New("Only one item is allowed for the Wildcard option")
)
