package wafexceptions

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetAll       = "failed to list WAF Exceptions: %w"
	ErrorConvertWafId = errors.New("invalid WAF ID. The value must be an integer")
)
