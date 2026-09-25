package personaltoken

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorList = errors.New("Failed to list your personal tokens: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
