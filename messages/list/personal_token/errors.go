package personaltoken

import "errors"

var (
	ErrorList = errors.New("Failed to list your personal tokens: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
