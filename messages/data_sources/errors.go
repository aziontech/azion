package datasources

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetDataSources = errors.New("Failed to list the data sources: %s. Check your settings and try again. If the error persists, contact Azion support")
)
