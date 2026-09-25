package applications

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetAll = errors.New("Failed to list your Applications: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
