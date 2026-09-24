package wafexceptions

import "errors"

// Used only by the v4 command tree.
var (
	ErrorFailToDeleteException = errors.New("failed to delete the WAF Exception: %w")
	ErrorConvertWafID          = errors.New("invalid WAF ID. The value must be an integer")
	ErrorConvertExceptionID    = errors.New("invalid WAF Exception ID. The value must be an integer")
)
