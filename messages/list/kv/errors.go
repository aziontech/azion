package kv

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorGetKv = errors.New("Failed to list your KV namespaces: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
