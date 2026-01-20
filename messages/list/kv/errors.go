package kv

import "errors"

var (
	ErrorGetKv = errors.New("Failed to list your KV namespaces: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
