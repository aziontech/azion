package variables

import "errors"

var (
	ErrorGetVariables = errors.New("Failed to describe the origins: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
