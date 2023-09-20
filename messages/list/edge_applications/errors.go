package edge_applications

import "errors"

var (
	ErrorGetAll = errors.New("Failed to list the edge application: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
