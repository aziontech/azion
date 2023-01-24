package domains

import "errors"

var (
	ErrorGetDomains = errors.New("Failed to get the Domains: %s. Check your settings and try again. If the error persists, contact Azion support")
)
