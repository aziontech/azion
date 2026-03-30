package wafexceptions

import "errors"

var (
	ErrorGetAll       = "failed to list WAF Exceptions: %w"
	ErrorConvertWafId = errors.New("invalid WAF ID. The value must be an integer")
)
