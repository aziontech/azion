package applications

import "errors"

// Used by the v3, v4 and v6 command trees.
var (
	ErrorGetAll = errors.New("Failed to list your Applications: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
