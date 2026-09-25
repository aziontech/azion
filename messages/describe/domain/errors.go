package domain

import "errors"

// Used by both the v3 and the v4 command trees.
var ErrorGetDomain = errors.New("Failed to describe the Domain: %s. Check your settings and try again. If the error persists, contact Azion support.")
