package http

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorRequest = errors.New("Error while requesting graphql api")
)
