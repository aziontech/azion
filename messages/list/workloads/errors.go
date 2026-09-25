package workloads

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetDomains = errors.New("Failed to list your workloads. Check your settings and try again. If the error persists, contact Azion support.")
)
