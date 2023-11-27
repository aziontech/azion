package edge_applications

import "errors"

var (
	ErrorGetAll = errors.New("Failed to list your Edge Applications: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
