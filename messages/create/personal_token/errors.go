package personaltoken

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorCreate            = errors.New("Failed to create the Personal Token: %s. Check your settings and try again. If the error persists, contact Azion support.")
	ErrorMissingExpiration = errors.New("Failed to create the Personal Token: You must provide an expiration value.")
)
