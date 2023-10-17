package domains

import "errors"

var (
	ErrorGetDomains = errors.New("Failed to list your domains. Check your settings and try again. If the error persists, contact Azion support.")
)
